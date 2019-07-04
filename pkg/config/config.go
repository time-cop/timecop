package config

import (
	"os"
	"fmt"
	"path/filepath"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Paths struct {
	RootPath	string
	DatabasePath	string
}

type Config struct {
	SnoozeInMinutes	float32	`yaml:"snooze"`
	DefaultTaskLengthInMinutes	float32	`yaml:"defaultTaskLength"`
}


var (
	AppPaths Paths
	AppConfig Config
	ConfigFile	string
	DefaultConfig Config = Config{
		SnoozeInMinutes: 5.0,
		DefaultTaskLengthInMinutes: 20.0,
	}
)

func envOrDefault(varName, def string) string {
	// Don't use `os.LookupEnv()` as we want default when varName is empty
	result := os.Getenv(varName)
	if result == "" {
		return def
	}
	return result
}

// Init kicks off configuration loading and processing
// This is not `init()` as we call it manually, allowing for CLI overrides (eg. of file locations)
// TODO: refactor to be pure function with more override checks
func Init() error {
	xdgConfigPath := envOrDefault("XDG_CONFIG_HOME", "$HOME/.config")
	AppPaths = Paths{
		RootPath:	os.ExpandEnv(filepath.Join(xdgConfigPath, "timecop")),
		DatabasePath:	os.ExpandEnv(filepath.Join(xdgConfigPath, "databases")),
	}
	paths := []string{
		AppPaths.RootPath,
	}
	for _, path := range paths {
		err := os.MkdirAll(path, os.ModePerm)
		if os.IsExist(err) {
			fmt.Printf("path exists: %s", path)
		} else if err != nil {
			return fmt.Errorf("Unhandled error creating path '%s': %#v", path, err)
		}
	}
	ConfigFile = filepath.Join(AppPaths.RootPath, "timecop.yml")

	fileInfo, err := os.Stat(ConfigFile)
	if err != nil || fileInfo.IsDir() {
		// if the file existed, then something odd happened (excuse the double negative)
		if !os.IsNotExist(err) {
			return fmt.Errorf("Unhandled error loading config file: %#v", err)
		}
		// otherwise, just create a default file which we'll redundantly re-load in next code block
		data, err := yaml.Marshal(&DefaultConfig)
		if err != nil {
			return fmt.Errorf("Couldn't marshal default config struct: %#v", err)
		}
		err = ioutil.WriteFile(ConfigFile, data, os.ModePerm)
		if err != nil {
			return fmt.Errorf("Couldn't write default config file: %#v", err)
		}
	}

	configData, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		return fmt.Errorf("Couldn't read config file: %#v", err)
	}
	err = yaml.Unmarshal(configData, &AppConfig)
	if err != nil {
		return fmt.Errorf("Error parsing config file: %#v", err)
	}

	return nil
}
