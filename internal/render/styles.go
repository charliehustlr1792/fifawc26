package render

import "github.com/charliehustlr1792/fifawc26/internal/theme"

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