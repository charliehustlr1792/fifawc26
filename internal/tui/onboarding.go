package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charliehustlr1792/fifawc26/internal/config"
	"github.com/charliehustlr1792/fifawc26/internal/theme"
)

type onboardScreen int

const (
	obChoose onboardScreen = iota
	obKeyInput
	obDone
)

type OnboardModel struct {
	cfg      *config.Config
	screen   onboardScreen
	cursor   int
	input    textinput.Model
	err      string
	finished bool
}

func NewOnboarding(cfg *config.Config) OnboardModel {
	ti := textinput.New()
	ti.Placeholder = "paste your football-data.org token here"
	ti.CharLimit = 100
	ti.Width = 50

	return OnboardModel{
		cfg:    cfg,
		screen: obChoose,
		input:  ti,
	}
}

func (m OnboardModel) Init() tea.Cmd { return nil }

func (m OnboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.screen {
		case obChoose:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < 1 {
					m.cursor++
				}
			case "1":
				m.cursor = 0
			case "2":
				m.cursor = 1
			case "enter":
				if m.cursor == 0 {
					m.screen = obKeyInput
					m.input.Focus()
					return m, textinput.Blink
				}
				m.cfg.Tier = config.TierKeyless
				m.cfg.APIKey = ""
				if err := config.Save(m.cfg); err != nil {
					m.err = err.Error()
					return m, nil
				}
				m.finished = true
				return m, tea.Quit
			}

		case obKeyInput:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.screen = obChoose
				m.err = ""
				m.input.SetValue("")
				return m, nil
			case "enter":
				key := strings.TrimSpace(m.input.Value())
				if len(key) < 20 {
					m.err = "that doesn't look like a valid token (expected ~36 chars)"
					return m, nil
				}
				m.cfg.Tier = config.TierKeyed
				m.cfg.APIKey = key
				if err := config.Save(m.cfg); err != nil {
					m.err = err.Error()
					return m, nil
				}
				m.finished = true
				return m, tea.Quit
			}
			var cmd tea.Cmd
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m OnboardModel) View() string {
	switch m.screen {
	case obChoose:
		return m.viewChoose()
	case obKeyInput:
		return m.viewKeyInput()
	}
	return ""
}

func (m OnboardModel) Finished() bool { return m.finished }

var (
	obTitle      = theme.Title
	obDim        = theme.Subtle
	obSelected   = theme.Selection.Padding(0, 1)
	obUnselected = lipgloss.NewStyle().Padding(0, 1)
	obPanel      = theme.PanelAccent
	obErr        = theme.Error
)

func (m OnboardModel) viewChoose() string {
	var b strings.Builder
	b.WriteString(obTitle.Render("⚽ Welcome to fifawc26"))
	b.WriteString("\n\n")
	b.WriteString("Choose how you'd like to fetch World Cup data:\n\n")

	opts := []struct {
		title string
		desc  string
	}{
		{
			"1. With API key  (recommended — full data)",
			"   Free key from football-data.org. Standings, fixtures,\n   live scores, top scorers — everything.",
		},
		{
			"2. Keyless mode  (no setup — live scores limited)",
			"   Uses public sources. Works out of the box, but live\n   scores update slower and some data may be missing.",
		},
	}
	for i, o := range opts {
		if i == m.cursor {
			b.WriteString(obSelected.Render("▶ " + o.title))
		} else {
			b.WriteString(obUnselected.Render("  " + o.title))
		}
		b.WriteString("\n")
		b.WriteString(obDim.Render(o.desc))
		b.WriteString("\n\n")
	}

	b.WriteString(obDim.Render("[↑/↓] move   [enter] select   [q] quit"))
	if m.err != "" {
		b.WriteString("\n\n")
		b.WriteString(obErr.Render("error: " + m.err))
	}

	return obPanel.Render(b.String())
}

func (m OnboardModel) viewKeyInput() string {
	var b strings.Builder
	b.WriteString(obTitle.Render("🔑 Get your free API key"))
	b.WriteString("\n\n")
	b.WriteString("How to get one (takes 30 seconds):\n\n")
	b.WriteString(fmt.Sprintf("  %s open  %s\n", obDim.Render("1."), "https://www.football-data.org/client/register"))
	b.WriteString(fmt.Sprintf("  %s enter your email and tick the boxes\n", obDim.Render("2.")))
	b.WriteString(fmt.Sprintf("  %s check your inbox — they email the token\n", obDim.Render("3.")))
	b.WriteString(fmt.Sprintf("  %s paste it below\n", obDim.Render("4.")))
	b.WriteString("\n")
	b.WriteString("Token:\n")
	b.WriteString(m.input.View())
	b.WriteString("\n\n")
	b.WriteString(obDim.Render("[enter] save   [esc] back   [ctrl+c] quit"))
	if m.err != "" {
		b.WriteString("\n\n")
		b.WriteString(obErr.Render("⚠ " + m.err))
	}

	return obPanel.Render(b.String())
}

func RunOnboarding(cfg *config.Config) error {
	m := NewOnboarding(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return err
	}
	if r, ok := result.(OnboardModel); ok && !r.finished {
		return fmt.Errorf("onboarding cancelled")
	}
	return nil
}