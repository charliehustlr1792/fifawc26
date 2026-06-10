package render

import (
	"fmt"
	"strings"
	"time"

	"github.com/charliehustlr1792/fifawc26/internal/api"
	"github.com/charliehustlr1792/fifawc26/internal/theme"
	"github.com/jedib0t/go-pretty/v6/table"
)

func TeamDetail(t *api.TeamDetail, matches []api.Match, standings *api.StandingsResponse) string {
	var b strings.Builder

	b.WriteString(theme.Title.Render(t.Name))
	b.WriteString("\n")

	if t.Area.Name != "" || t.Founded > 0 {
		sub := t.Area.Name
		if t.Founded > 0 {
			if sub != "" {
				sub += fmt.Sprintf(" · founded %d", t.Founded)
			} else {
				sub = fmt.Sprintf("founded %d", t.Founded)
			}
		}
		b.WriteString(theme.Subtle.Render(sub))
		b.WriteString("\n")
	}

	if standings != nil {
		if grp, pos := teamGroupPosition(t, standings); grp != "" {
			b.WriteString(theme.Heading.Render(fmt.Sprintf("%s — %s", grp, ordinal(pos))))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")

	if len(matches) > 0 {
		b.WriteString(theme.Heading.Render("Fixtures"))
		b.WriteString("\n")
		tw := table.NewWriter()
		tw.SetOutputMirror(&b)
		tw.AppendHeader(table.Row{"Date", "Stage", "Status", "Opponent", "Score"})
		for _, m := range matches {
			opp := m.AwayTeam.Name
			if m.AwayTeam.ID == t.ID || m.AwayTeam.Name == t.Name {
				opp = m.HomeTeam.Name
			}
			if opp == "" {
				opp = "TBD"
			}
			tw.AppendRow(table.Row{
				m.UTCDate.Local().Format("Mon Jan 02 15:04"),
				shortStage(m.Stage, m.Group),
				StatusBadge(m.Status),
				opp,
				formatScore(m),
			})
		}
		tw.SetStyle(table.StyleRounded)
		tw.Render()
		b.WriteString("\n")
	}

	if len(t.Squad) > 0 {
		b.WriteString(theme.Heading.Render("Squad"))
		b.WriteString("\n")
		sw := table.NewWriter()
		sw.SetOutputMirror(&b)
		sw.AppendHeader(table.Row{"#", "Name", "Position", "DOB", "Nationality"})
		for i, p := range t.Squad {
			dob := p.DateOfBirth
			if dob == "" {
				dob = "—"
			} else if parsed, err := time.Parse("2006-01-02", dob); err == nil {
				dob = parsed.Format("02 Jan 2006")
			}
			sw.AppendRow(table.Row{i + 1, p.Name, p.Position, dob, p.Nationality})
		}
		sw.SetStyle(table.StyleRounded)
		sw.Render()
	}

	return b.String()
}

func teamGroupPosition(t *api.TeamDetail, standings *api.StandingsResponse) (group string, pos int) {
	for _, s := range standings.Standings {
		if s.Type != "TOTAL" {
			continue
		}
		for _, row := range s.Table {
			if (t.ID != 0 && row.Team.ID == t.ID) || (t.Name != "" && row.Team.Name == t.Name) {
				return prettyGroup(s.Group), row.Position
			}
		}
	}
	return "", 0
}

func ordinal(n int) string {
	switch n {
	case 1:
		return "1st"
	case 2:
		return "2nd"
	case 3:
		return "3rd"
	default:
		return fmt.Sprintf("%dth", n)
	}
}
