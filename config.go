package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
)

type Config struct {
	Interface InterfaceConfig `json:"interface"`
	API       APIConfig       `json:"api"`
}

type InterfaceConfig struct {
	Color bool `json:"color"`
}

type APIConfig struct {
	Endpoint string `json:"endpoint"`
}

func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	filePath := ConfigPath()

	if err := config.LoadFile(filePath); err != nil {
		return nil, errors.Wrapf(err, "cannot load %q", filePath)
	}

	return config, nil
}

func ConfigPath() string {
	if path := os.Getenv("EVCLI_CONFIG_PATH"); path != "" {
		return path
	}

	homePath, err := os.UserHomeDir()
	if err != nil {
		panic(errors.Wrap(err, "cannot locate user home directory"))
	}

	return path.Join(homePath, ".evcli", "config.json")
}

func DefaultConfig() *Config {
	return &Config{
		Interface: InterfaceConfig{
			Color: true,
		},

		API: APIConfig{
			Endpoint: "https://api.eventline.com",
		},
	}
}

func (c *Config) LoadFile(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return errors.Wrap(err, "cannot read file")
	}

	return c.LoadData(data)
}

func (c *Config) LoadData(data []byte) error {
	if err := json.Unmarshal(data, c); err != nil {
		return errors.Wrap(err, "cannot parse json data")
	}

	return nil
}
