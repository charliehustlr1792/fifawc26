package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	bolt "go.etcd.io/bbolt"
)

var bucketName = []byte("fifawc26")

type BoltCache struct {
	db *bolt.DB
}

type entry struct {
	Value     []byte    `json:"v"`
	ExpiresAt time.Time `json:"e"`
}

func NewBoltCache(dir string) (*BoltCache, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create cache dir: %w", err)
	}
	path := filepath.Join(dir, "cache.db")

	db, err := bolt.Open(path, 0o600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("open bolt: %w", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists(bucketName)
		return e
	})
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return &BoltCache{db: db}, nil
}

func (c *BoltCache) Get(key string) ([]byte, error) {
	var out []byte
	err := c.db.View(func(tx *bolt.Tx) error {
		raw := tx.Bucket(bucketName).Get([]byte(key))
		if raw == nil {
			return ErrMiss
		}
		var e entry
		if err := json.Unmarshal(raw, &e); err != nil {
			return ErrMiss
		}
		if time.Now().After(e.ExpiresAt) {
			return ErrMiss
		}
		out = e.Value
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *BoltCache) Set(key string, value []byte, ttl time.Duration) error {
	e := entry{Value: value, ExpiresAt: time.Now().Add(ttl)}
	buf, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return c.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketName).Put([]byte(key), buf)
	})
}

func (c *BoltCache) Close() error {
	return c.db.Close()
}