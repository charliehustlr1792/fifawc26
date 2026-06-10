package api

import "context"

type MatchFilter struct {
	Status   string
	DateFrom string
	DateTo   string
	Matchday int
}

type Provider interface {
	GetCompetition(ctx context.Context, code string) (*Competition, error)
	GetStandings(ctx context.Context, code string) (*StandingsResponse, error)
	GetMatches(ctx context.Context, code string, f MatchFilter) (*MatchesResponse, error)
	GetScorers(ctx context.Context, code string, limit int) (*ScorersResponse, error)
	GetTeam(ctx context.Context, id int) (*TeamDetail, error)
}