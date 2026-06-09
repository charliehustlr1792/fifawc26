package cache

import (
	"errors"
	"time"
)

var ErrMiss = errors.New("cache miss")

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte, ttl time.Duration) error
	Close() error
}