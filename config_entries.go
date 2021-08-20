package main

var ConfigEntries map[string]ConfigEntry

type ConfigEntry struct {
	Name string
	Get  func(*Config) string
	Set  func(*Config, string)
}

func init() {
	entries := []ConfigEntry{
		ConfigEntry{
			Name: "api.endpoint",
			Get:  func(c *Config) string { return c.API.Endpoint },
			Set:  func(c *Config, s string) { c.API.Endpoint = s },
		},
		ConfigEntry{
			Name: "api.key",
			Get:  func(c *Config) string { return c.API.Key },
			Set:  func(c *Config, s string) { c.API.Key = s },
		},
	}

	ConfigEntries = make(map[string]ConfigEntry)
	for _, e := range entries {
		ConfigEntries[e.Name] = e
	}
}
