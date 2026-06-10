package render

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charliehustlr1792/fifawc26/internal/api"
	"github.com/charmbracelet/lipgloss"
	"github.com/charliehustlr1792/fifawc26/internal/theme"
)

func MatchesList(matches []api.Match, selected int) string {
	if len(matches) == 0 {
		return Subtle.Render("No matches.")
	}
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].UTCDate.Before(matches[j].UTCDate)
	})

	rowStyle := lipgloss.NewStyle()
	selStyle := theme.Selection

	var b strings.Builder
	for i, m := range matches {
		when := m.UTCDate.Local().Format("Mon Jan 02 15:04")
		stage := shortStage(m.Stage, m.Group)
		score := formatScore(m)
		home := m.HomeTeam.Name
		away := m.AwayTeam.Name
		if home == "" {
			home = "TBD"
		}
		if away == "" {
			away = "TBD"
		}
		line := fmt.Sprintf("%-18s  %-8s  %s  %-22s %s  %s",
			when, stage, StatusBadge(m.Status), home, score, away)
		if i == selected {
			b.WriteString(selStyle.Render("▶ " + line))
		} else {
			b.WriteString(rowStyle.Render("  " + line))
		}
		b.WriteString("\n")
	}
	return b.String()
}