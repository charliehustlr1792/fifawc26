package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charliehustlr1792/fifawc26/internal/api"
	"github.com/charliehustlr1792/fifawc26/internal/render"
)

type tab int

const (
	tabStandings tab = iota
	tabMatches
	tabScorers
	tabCount
)

var tabNames = []string{"Standings", "Matches", "Scorers"}

type Model struct {
	client api.Provider
	active tab
	width  int
	height int

	standings *api.StandingsResponse
	matches   *api.MatchesResponse
	scorers   *api.ScorersResponse

	loading bool
	err     error
}

type standingsMsg struct {
	data *api.StandingsResponse
	err  error
}

type matchesMsg struct {
	data *api.MatchesResponse
	err  error
}

type scorersMsg struct {
	data *api.ScorersResponse
	err  error
}

func NewModel(client api.Provider) Model {
	return Model{client: client, active: tabStandings, loading: true}
}

func (m Model) Init() tea.Cmd {
	return m.fetchActive()
}

func (m Model) fetchActive() tea.Cmd {
	switch m.active {
	case tabStandings:
		return func() tea.Msg {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			d, err := m.client.GetStandings(ctx, "WC")
			return standingsMsg{d, err}
		}
	case tabMatches:
		return func() tea.Msg {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			d, err := m.client.GetMatches(ctx, "WC", api.MatchFilter{})
			return matchesMsg{d, err}
		}
	case tabScorers:
		return func() tea.Msg {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			d, err := m.client.GetScorers(ctx, "WC", 20)
			return scorersMsg{d, err}
		}
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab", "right", "l":
			m.active = (m.active + 1) % tabCount
			return m.maybeFetch()
		case "shift+tab", "left", "h":
			m.active = (m.active + tabCount - 1) % tabCount
			return m.maybeFetch()
		case "1":
			m.active = tabStandings
			return m.maybeFetch()
		case "2":
			m.active = tabMatches
			return m.maybeFetch()
		case "3":
			m.active = tabScorers
			return m.maybeFetch()
		case "r":
			m.loading = true
			m.err = nil
			return m, m.fetchActive()
		}

	case standingsMsg:
		m.loading = false
		m.standings = msg.data
		m.err = msg.err
	case matchesMsg:
		m.loading = false
		m.matches = msg.data
		m.err = msg.err
	case scorersMsg:
		m.loading = false
		m.scorers = msg.data
		m.err = msg.err
	}
	return m, nil
}

func (m Model) maybeFetch() (tea.Model, tea.Cmd) {
	if m.hasData() {
		m.loading = false
		m.err = nil
		return m, nil
	}
	m.loading = true
	m.err = nil
	return m, m.fetchActive()
}

func (m Model) hasData() bool {
	switch m.active {
	case tabStandings:
		return m.standings != nil
	case tabMatches:
		return m.matches != nil
	case tabScorers:
		return m.scorers != nil
	}
	return false
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString(m.renderHeader())
	b.WriteString("\n")
	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")

	switch {
	case m.loading:
		b.WriteString(render.Subtle.Render("Fetching data..."))
	case m.err != nil:
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4D4D")).Render("Error: " + m.err.Error()))
	default:
		b.WriteString(m.renderBody())
	}

	b.WriteString("\n\n")
	b.WriteString(m.renderFooter())
	return b.String()
}

func (m Model) renderHeader() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFD23F")).
		Render("⚽ FIFA World Cup 2026")
	return title
}

func (m Model) renderTabs() string {
	active := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#0B0E14")).
		Background(lipgloss.Color("#FFD23F")).
		Padding(0, 2)
	inactive := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Padding(0, 2)

	parts := make([]string, 0, len(tabNames))
	for i, name := range tabNames {
		label := fmt.Sprintf("%d %s", i+1, name)
		if tab(i) == m.active {
			parts = append(parts, active.Render(label))
		} else {
			parts = append(parts, inactive.Render(label))
		}
	}
	return strings.Join(parts, " ")
}

func (m Model) renderBody() string {
	var sb strings.Builder
	switch m.active {
	case tabStandings:
		if m.standings != nil {
			render.Standings(&sb, m.standings)
		}
	case tabMatches:
		if m.matches != nil {
			render.Matches(&sb, m.matches, "")
		}
	case tabScorers:
		if m.scorers != nil {
			render.Scorers(&sb, m.scorers)
		}
	}
	return sb.String()
}

func (m Model) renderFooter() string {
	return render.Subtle.Render("[1/2/3 or ←/→] tabs   [r] refresh   [q] quit")
}

func Run(client api.Provider) error {
	p := tea.NewProgram(NewModel(client), tea.WithAltScreen())
	_, err := p.Run()
	return err
}