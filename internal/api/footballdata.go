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
	"strings"
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
    return fmt.Errorf("rate limited — auto-retry will resume at %s", c.resetAt.Format("15:04:05"))
}
if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
    return fmt.Errorf("auth failed (HTTP %d) — check FIFAWC26_API_KEY", resp.StatusCode)
}
if resp.StatusCode == http.StatusNotFound {
    return fmt.Errorf("not found (HTTP 404) — endpoint or resource unavailable")
}
if resp.StatusCode < 200 || resp.StatusCode >= 300 {
    body, _ := io.ReadAll(resp.Body)
    return fmt.Errorf("api error %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
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

func (c *FootballDataClient) GetStandings(ctx context.Context, code string) (*StandingsResponse, error) {
	var out StandingsResponse
	if err := c.get(ctx, "/competitions/"+code+"/standings", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *FootballDataClient) GetMatches(ctx context.Context, code string, f MatchFilter) (*MatchesResponse, error) {
	path := "/competitions/" + code + "/matches"
	q := buildQuery(map[string]string{
		"status":   f.Status,
		"dateFrom": f.DateFrom,
		"dateTo":   f.DateTo,
		"matchday": intToStr(f.Matchday),
	})
	if q != "" {
		path += "?" + q
	}
	var out MatchesResponse
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *FootballDataClient) GetTeam(ctx context.Context, id int) (*TeamDetail, error) {
	var t TeamDetail
	if err := c.get(ctx, fmt.Sprintf("/teams/%d", id), &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (c *FootballDataClient) GetScorers(ctx context.Context, code string, limit int) (*ScorersResponse, error) {
	path := "/competitions/" + code + "/scorers"
	if limit > 0 {
		path += "?limit=" + intToStr(limit)
	}
	var out ScorersResponse
	if err := c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func buildQuery(params map[string]string) string {
	first := true
	out := ""
	for k, v := range params {
		if v == "" {
			continue
		}
		if !first {
			out += "&"
		}
		out += k + "=" + v
		first = false
	}
	return out
}

func intToStr(i int) string {
	if i == 0 {
		return ""
	}
	return strconv.Itoa(i)
}