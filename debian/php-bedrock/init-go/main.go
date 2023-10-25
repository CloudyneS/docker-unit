package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"syscall"
)

type Initializer interface {
	RunAsUser(username string, command string, arguments []string) error
	RunWpCli(args []string) error
	PathExists(path string) error
	DownloadFile(url string, path string) error
	InstallDatabase() error
	InstallFiles() error
	ExtractFiles(source string, destination string) error
	UpdateSalts() error
	SetTheme() error
	PerformChown() error
	PerformChmod() error
	HandleErrors(err error)
	Run() error
}

type WpInitializer struct {
	ExitOnError          bool   `json:"exitOnError" default:"false"`
	WebserverUser        string `json:"webserverUser" default:"unit"`
	WebserverGroup       string `json:"webserverGroup" default:"root"`
	ApplicationDir       string `json:"applicationDir" default:"/app"`
	Permissions          string `json:"permissions" default:"770"`
	ImportDatabase       bool   `json:"importDatabase" default:"false"`
	DatabasePath         string `json:"databasePath" default:"/mnt/cloud/import.sql"`
	DatabaseUrl          string `json:"databaseUrl" default:""`
	ImportContent        bool   `json:"importContent" default:"false"`
	ContentPath          string `json:"contentPath" default:"/mnt/cloud/import.zip"`
	ContentUrl           string `json:"contentUrl" default:""`
	OverwriteDatabase    bool   `json:"overwriteDatabase" default:"false"`
	UpdatePermissions    bool   `json:"updatePermissions" default:"true"`
	ConvertUploadsToWebp bool   `json:"convertUploadsToWebp" default:"false"`
	ConvertMissingOnly   bool   `json:"convertMissingOnly" default:"true"`
	GenerateSalts        bool   `json:"generateSalts" default:"true"`
	ActivateTheme        string `json:"activateTheme" default:""`
}

func (ini *WpInitializer) RunAsUser(username string, command string, arguments []string) error {
	u, err := user.Lookup(username)
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return err
	}

	cmd := exec.Command(command, arguments...)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("command execution failed (wp %s): %v, stderr: %v", strings.Join(arguments, ", "), err, stderr.String())
	}
	fmt.Println(stdout.String())

	return nil
}

func (ini *WpInitializer) RunWpCli(args []string) error {
	return ini.RunAsUser(
		ini.WebserverUser, "wp", args,
	)
}

func (ini *WpInitializer) PathExists(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("path %s does not exist", path)
	}
	return err
}

func (ini *WpInitializer) DownloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (ini *WpInitializer) InstallDatabase() error {
	if ini.ImportDatabase {
		if ini.DatabasePath == "" {
			if ini.DatabaseUrl == "" {
				return fmt.Errorf("no path to import the database from specified")
			}
			// Download file
			ini.DatabasePath = "/tmp/temp-import-database.sql"
			err := ini.DownloadFile(ini.DatabasePath, ini.DatabaseUrl)
			if err != nil {
				return err
			}
		}
		err := ini.PathExists(ini.DatabasePath)
		if err != nil {
			return err
		}

		if ini.OverwriteDatabase {
			_ = ini.RunWpCli([]string{"db", "reset", "--yes"})
		}

		isInstalled := ini.RunWpCli([]string{"core", "is-installed"})
		fmt.Println(isInstalled)

		if isInstalled != nil {
			return ini.RunWpCli([]string{"db", "import", ini.DatabasePath})
		}
	}
	return nil
}

func (ini *WpInitializer) ExtractFiles(source string, destination string) error {
	// Check if string ends with .zip
	if strings.HasSuffix(source, ".zip") {
		return ini.RunAsUser(
			ini.WebserverUser,
			"unzip",
			[]string{
				"-o",
				"-d",
				destination,
				source,
			},
		)
	}

	if strings.HasSuffix(source, ".tar.gz") {
		return ini.RunAsUser(
			ini.WebserverUser,
			"tar",
			[]string{
				"--overwrite",
				"--no-overwrite-dir",
				"-xzf",
				source,
				"-C",
				destination,
			},
		)
	}

	return fmt.Errorf("unsupported file format")
}

func (ini *WpInitializer) InstallFiles() error {
	if ini.ImportContent {
		if ini.ContentPath == "" {
			if ini.ContentUrl == "" {
				return fmt.Errorf("no path to import the content from specified")
			}
			// Download file
			ini.ContentPath = "/tmp/temp-import-content.zip"
			if strings.HasSuffix(ini.ContentUrl, ".tar.gz") {
				ini.ContentPath = "/tmp/temp-import-content.tar.gz"
			}
			err := ini.DownloadFile(ini.ContentPath, ini.ContentUrl)
			if err != nil {
				return err
			}
		}
		err := ini.PathExists(ini.ContentPath)
		if err != nil {
			return err
		}

		return ini.ExtractFiles(ini.ContentPath, "/app/web/app/")
	}
	return nil
}

func (ini *WpInitializer) UpdateSalts() error {
	if ini.GenerateSalts {
		return ini.RunWpCli([]string{"dotenv", "salts", "generate", "--force"})
	}
	return nil
}

func (ini *WpInitializer) SetTheme() error {
	if ini.ActivateTheme != "" {
		return ini.RunWpCli([]string{"theme", "activate", ini.ActivateTheme})
	}
	return nil
}

func (ini *WpInitializer) PerformChown() error {
	if ini.UpdatePermissions {
		return ini.RunAsUser(
			"root",
			"chown",
			[]string{
				"-R",
				ini.WebserverUser + ":" + ini.WebserverGroup,
				ini.ApplicationDir,
			},
		)
	}
	return nil
}

func (ini *WpInitializer) PerformChmod() error {
	if ini.UpdatePermissions {
		return ini.RunAsUser(
			"root",
			"chmod",
			[]string{
				"-R",
				ini.Permissions,
				ini.ApplicationDir,
			},
		)
	}
	return nil
}

func (ini *WpInitializer) HandleErrors(err error) {
	if err != nil {
		fmt.Println(err)
		if ini.ExitOnError {
			os.Exit(1)
		}
	}
}

func (ini *WpInitializer) Run() {
	fmt.Println("Installing Database")
	ini.HandleErrors(ini.InstallDatabase())

	fmt.Println("Installing Files")
	ini.HandleErrors(ini.InstallFiles())

	fmt.Println("Setting Permissions")
	ini.HandleErrors(ini.PerformChmod())

	fmt.Println("Setting Ownership")
	ini.HandleErrors(ini.PerformChown())

	// Handle WebP Conversions
	fmt.Println("Converting Uploads to WebP")
	if ini.ConvertUploadsToWebp {
		if ini.ConvertMissingOnly {
			ini.HandleErrors(ini.RunWpCli([]string{"cloudyne-webp", "convert"}))
		} else {
			ini.HandleErrors(ini.RunWpCli([]string{"cloudyne-webp", "convert", "--force-all=true"}))
		}
	}

	fmt.Println("Updating Salts")
	ini.HandleErrors(ini.UpdateSalts())

	fmt.Println("Activating Theme")
	ini.HandleErrors(ini.SetTheme())
}

func main() {
	fmt.Println("Starting initialization")
	path := ""

	if val, ok := os.LookupEnv("CD_CONFIG"); ok {
		path = val
	} else {
		fmt.Println("No configuration file specified via environment variable, checking CLI arguments...")
		if len(os.Args) < 2 {
			fmt.Println("No configuration file specified via CLI, exiting...")
			os.Exit(1)
		} else {
			path = os.Args[1]
		}
	}

	fmt.Println("Using configuration file: " + path)
	initializer := &WpInitializer{}

	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error opening configuration file: " + err.Error())
		os.Exit(1)
	}

	err = json.Unmarshal(file, initializer)
	if err != nil {
		fmt.Println("Error parsing configuration file: " + err.Error())
		os.Exit(1)
	}

	initializer.Run()
}
