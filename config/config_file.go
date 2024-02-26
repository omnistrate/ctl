package config

import (
	"bytes"
	"errors"
	"github.com/mitchellh/go-homedir"

	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"path/filepath"
)

// AuthType auth type
type AuthType string

const (
	//JWTAuthType jwt authentication type
	JWTAuthType = "jwt_auth"

	// ConfigLocationEnv is the name of the env variable used
	// to configure the location of the omnistrate-cli config folder.
	// When not set, DefaultDir location is used.
	ConfigLocationEnv string = "OMNISTRATE_CONFIG"

	DefaultDir         string      = "~/.omnistrate"
	DefaultFile        string      = "config.yml"
	DefaultPermissions os.FileMode = 0700

	// DefaultCIDir creates the 'omnistrate' directory in the current directory
	// if running in a CI environment.
	DefaultCIDir string = ".omnistrate"
	// DefaultCIPermissions creates the config file with elevated permissions
	// for it to be read by multiple users when running in a CI environment.
	DefaultCIPermissions os.FileMode = 0744
)

// ConfigFile for Omnistrate CLI exclusively.
type ConfigFile struct {
	AuthConfigs []AuthConfig `yaml:"auths"`
	FilePath    string       `yaml:"-"`
}

type AuthConfig struct {
	Auth    AuthType `yaml:"auth,omitempty"`
	Email   string   `yaml:"email,omitempty"`
	Token   string   `yaml:"token,omitempty"`
	Options []Option `yaml:"options,omitempty"`
}

type Option struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

var ErrConfigNotFound = errors.New("config file not found")

type AuthConfigNotFoundError struct {
}

func (e *AuthConfigNotFoundError) Error() string {
	return "no auth config found"
}

// New initializes a config file for the given file path
func New(filePath string) (*ConfigFile, error) {
	if filePath == "" {
		return nil, fmt.Errorf("can't create config with empty filePath")
	}
	conf := &ConfigFile{
		AuthConfigs: make([]AuthConfig, 0),
		FilePath:    filePath,
	}

	return conf, nil
}

// ConfigDir returns the path to the omnistrate-cli config directory.
// When
// 1. CI = "true" and OMNISTRATE_CONFIG="", then it will return `.omnistrate`, which is located in the current working directory.
// 2. CI = "true" and OMNISTRATE_CONFIG="<path>", then it will return the path value in  OMNISTRATE_CONFIG
// 3. CI = "" and OMNISTRATE_CONFIG="", then it will return the default location ~/.omnistrate
func ConfigDir() string {
	override := os.Getenv(ConfigLocationEnv)
	ci := isRunningInCI()

	switch {
	// case (1) from docs string
	case ci && override == "":
		return DefaultCIDir
	// case (2) from the doc string
	case override != "":
		// case (3) from the doc string
		return override
	default:
		return DefaultDir
	}
}

// isRunningInCI checks the ENV var CI and returns true if it's set to true or 1
func isRunningInCI() bool {
	if env, ok := os.LookupEnv("CI"); ok {
		if env == "true" || env == "1" {
			return true
		}
	}
	return false
}

// EnsureFile creates the root dir and config file
func EnsureFile() (string, error) {
	permission := DefaultPermissions
	dir := ConfigDir()
	if isRunningInCI() {
		permission = DefaultCIPermissions
	}
	dirPath, err := homedir.Expand(dir)
	if err != nil {
		return "", err
	}

	filePath := path.Clean(filepath.Join(dirPath, DefaultFile))
	if err = os.MkdirAll(filepath.Dir(filePath), permission); err != nil {
		return "", err
	}

	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		var file *os.File
		file, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return "", err
		}
		defer file.Close()
	}

	return filePath, nil
}

// FileExists returns true if the config file is located at the default path
func fileExists() bool {
	dir := ConfigDir()
	dirPath, err := homedir.Expand(dir)
	if err != nil {
		return false
	}

	filePath := path.Clean(filepath.Join(dirPath, DefaultFile))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

// Save writes the config to disk
func (configFile *ConfigFile) save() error {
	file, err := os.OpenFile(configFile.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	var buff bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&buff)
	yamlEncoder.SetIndent(2) // this is what you're looking for
	if err = yamlEncoder.Encode(&configFile); err != nil {
		return err
	}

	_, err = file.Write(buff.Bytes())
	return err
}

// Load reads the yml file from disk
func (configFile *ConfigFile) load() error {
	conf := &ConfigFile{}

	if _, err := os.Stat(configFile.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("can't load config from non existent filePath")
	}

	data, err := os.ReadFile(configFile.FilePath)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(data, conf); err != nil {
		return err
	}

	if len(conf.AuthConfigs) > 0 {
		configFile.AuthConfigs = conf.AuthConfigs
	}
	return nil
}

// UpdateAuthConfig creates or updates the username and password
func UpdateAuthConfig(authConfig AuthConfig) error {
	configPath, err := EnsureFile()
	if err != nil {
		return err
	}

	cfg, err := New(configPath)
	if err != nil {
		return err
	}

	if err = cfg.load(); err != nil {
		return err
	}

	if len(cfg.AuthConfigs) == 0 {
		cfg.AuthConfigs = append(cfg.AuthConfigs, authConfig)
	} else {
		cfg.AuthConfigs[0] = authConfig
	}

	if err = cfg.save(); err != nil {
		return err
	}

	return nil
}

// LookupAuthConfig returns the username and password
func LookupAuthConfig() (AuthConfig, error) {
	var authConfig AuthConfig

	if !fileExists() {
		return authConfig, ErrConfigNotFound
	}

	configPath, err := EnsureFile()
	if err != nil {
		return authConfig, err
	}

	cfg, err := New(configPath)
	if err != nil {
		return authConfig, err
	}

	if err = cfg.load(); err != nil {
		return authConfig, err
	}

	if len(cfg.AuthConfigs) > 0 {
		return cfg.AuthConfigs[0], nil
	}

	return authConfig, &AuthConfigNotFoundError{}
}
