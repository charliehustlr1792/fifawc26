package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const baseURL = "https://api.football-data.org/v4"

type FootballDataClient struct {
	apiKey string
	http   *http.Client

	mu             sync.Mutex
	available      int
	resetAt        time.Time
	throttleKnown  bool
}

func NewFootballDataClient(apiKey string) *FootballDataClient {
	return &FootballDataClient{
		apiKey: apiKey,
		http:   &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *FootballDataClient) GetCompetition(ctx context.Context, code string) (*Competition, error) {
	var comp Competition
	if err := c.get(ctx, "/competitions/"+code, &comp); err != nil {
		return nil, err
	}
	return &comp, nil
}

func (c *FootballDataClient) get(ctx context.Context, path string, out any) error {
	if err := c.waitIfThrottled(ctx); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Auth-Token", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("http call: %w", err)
	}
	defer resp.Body.Close()

	c.updateThrottle(resp.Header)

	if resp.StatusCode == http.StatusTooManyRequests {
		return fmt.Errorf("rate limited (429); try again after %s", c.resetAt.Format(time.RFC3339))
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("api error %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}

func (c *FootballDataClient) updateThrottle(h http.Header) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if v := h.Get("X-Requests-Available-Minute"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			c.available = n
			c.throttleKnown = true
		}
	}
	if v := h.Get("X-RequestCounter-Reset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			c.resetAt = time.Now().Add(time.Duration(n) * time.Second)
		}
	}
}

func (c *FootballDataClient) waitIfThrottled(ctx context.Context) error {
	c.mu.Lock()
	needWait := c.throttleKnown && c.available <= 0 && time.Now().Before(c.resetAt)
	wait := time.Until(c.resetAt)
	c.mu.Unlock()

	if !needWait {
		return nil
	}
	select {
	case <-time.After(wait):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}