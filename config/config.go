package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type SourceType string

const (
	RSS  SourceType = "rss"
	ATOM SourceType = "atom"
	JSON SourceType = "json"
	HTML SourceType = "html"
)

type Source struct {
	Name string
	Type SourceType
	URL  string
}

type Config struct {
	ConfigFilePath string
	Sources        []Source
}

func NewConfig(configFilePath string) (*Config, error) {
	var config = &Config{}
	err := config.loadConfigFile(configFilePath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) loadConfigFile(path string) error {
	configFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configFile, c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) Endpoints() []string {
	endpoints := make([]string, 0)
	for _, source := range c.Sources {
		endpoints = append(endpoints, source.URL)
	}
	return endpoints
}
