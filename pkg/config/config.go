package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

const ConfigFileName = "gogo.toml"

// Config represents the gogo.toml configuration file
type Config struct {
	Database DatabaseConfig `toml:"database"`
}

// DatabaseConfig contains database-related configuration
type DatabaseConfig struct {
	Migrations string          `toml:"migrations"`
	SQLBoiler  SQLBoilerConfig `toml:"sqlboiler"`
}

// SQLBoilerConfig contains sqlboiler-specific configuration
type SQLBoilerConfig struct {
	PkgName   string   `toml:"pkgname"`
	Output    string   `toml:"output"`
	Wipe      bool     `toml:"wipe"`
	Blacklist []string `toml:"blacklist"`
}

// Load reads and parses the gogo.toml config file from the given directory
func Load(dir string) (*Config, error) {
	configPath := filepath.Join(dir, ConfigFileName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	if cfg.Database.Migrations == "" {
		cfg.Database.Migrations = "migrations"
	}
	if cfg.Database.SQLBoiler.PkgName == "" {
		cfg.Database.SQLBoiler.PkgName = "models"
	}
	if cfg.Database.SQLBoiler.Output == "" {
		cfg.Database.SQLBoiler.Output = "pkg/models"
	}

	return &cfg, nil
}
