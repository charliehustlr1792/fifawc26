package render

import (
	"github.com/charliehustlr1792/fifawc26/internal/theme"
)

func truncate(s string, max int) string {
	if len([]rune(s)) <= max {
		return s
	}
	return string([]rune(s)[:max-1]) + "…"
}

var (
	Title           = theme.Title
	Subtle          = theme.Subtle
	StatusLive      = theme.StatusLive
	StatusScheduled = theme.StatusScheduled
	StatusFinished  = theme.StatusFinished
	StatusOther     = theme.StatusOther
	TeamHome        = theme.Heading
	Score           = theme.Score
)

func StatusBadge(status string) string { return theme.StatusBadge(status) }