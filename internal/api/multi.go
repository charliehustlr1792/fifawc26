package api

import (
	"context"
	"fmt"
)

type MultiProvider struct {
	providers []Provider
}

func NewMultiProvider(providers ...Provider) *MultiProvider {
	return &MultiProvider{providers: providers}
}

func (m *MultiProvider) GetCompetition(ctx context.Context, code string) (*Competition, error) {
	var lastErr error
	for _, p := range m.providers {
		if c, err := p.GetCompetition(ctx, code); err == nil {
			return c, nil
		} else {
			lastErr = err
		}
	}
	return nil, errOrFallback("competition", lastErr)
}

func (m *MultiProvider) GetStandings(ctx context.Context, code string) (*StandingsResponse, error) {
	var lastErr error
	for i, p := range m.providers {
		s, err := p.GetStandings(ctx, code)
		if err != nil {
			lastErr = err
			continue
		}
		if s != nil && len(s.Standings) > 0 {
			return s, nil
		}
		if i == len(m.providers)-1 {
			if s != nil {
				return s, nil
			}
			return &StandingsResponse{Standings: []Standing{}}, nil
		}
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return &StandingsResponse{Standings: []Standing{}}, nil
}

func (m *MultiProvider) GetMatches(ctx context.Context, code string, f MatchFilter) (*MatchesResponse, error) {
	var lastErr error
	for _, p := range m.providers {
		if r, err := p.GetMatches(ctx, code, f); err == nil {
			if r != nil && len(r.Matches) > 0 {
				return r, nil
			}
			if r != nil && len(m.providers) == 1 {
				return r, nil
			}
		} else {
			lastErr = err
		}
	}
	return nil, errOrFallback("matches", lastErr)
}

func (m *MultiProvider) GetScorers(ctx context.Context, code string, limit int) (*ScorersResponse, error) {
	for _, p := range m.providers {
		if s, err := p.GetScorers(ctx, code, limit); err == nil {
			if s != nil && len(s.Scorers) > 0 {
				return s, nil
			}
			if s != nil && len(m.providers) == 1 {
				return s, nil
			}
		}
	}
	return &ScorersResponse{Scorers: []Scorer{}}, nil
}

func (m *MultiProvider) GetTeam(ctx context.Context, id int) (*TeamDetail, error) {
	var lastErr error
	for _, p := range m.providers {
		if t, err := p.GetTeam(ctx, id); err == nil {
			if t != nil && t.Name != "" {
				return t, nil
			}
		} else {
			lastErr = err
		}
	}
	return nil, errOrFallback("team", lastErr)
}

func errOrFallback(what string, last error) error {
	if last != nil {
		return fmt.Errorf("%s unavailable from all providers: %w", what, last)
	}
	return fmt.Errorf("%s unavailable from all providers", what)
}