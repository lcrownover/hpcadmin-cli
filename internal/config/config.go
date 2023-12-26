package config

import (
	"fmt"
	"log/slog"
	"os"
	"os/user"

	"gopkg.in/yaml.v3"
)

type CLIConfig struct {
	Server string `yaml:"server"`
	Port   int    `yaml:"port"`
	UseTLS bool   `yaml:"use_tls"`
	APIKey string `yaml:"api_key"`

	// Though these are processed by ProcessConfig
	// they could also be provided if needed.
	Protocol string `yaml:"protocol"`
	BaseURL  string `yaml:"base_url"`

	// These are not used by the CLI, but are used by the application
	ConfigDir string `yaml:"-"`
}

// EnsureCLIConfigDir ensures that the CLI config directory exists.
// If it does not exist, it will be created.
// Returns the path to the directory.
func EnsureCLIConfigDir() (string, error) {
	slog.Debug("ensuring config directory", "method", "EnsureCLIConfigDir")
	var dir string
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	// Look config dir in this order:
	// 1. ~/.config/hpcadmin/
	dir = usr.HomeDir + "/.config/hpcadmin/"
	slog.Debug(fmt.Sprintf("checking for config dir: %s", dir), "method", "EnsureCLIConfigDir")
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		slog.Debug("found config dir", "method", "EnsureCLIConfigDir")
		return dir, nil
	}
	// 2. ~/.hpcadmin
	dir = usr.HomeDir + "/.hpcadmin"
	slog.Debug(fmt.Sprintf("checking for config dir: %s", dir), "method", "EnsureCLIConfigDir")
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		slog.Debug("found config dir", "method", "EnsureCLIConfigDir")
		return dir, nil
	}
	// Neither are found, so let's create the first one
	dir = usr.HomeDir + "/.config/hpcadmin"
	slog.Debug(fmt.Sprintf("creating config dir: %s", dir), "method", "EnsureCLIConfigDir")
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return "", err
	}
	return dir, nil
}

// ReadConfigFile reads the CLI config file and returns a CLIConfig struct.
func ReadConfigFile(configPath string) (*CLIConfig, error) {
	var cfg CLIConfig
	slog.Debug("Reading config file", "method", "ReadConfigFile", "path", configPath)
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}

	slog.Debug("Parsing YAML", "method", "ReadConfigFile", "path", configPath)
	err = yaml.Unmarshal(configData, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %v", err)
	}

	return &cfg, nil
}

// ProcessConfig processes the CLI config provided from the config file.
// It sets default values for any missing fields.
// It also sets fields based on other fields.
func ProcessConfig(config *CLIConfig) error {
	slog.Debug("processing config", "method", "ProcessConfig")
	if config.UseTLS {
		slog.Debug("TLS is enabled", "method", "ProcessConfig")
		if config.Port == 0 {
			slog.Debug("setting default port to 443", "method", "ProcessConfig")
			config.Port = 443
		}
	} else {
		slog.Debug("TLS is disabled", "method", "ProcessConfig")
		if config.Port == 0 {
			slog.Debug("setting default port to 80", "method", "ProcessConfig")
			config.Port = 80
		}
	}

	if config.Protocol == "" {
		if config.UseTLS {
			slog.Debug("setting default protocol to https", "method", "ProcessConfig")
			config.Protocol = "https"
		} else {
			slog.Debug("setting default protocol to http", "method", "ProcessConfig")
			config.Protocol = "http"
		}
	}

	if config.BaseURL == "" {
		config.BaseURL = fmt.Sprintf("%s://%s:%d", config.Protocol, config.Server, config.Port)
		slog.Debug("setting default base URL", "baseurl", config.BaseURL, "method", "ProcessConfig")
	}

	slog.Debug("config processed", "method", "ProcessConfig")
	return nil
}

// EnsureCLIConfigFile ensures that the CLI config file exists.
// If it does not exist, it will be created.
// Returns the path to the file.
func EnsureCLIConfigFile(configDir string) (string, error) {
	filePath := configDir + "/config.yaml"
	slog.Debug("ensuring config file exists", "filePath", filePath, "method", "EnsureCLIConfigFile")
	os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0600)
	return filePath, nil

}

func GetCLIConfig() (*CLIConfig, error) {
	// Ensure config directory exists
	configDir, err := EnsureCLIConfigDir()
	if err != nil {
		return nil, fmt.Errorf("reading configuration directory: %v", err)
	}

	// Ensure config file exists
	configFile, err := EnsureCLIConfigFile(configDir)
	if err != nil {
		return nil, fmt.Errorf("reading configuration file: %v", err)
	}

	// Read config file
	cfg, err := ReadConfigFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("reading configuration file: %v", err)
	}

	// Server and Port are required to be set
	if cfg.Server == "" || cfg.Port == 0 {
		return nil, fmt.Errorf("server and port must be defined in config file")
	}

	// Set config directory
	cfg.ConfigDir = configDir

	// Process config file
	err = ProcessConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("processing configuration file: %v", err)
	}

	return cfg, nil
}
