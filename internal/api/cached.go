package api

import (
	"context"
	"encoding/json"
	"errors"
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
}

func DefaultTTLs() TTLs {
	return TTLs{
		Competition: 24 * time.Hour,
	}
}

func NewCachedProvider(inner Provider, c cache.Cache, ttls TTLs) *CachedProvider {
	return &CachedProvider{inner: inner, c: c, ttls: ttls}
}

func (p *CachedProvider) GetCompetition(ctx context.Context, code string) (*Competition, error) {
	key := "competition:" + code
	if raw, err := p.c.Get(key); err == nil {
		var comp Competition
		if json.Unmarshal(raw, &comp) == nil {
			return &comp, nil
		}
	} else if !errors.Is(err, cache.ErrMiss) {
		return nil, err
	}

	comp, err := p.inner.GetCompetition(ctx, code)
	if err != nil {
		return nil, err
	}
	if buf, mErr := json.Marshal(comp); mErr == nil {
		_ = p.c.Set(key, buf, p.ttls.Competition)
	}
	return comp, nil
}