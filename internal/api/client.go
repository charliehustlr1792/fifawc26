package api

import "context"

type Provider interface {
	GetCompetition(ctx context.Context, code string) (*Competition, error)
}