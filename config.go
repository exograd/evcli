package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Config struct {
	Interface InterfaceConfig `json:"interface,omitempty"`
	API       APIConfig       `json:"api,omitempty"`
}

type InterfaceConfig struct {
	Color bool `json:"color,omitempty"`
}

type APIConfig struct {
	Endpoint string `json:"endpoint,omitempty"`
	Key      string `json:"key,omitempty"`
}

func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	filePath := ConfigPath()

	trace("loading configuration from %s", filePath)

	if err := config.LoadFile(filePath); err != nil {
		return nil, fmt.Errorf("cannot load %q: %w", filePath, err)
	}

	return config, nil
}

func ConfigPath() string {
	if path := os.Getenv("EVCLI_CONFIG_PATH"); path != "" {
		return path
	}

	homePath, err := os.UserHomeDir()
	if err != nil {
		die("cannot locate user home directory: %v", err)
	}

	return path.Join(homePath, ".evcli", "config.json")
}

func DefaultConfig() *Config {
	return &Config{
		Interface: InterfaceConfig{
			Color: true,
		},

		API: APIConfig{
			Endpoint: "https://api.eventline.net",
		},
	}
}

func (c *Config) LoadFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	return c.LoadData(data)
}

func (c *Config) LoadData(data []byte) error {
	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("cannot parse json data: %w", err)
	}

	return nil
}
