package render

import "github.com/charmbracelet/lipgloss"

var (
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD23F")).
		MarginTop(1).
		MarginBottom(1)

	Subtle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))

	StatusLive = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF4D4D"))

	StatusScheduled = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4DD0E1"))

	StatusFinished = lipgloss.NewStyle().
		Foreground(lipgloss.Color("245"))

	StatusOther = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFB454"))

	TeamHome = lipgloss.NewStyle().Bold(true)
	Score    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF"))
)

func StatusBadge(status string) string {
	switch status {
	case "IN_PLAY", "PAUSED", "EXTRA_TIME", "PENALTY_SHOOTOUT":
		return StatusLive.Render(" LIVE ")
	case "FINISHED":
		return StatusFinished.Render("FT")
	case "SCHEDULED", "TIMED":
		return StatusScheduled.Render("•")
	default:
		return StatusOther.Render(status)
	}
}