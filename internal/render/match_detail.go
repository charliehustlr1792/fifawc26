package render

import (
	"fmt"
	"strings"

	"github.com/charliehustlr1792/fifawc26/internal/api"
	"github.com/charmbracelet/lipgloss"
)

func MatchDetail(m api.Match) string {
	var b strings.Builder

	heading := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD23F")).
		Render(fmt.Sprintf("%s  vs  %s", m.HomeTeam.Name, m.AwayTeam.Name))
	b.WriteString(heading)
	b.WriteString("\n\n")

	fmt.Fprintf(&b, "Stage:      %s\n", shortStage(m.Stage, m.Group))
	if m.Matchday > 0 {
		fmt.Fprintf(&b, "Matchday:   %d\n", m.Matchday)
	}
	fmt.Fprintf(&b, "Kickoff:    %s\n", m.UTCDate.Local().Format("Mon, 02 Jan 2006 15:04 MST"))
	fmt.Fprintf(&b, "Status:     %s\n", StatusBadge(m.Status))

	fmt.Fprintf(&b, "Score:      %s\n", scoreLine(m))
	if m.Score.HalfTime.Home != nil && m.Score.HalfTime.Away != nil {
		fmt.Fprintf(&b, "Half-time:  %d : %d\n", *m.Score.HalfTime.Home, *m.Score.HalfTime.Away)
	}
	if m.Score.Winner != "" {
		fmt.Fprintf(&b, "Winner:     %s\n", m.Score.Winner)
	}
	if m.Score.Duration != "" {
		fmt.Fprintf(&b, "Duration:   %s\n", m.Score.Duration)
	}

	b.WriteString("\n")
	b.WriteString(Subtle.Render("[esc] back to list"))
	return b.String()
}

func scoreLine(m api.Match) string {
	h := m.Score.FullTime.Home
	a := m.Score.FullTime.Away
	if h == nil || a == nil {
		return Subtle.Render("not played yet")
	}
	return Score.Render(fmt.Sprintf("%d : %d", *h, *a))
}