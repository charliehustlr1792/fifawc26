package theme

import "github.com/charmbracelet/lipgloss"

const (
	Red       = lipgloss.Color("#E60026")
	Blue      = lipgloss.Color("#0066B3")
	Green     = lipgloss.Color("#009A44")
	Cream     = lipgloss.Color("#F4E9D8")
	OffBlack  = lipgloss.Color("#0B0E14")
	White     = lipgloss.Color("#FFFFFF")
	MutedText = lipgloss.Color("245")
	DimText   = lipgloss.Color("240")
)

var (
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Cream)

	Heading = lipgloss.NewStyle().
		Bold(true).
		Foreground(Blue)

	Subtle = lipgloss.NewStyle().
		Foreground(MutedText)

	Dim = lipgloss.NewStyle().
		Foreground(DimText)

	StatusLive = lipgloss.NewStyle().
		Bold(true).
		Foreground(White).
		Background(Red).
		Padding(0, 1)

	StatusScheduled = lipgloss.NewStyle().
		Foreground(Blue)

	StatusFinished = lipgloss.NewStyle().
		Foreground(MutedText)

	StatusOther = lipgloss.NewStyle().
		Foreground(Green)

	TabActive = lipgloss.NewStyle().
		Bold(true).
		Foreground(White).
		Background(Red).
		Padding(0, 2)

	TabInactive = lipgloss.NewStyle().
		Foreground(MutedText).
		Padding(0, 2)

	Selection = lipgloss.NewStyle().
		Bold(true).
		Foreground(White).
		Background(Blue)

	Score = lipgloss.NewStyle().
		Bold(true).
		Foreground(Cream)

	Points = lipgloss.NewStyle().
		Bold(true).
		Foreground(Green)

	Error = lipgloss.NewStyle().
		Foreground(Red)

	Panel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Blue).
		Padding(0, 1)

	PanelAccent = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(Red).
		Padding(1, 2)
)

func StatusBadge(status string) string {
	switch status {
	case "IN_PLAY", "PAUSED", "EXTRA_TIME", "PENALTY_SHOOTOUT":
		return StatusLive.Render("LIVE")
	case "FINISHED":
		return StatusFinished.Render("FT")
	case "SCHEDULED", "TIMED":
		return StatusScheduled.Render("•")
	default:
		return StatusOther.Render(status)
	}
}