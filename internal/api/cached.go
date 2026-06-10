package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/charliehustlr1792/fifawc26/internal/cache"
)

type CachedProvider struct {
	inner Provider
	c     cache.Cache
	ttls  TTLs
}

type TTLs struct {
	Competition time.Duration
	Standings   time.Duration
	Matches     time.Duration
	Scorers     time.Duration
	Team        time.Duration
}

func DefaultTTLs() TTLs {
	return TTLs{
		Competition: 24 * time.Hour,
		Standings:   5 * time.Minute,
		Matches:     60 * time.Second,
		Scorers:     5 * time.Minute,
		Team:        24 * time.Hour,
	}
}

func NewCachedProvider(inner Provider, c cache.Cache, ttls TTLs) *CachedProvider {
	return &CachedProvider{inner: inner, c: c, ttls: ttls}
}

func cacheGet[T any](c cache.Cache, key string) (*T, bool) {
	raw, err := c.Get(key)
	if err != nil {
		return nil, false
	}
	var v T
	if json.Unmarshal(raw, &v) != nil {
		return nil, false
	}
	return &v, true
}

func cacheSet(c cache.Cache, key string, v any, ttl time.Duration) {
	if buf, err := json.Marshal(v); err == nil {
		_ = c.Set(key, buf, ttl)
	}
}

func (p *CachedProvider) GetCompetition(ctx context.Context, code string) (*Competition, error) {
	key := "competition:" + code
	if v, ok := cacheGet[Competition](p.c, key); ok {
		return v, nil
	}
	v, err := p.inner.GetCompetition(ctx, code)
	if err != nil {
		return nil, err
	}
	cacheSet(p.c, key, v, p.ttls.Competition)
	return v, nil
}

func (p *CachedProvider) GetStandings(ctx context.Context, code string) (*StandingsResponse, error) {
	key := "standings:" + code
	if v, ok := cacheGet[StandingsResponse](p.c, key); ok {
		return v, nil
	}
	v, err := p.inner.GetStandings(ctx, code)
	if err != nil {
		return nil, err
	}
	cacheSet(p.c, key, v, p.ttls.Standings)
	return v, nil
}

func (p *CachedProvider) GetMatches(ctx context.Context, code string, f MatchFilter) (*MatchesResponse, error) {
	key := fmt.Sprintf("matches:%s:%s:%s:%s:%d", code, f.Status, f.DateFrom, f.DateTo, f.Matchday)
	if v, ok := cacheGet[MatchesResponse](p.c, key); ok {
		return v, nil
	}
	v, err := p.inner.GetMatches(ctx, code, f)
	if err != nil {
		return nil, err
	}
	cacheSet(p.c, key, v, p.ttls.Matches)
	return v, nil
}

func (p *CachedProvider) GetScorers(ctx context.Context, code string, limit int) (*ScorersResponse, error) {
	key := fmt.Sprintf("scorers:%s:%d", code, limit)
	if v, ok := cacheGet[ScorersResponse](p.c, key); ok {
		return v, nil
	}
	v, err := p.inner.GetScorers(ctx, code, limit)
	if err != nil {
		return nil, err
	}
	cacheSet(p.c, key, v, p.ttls.Scorers)
	return v, nil
}

func (p *CachedProvider) GetTeam(ctx context.Context, id int) (*TeamDetail, error) {
	key := fmt.Sprintf("team:%d", id)
	if v, ok := cacheGet[TeamDetail](p.c, key); ok {
		return v, nil
	}
	v, err := p.inner.GetTeam(ctx, id)
	if err != nil {
		return nil, err
	}
	cacheSet(p.c, key, v, p.ttls.Team)
	return v, nil
}

var _ = errors.New