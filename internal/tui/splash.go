package tui

import (
	"math/rand"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charliehustlr1792/fifawc26/internal/theme"
)

const splashDuration = 1500 * time.Millisecond

type splashDoneMsg struct{}

type SplashModel struct {
	width    int
	height   int
	tip      string
	finished bool
}

var splashTips = []string{
	"press 1 / 2 / 3 to switch tabs",
	"hit A-L to filter standings by group",
	"use enter on a match for the full recap",
	"data refreshes every 45 seconds automatically",
	"run fifawc26 setup to change your API key",
	"press r anywhere to force a refresh",
}

func NewSplash() SplashModel {
	return SplashModel{
		tip: splashTips[rand.Intn(len(splashTips))],
	}
}

func (m SplashModel) Init() tea.Cmd {
	return tea.Tick(splashDuration, func(time.Time) tea.Msg { return splashDoneMsg{} })
}

func (m SplashModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
		m.finished = true
		return m, tea.Quit
	case splashDoneMsg:
		m.finished = true
		return m, tea.Quit
	}
	return m, nil
}

func (m SplashModel) Finished() bool { return m.finished }

func (m SplashModel) View() string {
	wordmark := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Red).
		Render("F I F A   W C   2 0 2 6")

	tagline := lipgloss.NewStyle().
		Foreground(theme.Blue).
		Render("Elite Ball Knowledge")

	divider := lipgloss.NewStyle().
		Foreground(theme.Green).
		Render(strings.Repeat("─", 28))

	tip := theme.Subtle.Render("tip: " + m.tip)

	stack := lipgloss.JoinVertical(lipgloss.Center,
		wordmark,
		"",
		tagline,
		divider,
		"",
		tip,
	)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Blue).
		Padding(1, 6).
		Render(stack)

	if m.width == 0 || m.height == 0 {
		return box
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box)
}

func RunSplash() error {
	p := tea.NewProgram(NewSplash())
	_, err := p.Run()
	return err
}