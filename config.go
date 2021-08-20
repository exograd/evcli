package main

import (
	"bytes"
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
	Color bool `json:"color"`
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

func (c *Config) Write() error {
	filePath := ConfigPath()

	trace("writing configuration to %s", filePath)

	return c.WriteFile(filePath)
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

func (c *Config) WriteFile(filePath string) error {
	var buf bytes.Buffer

	e := json.NewEncoder(&buf)
	e.SetIndent("", "  ")

	if err := e.Encode(c); err != nil {
		return fmt.Errorf("cannot encode configuration: %w", err)
	}

	if err := ioutil.WriteFile(filePath, buf.Bytes(), 0600); err != nil {
		return fmt.Errorf("cannot write file: %w", err)
	}

	return nil
}

func (c *Config) GetEntry(name string) (string, error) {
	e, found := ConfigEntries[name]
	if !found {
		return "", fmt.Errorf("unknown configuration entry %q", name)
	}

	return e.Get(c), nil
}

func (c *Config) SetEntry(name, value string) error {
	e, found := ConfigEntries[name]
	if !found {
		return fmt.Errorf("unknown configuration entry %q", name)
	}

	if err := e.Set(c, value); err != nil {
		return fmt.Errorf("invalid value: %v", err)
	}

	return nil
}
