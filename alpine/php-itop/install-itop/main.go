package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

var RequiredEnv = []string{
	"DB_HOST",
	"DB_PORT",
	"DB_USER",
	"DB_PASSWORD",
}

func main() {
	osenv := os.Environ()

	envVars := make(map[string]string)

	for _, env := range osenv {
		envKey := strings.Split(env, "=")[0]

		if strings.HasPrefix(envKey, "ITOP_") {
			envKey = strings.TrimPrefix(envKey, "ITOP_")
			envKey = strings.ToUpper(envKey)

			envVars[envKey] = os.Getenv(env)
		}
	}

	for _, env := range RequiredEnv {
		if _, ok := envVars[env]; !ok {
			log.Fatalf("Environment variable %s is required", env)
		}
	}

	dbConfig := ItopInstallConfigDatabase{
		Name: "itop",
	}
	adminAccount := ItopInstallConfigAdminAccount{}
	url := ""
	extensions := []string{}
	modules := []string{}
	language := ""

	for env, value := range envVars {
		switch env {
		case "DB_HOST":
			dbConfig.Server = value + dbConfig.Server
		case "DB_PORT":
			dbConfig.Server = dbConfig.Server + ":" + value
		case "DB_USER":
			dbConfig.User = value
		case "DB_PASSWORD":
			dbConfig.Pwd = value
		case "DB_PREFIX":
			dbConfig.Prefix = value
		case "DB_DATABASE", "DB_NAME":
			dbConfig.Name = value
		case "ADMIN_USER":
			adminAccount.User = value
		case "ADMIN_PASSWORD":
			adminAccount.Pwd = value
		case "URL":
			url = value
		case "EXTENSIONS":
			extensions = strings.Split(value, ",")
		case "MODULES":
			modules = strings.Split(value, ",")
		case "LANGUAGE":
			language = value
		}
	}

	config := GetConfig(&dbConfig, url, &adminAccount, language)
	if len(extensions) > 0 {
		config.SetExtensions(extensions)
	}
	if len(modules) > 0 {
		config.SetModules(modules)
	}

	err := config.Save("/tmp/setup.xml")
	if err != nil {
		log.Fatalf("Failed to save config: %s", err)
	}

	var cmd *exec.Cmd

	if os.Stat("/app/config/production/config-itop.php"); err == nil {
		cmd = exec.Command("php", "/app/toolkit/unattended_install.php", "--response-file=/tmp/setup.xml", "--clean=0", "--use-itop-config=1")
	} else {
		cmd = exec.Command("php", "/app/toolkit/unattended_install.php", "--response-file=/tmp/setup.xml", "--clean=0")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to run command: %s", err)
	}

	os.Remove("/tmp/setup.xml")
}
