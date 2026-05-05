package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level logpipe configuration.
type Config struct {
	Server  ServerConfig   `yaml:"server"`
	Sources []SourceConfig `yaml:"sources"`
	Sinks   []SinkConfig   `yaml:"sinks"`
}

// ServerConfig defines the HTTP/gRPC listener settings.
type ServerConfig struct {
	Addr string `yaml:"addr"`
	Port int    `yaml:"port"`
}

// SourceConfig defines an inbound log source.
type SourceConfig struct {
	Name   string            `yaml:"name"`
	Type   string            `yaml:"type"` // e.g. "docker", "file", "stdin"
	Labels map[string]string `yaml:"labels"`
}

// SinkConfig defines an outbound log destination.
type SinkConfig struct {
	Name    string            `yaml:"name"`
	Type    string            `yaml:"type"` // e.g. "stdout", "file", "http"
	Target  string            `yaml:"target"`
	Options map[string]string `yaml:"options"`
}

// Load reads and parses a YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse yaml: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config: validation failed: %w", err)
	}

	return &cfg, nil
}

// Validate performs basic semantic checks on the configuration.
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535, got %d", c.Server.Port)
	}
	names := make(map[string]bool)
	for _, s := range c.Sources {
		if s.Name == "" {
			return fmt.Errorf("each source must have a non-empty name")
		}
		if names[s.Name] {
			return fmt.Errorf("duplicate source name %q", s.Name)
		}
		names[s.Name] = true
	}
	for _, s := range c.Sinks {
		if s.Name == "" {
			return fmt.Errorf("each sink must have a non-empty name")
		}
	}
	return nil
}
