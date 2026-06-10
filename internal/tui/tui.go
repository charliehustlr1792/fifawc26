package tui

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/charliehustlr1792/fifawc26/internal/api"
	"github.com/charliehustlr1792/fifawc26/internal/render"
	"github.com/charliehustlr1792/fifawc26/internal/theme"
)

type tab int

const (
	tabStandings tab = iota
	tabMatches
	tabScorers
	tabCount
)

type screen int

const (
	screenTabs screen = iota
	screenMatchDetail
	screenTeamDetail
)

const chromeHeight = 8

var tabNames = []string{"Standings", "Matches", "Scorers"}

type Model struct {
	client api.Provider
	active tab
	screen screen
	width  int
	height int
	vp     viewport.Model
	ready  bool

	groupFilter     string
	matchCursor     int
	standingsCursor int
	selectedMatch   *api.Match
	selectedTeam    *api.TeamDetail
	selectedTeamID  int
	teamFixtures    []api.Match

	standings *api.StandingsResponse
	matches   *api.MatchesResponse
	scorers   *api.ScorersResponse

	loading bool
	err     error

	lastUpdated time.Time
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
type teamMsg struct {
	data *api.TeamDetail
	err  error
}

type tickMsg time.Time

func NewModel(client api.Provider) Model {
	return Model{client: client, active: tabStandings, loading: true}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.fetchActive(), tickCmd(45*time.Second))
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

func (m Model) fetchMatches() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		d, err := m.client.GetMatches(ctx, "WC", api.MatchFilter{})
		return matchesMsg{d, err}
	}
}

func (m Model) fetchTeam(id int) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		t, err := m.client.GetTeam(ctx, id)
		return teamMsg{t, err}
	}
}

func tickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m Model) sortedMatches() []api.Match {
	if m.matches == nil {
		return nil
	}
	out := make([]api.Match, len(m.matches.Matches))
	copy(out, m.matches.Matches)
	sort.Slice(out, func(i, j int) bool { return out[i].UTCDate.Before(out[j].UTCDate) })
	return out
}

func filterMatchesByTeam(matches []api.Match, teamID int) []api.Match {
	out := make([]api.Match, 0, 8)
	for _, m := range matches {
		if m.HomeTeam.ID == teamID || m.AwayTeam.ID == teamID {
			out = append(out, m)
		}
	}
	return out
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		h := msg.Height - chromeHeight
		if h < 3 {
			h = 3
		}
		innerW := msg.Width - 6
		if innerW < 20 {
			innerW = 20
		}
		if !m.ready {
			m.vp = viewport.New(innerW, h)
			m.ready = true
		} else {
			m.vp.Width = innerW
			m.vp.Height = h
		}
		m.refreshViewportContent()
		return m, nil

	case tea.KeyMsg:
		if m.screen == screenMatchDetail {
			switch msg.String() {
			case "esc", "q":
				m.screen = screenTabs
				m.selectedMatch = nil
				m.refreshViewportContent()
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			case "t":
				if m.selectedMatch != nil {
					id := m.selectedMatch.HomeTeam.ID
					if id == 0 {
						break
					}
					m.selectedTeamID = id
					m.screen = screenTeamDetail
					m.loading = true
					m.refreshViewportContent()
					var cmds []tea.Cmd
					cmds = append(cmds, m.fetchTeam(id))
					if m.matches == nil {
						cmds = append(cmds, m.fetchMatches())
					}
					return m, tea.Batch(cmds...)
				}
			}
			var cmd tea.Cmd
			m.vp, cmd = m.vp.Update(msg)
			return m, cmd
		}

		if m.screen == screenTeamDetail {
			switch msg.String() {
			case "esc", "q":
				m.screen = screenTabs
				m.selectedTeam = nil
				m.refreshViewportContent()
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			}
			var cmd tea.Cmd
			m.vp, cmd = m.vp.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.active = (m.active + 1) % tabCount
			m.resetTabState()
			return m.maybeFetch()
		case "shift+tab":
			m.active = (m.active + tabCount - 1) % tabCount
			m.resetTabState()
			return m.maybeFetch()
		case "1":
			m.active = tabStandings
			m.resetTabState()
			return m.maybeFetch()
		case "2":
			m.active = tabMatches
			m.resetTabState()
			return m.maybeFetch()
		case "3":
			m.active = tabScorers
			m.resetTabState()
			return m.maybeFetch()
		case "r":
			m.loading = true
			m.err = nil
			m.refreshViewportContent()
			return m, m.fetchActive()
		}

		if m.active == tabStandings {
			s := msg.String()
			if len(s) == 1 && s[0] >= 'a' && s[0] <= 'l' {
				m.groupFilter = strings.ToUpper(s)
				m.standingsCursor = 0
				m.refreshViewportContent()
				m.vp.GotoTop()
				return m, nil
			}
			if len(s) == 1 && s[0] >= 'A' && s[0] <= 'L' {
				m.groupFilter = s
				m.standingsCursor = 0
				m.refreshViewportContent()
				m.vp.GotoTop()
				return m, nil
			}
			if s == "0" || s == "esc" {
				m.groupFilter = ""
				m.standingsCursor = 0
				m.refreshViewportContent()
				m.vp.GotoTop()
				return m, nil
			}
			if m.standings != nil {
				rows := render.VisibleStandingsRows(m.standings, m.groupFilter)
				switch s {
				case "up", "k":
					if m.standingsCursor > 0 {
						m.standingsCursor--
					}
					m.refreshViewportContent()
					return m, nil
				case "down", "j":
					if m.standingsCursor < len(rows)-1 {
						m.standingsCursor++
					}
					m.refreshViewportContent()
					return m, nil
				case "enter":
					if len(rows) > 0 && m.standingsCursor < len(rows) {
						id := rows[m.standingsCursor].Team.ID
						if id == 0 {
							break
						}
						m.selectedTeamID = id
						m.screen = screenTeamDetail
						m.loading = true
						m.refreshViewportContent()
						var cmds []tea.Cmd
						cmds = append(cmds, m.fetchTeam(id))
						if m.matches == nil {
							cmds = append(cmds, m.fetchMatches())
						}
						return m, tea.Batch(cmds...)
					}
				}
			}
		}

		if m.active == tabMatches && m.matches != nil {
			list := m.sortedMatches()
			switch msg.String() {
			case "up", "k":
				if m.matchCursor > 0 {
					m.matchCursor--
				}
				m.refreshViewportContent()
				return m, nil
			case "down", "j":
				if m.matchCursor < len(list)-1 {
					m.matchCursor++
				}
				m.refreshViewportContent()
				return m, nil
			case "enter":
				if m.matchCursor >= 0 && m.matchCursor < len(list) {
					sel := list[m.matchCursor]
					m.selectedMatch = &sel
					m.screen = screenMatchDetail
					m.refreshViewportContent()
					m.vp.GotoTop()
					return m, nil
				}
			}
		}

		var cmd tea.Cmd
		m.vp, cmd = m.vp.Update(msg)
		return m, cmd

	case standingsMsg:
		m.loading = false
		m.standings = msg.data
		if msg.err == nil {
			m.lastUpdated = time.Now()
		}
		m.err = msg.err
		m.refreshViewportContent()
		m.vp.GotoTop()
		return m, nil

	case matchesMsg:
		m.loading = false
		m.matches = msg.data
		m.err = msg.err
		if msg.err == nil {
			m.lastUpdated = time.Now()
		}
		if m.screen != screenMatchDetail && m.screen != screenTeamDetail {
			m.matchCursor = 0
		}
		if m.screen == screenMatchDetail && m.selectedMatch != nil && msg.data != nil {
			for _, fresh := range msg.data.Matches {
				if fresh.ID == m.selectedMatch.ID {
					freshCopy := fresh
					m.selectedMatch = &freshCopy
					break
				}
			}
		}
		if m.screen == screenTeamDetail && m.selectedTeam != nil && msg.data != nil {
			m.teamFixtures = filterMatchesByTeam(msg.data.Matches, m.selectedTeam.ID)
			m.refreshViewportContent()
		} else if m.screen != screenMatchDetail && m.screen != screenTeamDetail {
			m.refreshViewportContent()
			m.vp.GotoTop()
		}
		return m, nil

	case scorersMsg:
		m.loading = false
		m.scorers = msg.data
		if msg.err == nil {
			m.lastUpdated = time.Now()
		}
		m.err = msg.err
		m.refreshViewportContent()
		m.vp.GotoTop()
		return m, nil

	case teamMsg:
		m.loading = false
		if msg.err != nil {
			m.err = fmt.Errorf("team data unavailable (keyless mode has limited coverage)")
		} else {
			m.selectedTeam = msg.data
			m.err = nil
			if msg.data != nil && m.matches != nil {
				m.teamFixtures = filterMatchesByTeam(m.matches.Matches, msg.data.ID)
			}
		}
		m.refreshViewportContent()
		m.vp.GotoTop()
		return m, nil

	case tickMsg:
		cmds := []tea.Cmd{tickCmd(45 * time.Second)}
		if !m.loading {
			cmds = append(cmds, m.fetchActive())
		}
		return m, tea.Batch(cmds...)
	}

	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m *Model) resetTabState() {
	m.groupFilter = ""
	m.matchCursor = 0
	m.standingsCursor = 0
	m.screen = screenTabs
	m.selectedMatch = nil
	m.selectedTeam = nil
}

func (m Model) maybeFetch() (tea.Model, tea.Cmd) {
	if m.hasData() {
		m.loading = false
		m.err = nil
		m.refreshViewportContent()
		m.vp.GotoTop()
		return m, nil
	}
	m.loading = true
	m.err = nil
	m.refreshViewportContent()
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

func (m *Model) refreshViewportContent() {
	if !m.ready {
		return
	}
	switch {
	case m.screen == screenTeamDetail && m.selectedTeam != nil:
		m.vp.SetContent(render.TeamDetail(m.selectedTeam, m.teamFixtures, m.standings))
		return
	case m.screen == screenMatchDetail && m.selectedMatch != nil:
		m.vp.SetContent(render.MatchDetail(*m.selectedMatch))
		return
	case m.loading:
		m.vp.SetContent(render.Subtle.Render("Fetching data..."))
	case m.err != nil:
		m.vp.SetContent(theme.Error.Render("Error: " + m.err.Error()))
	default:
		m.vp.SetContent(m.renderBody())
	}
}

func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}
	var b strings.Builder
	b.WriteString(m.renderHeader())
	b.WriteString("\n")
	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")

	bodyW := m.width - 2
	if bodyW < 20 {
		bodyW = 20
	}
	body := theme.Panel.Width(bodyW).Render(m.vp.View())
	b.WriteString(body)

	b.WriteString("\n")
	b.WriteString(m.renderFooter())
	return b.String()
}

func (m Model) renderHeader() string {
	return theme.Title.Render("FIFA World Cup 2026 ⚽")
}

func (m Model) renderTabs() string {
	parts := make([]string, 0, len(tabNames))
	for i, name := range tabNames {
		label := fmt.Sprintf("%d %s", i+1, name)
		if tab(i) == m.active {
			parts = append(parts, theme.TabActive.Render(label))
		} else {
			parts = append(parts, theme.TabInactive.Render(label))
		}
	}
	return strings.Join(parts, " ")
}

func (m Model) renderBody() string {
	var sb strings.Builder
	switch m.active {
	case tabStandings:
		if m.standings != nil {
			sb.WriteString(render.StandingsListTUI(m.standings, m.groupFilter, m.standingsCursor))
		}
	case tabMatches:
		if m.matches != nil {
			sb.WriteString(render.MatchesList(m.sortedMatches(), m.matchCursor))
		}
	case tabScorers:
		if m.scorers != nil {
			render.Scorers(&sb, m.scorers)
		}
	}
	return sb.String()
}

func (m Model) renderFooter() string {
	var help string
	switch {
	case m.screen == screenTeamDetail:
		help = "[esc] back   [↑/↓] scroll   [q] quit"
	case m.screen == screenMatchDetail:
		help = "[esc] back   [↑/↓] scroll   [t] home team   [q] quit"
	case m.active == tabStandings:
		help = "[1/2/3] tabs   [A–L] filter group   [0] all   [↑/↓] team   [enter] detail   [r] refresh   [q] quit"
	case m.active == tabMatches:
		help = "[1/2/3] tabs   [↑/↓ or j/k] select   [enter] details   [r] refresh   [q] quit"
	default:
		help = "[1/2/3] tabs   [↑/↓ pgup/pgdn] scroll   [r] refresh   [q] quit"
	}

	right := ""
	if !m.lastUpdated.IsZero() {
		ago := time.Since(m.lastUpdated).Truncate(time.Second)
		right = render.Subtle.Render(fmt.Sprintf("updated %s ago", ago))
	}
	if m.loading {
		right = render.Subtle.Render("refreshing...")
	}

	leftR := render.Subtle.Render(help)
	if right == "" || m.width == 0 {
		return leftR
	}
	gap := m.width - lipgloss.Width(leftR) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}
	return leftR + strings.Repeat(" ", gap) + right
}

func Run(client api.Provider) error {
	p := tea.NewProgram(NewModel(client), tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
