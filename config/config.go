package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// SourceType represents the type of content source
type SourceType string

const (
	SourceTypeRSS SourceType = "rss"
	SourceTypeApi SourceType = "api"
)

// Source represents a news source configuration
type Source struct {
	Name string     `yaml:"name"`
	Type SourceType `yaml:"type"`
	URL  string     `yaml:"url"`
}

// Config represents the application configuration
type Config struct {
	Sources []Source `yaml:"sources"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
