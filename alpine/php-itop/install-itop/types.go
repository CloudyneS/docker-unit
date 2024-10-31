package main

import (
	"encoding/xml"
	"os"
)

type ItopInstallConfig struct {
	xml.Name                  `xml:"installation"`
	Mode                      string                        `xml:"mode"`
	Preinstall                ItopInstallConfigPreinstall   `xml:"preinstall"`
	SourceDir                 string                        `xml:"source_dir"`
	DatamodelVersion          string                        `xml:"datamodel_version"`
	PreviousConfigurationFile string                        `xml:"previous_configuration_file"`
	ExtensionsDir             string                        `xml:"extensions_dir"`
	TargetEnv                 string                        `xml:"target_env"`
	WorkspaceDir              string                        `xml:"workspace_dir"`
	Database                  ItopInstallConfigDatabase     `xml:"database"`
	Url                       string                        `xml:"url"`
	GraphvizPath              string                        `xml:"graphviz_path"`
	AdminAccount              ItopInstallConfigAdminAccount `xml:"admin_account"`
	Language                  string                        `xml:"language"`
	SelectedModules           []string                      `xml:"selected_modules>item"`
	SelectedExtensions        []string                      `xml:"selected_extensions>item"`
	SampleData                string                        `xml:"sample_data"`
	OldAddon                  string                        `xml:"old_addon"`
	Options                   map[string]string             `xml:"options"`
	MysqlBindir               string                        `xml:"mysql_bindir"`
}

type ItopInstallConfigPreinstall struct {
	Copies ItopInstallConfigPreinstallCopies `xml:"copies"`
}

type ItopInstallConfigPreinstallCopies struct {
	Type string `xml:"type,attr"`
}

var DefaultItopInstallConfigPreinstall = ItopInstallConfigPreinstall{
	Copies: ItopInstallConfigPreinstallCopies{
		Type: "array",
	},
}

type ItopInstallConfigDatabase struct {
	Server       string `xml:"server"`
	User         string `xml:"user"`
	Pwd          string `xml:"pwd"`
	Name         string `xml:"name"`
	DbTlsEnabled string `xml:"db_tls_enabled"`
	DbTlsCa      string `xml:"db_tls_ca"`
	Prefix       string `xml:"prefix"`
}

type ItopInstallConfigAdminAccount struct {
	User     string `xml:"user"`
	Pwd      string `xml:"pwd"`
	Language string `xml:"language"`
}

var DefaultItopInstallConfig = ItopInstallConfig{
	Mode:                      "upgrade",
	Preinstall:                DefaultItopInstallConfigPreinstall,
	SourceDir:                 "datamodels/2.x/",
	DatamodelVersion:          "3.0.0",
	PreviousConfigurationFile: "/app/config/production/config-itop.php",
	ExtensionsDir:             "extensions",
	TargetEnv:                 "production",
	WorkspaceDir:              "",
	Database: ItopInstallConfigDatabase{
		Server:       "",
		User:         "",
		Pwd:          "",
		Name:         "",
		DbTlsEnabled: "",
		DbTlsCa:      "",
		Prefix:       "",
	},
	Url:          "",
	GraphvizPath: "/usr/bin/dot",
	AdminAccount: ItopInstallConfigAdminAccount{
		User:     "",
		Pwd:      "",
		Language: "",
	},
	Language: "",
	SelectedModules: []string{
		"authent-cas",
		"authent-external",
		"authent-ldap",
		"authent-local",
		"combodo-backoffice-darkmoon-theme",
		"itop-backup",
		"itop-config",
		"itop-files-information",
		"itop-portal-base",
		"itop-profiles-itil",
		"itop-sla-computation",
		"itop-structure",
		"itop-welcome-itil",
		"itop-config-mgmt",
		"itop-attachments",
		"itop-tickets",
		"combodo-db-tools",
		"combodo-webhook-integration",
		"itop-core-update",
		"itop-hub-connector",
		"itop-endusers-devices",
		"itop-service-mgmt-provider",
		"itop-bridge-cmdb-ticket",
		"itop-faq-light",
		"itop-knownerror-mgmt",
		"combodo-saml",
	},
	SelectedExtensions: []string{
		"itop-config-mgmt-core",
		"itop-config-mgmt-end-user",
		"itop-service-mgmt-service-provider",
		"itop-ticket-mgmt-none",
		"itop-change-mgmt-none",
		"itop-known-error-mgmt-none",
		"combodo-saml",
	},
	SampleData:  "",
	OldAddon:    "",
	Options:     map[string]string{},
	MysqlBindir: "",
}

func GetConfig(dbConfig *ItopInstallConfigDatabase, url string, adminAccount *ItopInstallConfigAdminAccount, language string) *ItopInstallConfig {

	DefaultItopInstallConfig.Database = *dbConfig
	DefaultItopInstallConfig.Url = url
	DefaultItopInstallConfig.AdminAccount = *adminAccount
	DefaultItopInstallConfig.Language = language

	return &DefaultItopInstallConfig
}

func (config *ItopInstallConfig) SetModules(modules []string) {
	config.SelectedModules = modules
}

func (config *ItopInstallConfig) SetExtensions(extensions []string) {
	config.SelectedExtensions = extensions
}

func (config *ItopInstallConfig) AddModule(module string) {
	config.SelectedModules = append(config.SelectedModules, module)
}

func (config *ItopInstallConfig) AddExtension(extension string) {
	config.SelectedExtensions = append(config.SelectedExtensions, extension)
}

func (config *ItopInstallConfig) Save(path string) error {
	xml, err := xml.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, xml, 0644)
}
