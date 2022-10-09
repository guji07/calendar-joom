package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	Server struct {
		Host    string `yaml:"host"`
		Port    string `yaml:"port"`
		Timeout struct {
			// Server is server timeout to use for graceful shutdowns
			Server time.Duration `yaml:"server"`

			// Write is time until an http write opperation will be closed
			Write time.Duration `yaml:"write"`

			// Read is time before session with read operation will be close
			Read time.Duration `yaml:"read"`

			// Read is time before idle http session will be closed
			Idle time.Duration `yaml:"idle"`
		} `yaml:"timeout"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Url      string `yaml:"url"`
	} `yaml:"database"`
}

func NewConfig(configPath string) (*AppConfig, error) {
	config := &AppConfig{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}
