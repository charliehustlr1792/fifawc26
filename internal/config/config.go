package config

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	APIKey       string
	PollInterval time.Duration
	CacheDir     string
}

func Load() (*Config, error) {
	apiKey := os.Getenv("FIFAWC26_API_KEY")
	if apiKey == "" {
		return nil, errors.New("FIFAWC26_API_KEY environment variable is not set")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cacheDir := filepath.Join(home, ".fifawc26")

	return &Config{
		APIKey:       apiKey,
		PollInterval: 45 * time.Second,
		CacheDir:     cacheDir,
	}, nil
}