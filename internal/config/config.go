package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Tier string

const (
	TierKeyed   Tier = "keyed"
	TierKeyless Tier = "keyless"
)

type Config struct {
	Tier         Tier
	APIKey       string
	PollInterval time.Duration
	CacheDir     string
	ConfigPath   string
}

func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".fifawc26"), nil
}

func Load() (*Config, error) {
	dir, err := configDir()
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create config dir: %w", err)
	}

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(dir)
	v.SetDefault("poll_interval_seconds", 45)

	cfgPath := filepath.Join(dir, "config.toml")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("read config: %w", err)
			}
		}
	}

	if key := os.Getenv("FIFAWC26_API_KEY"); key != "" && v.GetString("api_key") == "" {
		v.Set("api_key", key)
		v.Set("tier", string(TierKeyed))
	}

	return &Config{
		Tier:         Tier(v.GetString("tier")),
		APIKey:       v.GetString("api_key"),
		PollInterval: time.Duration(v.GetInt("poll_interval_seconds")) * time.Second,
		CacheDir:     dir,
		ConfigPath:   cfgPath,
	}, nil
}

func Save(c *Config) error {
	v := viper.New()
	v.SetConfigFile(c.ConfigPath)
	v.Set("tier", string(c.Tier))
	v.Set("api_key", c.APIKey)
	v.Set("poll_interval_seconds", int(c.PollInterval/time.Second))
	return v.WriteConfig()
}

func (c *Config) NeedsOnboarding() bool {
	if c.Tier == TierKeyed {
		return c.APIKey == ""
	}
	if c.Tier == TierKeyless {
		return false
	}
	return true
}