package render

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charliehustlr1792/fifawc26/internal/api"
	"github.com/charliehustlr1792/fifawc26/internal/theme"
)

func StandingsListTUI(resp *api.StandingsResponse, groupLetter string, cursor int) string {
	groups := make([]api.Standing, 0, len(resp.Standings))
	for _, s := range resp.Standings {
		if s.Type != "TOTAL" {
			continue
		}
		if groupLetter != "" && !strings.EqualFold(s.Group, "Group "+groupLetter) {
			continue
		}
		groups = append(groups, s)
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].Group < groups[j].Group })

	if len(groups) == 0 {
		return Subtle.Render("No groups match the filter.")
	}

	var b strings.Builder
	rowIdx := 0
	for gi, g := range groups {
		if gi > 0 {
			b.WriteString("\n")
		}
		b.WriteString(theme.Heading.Render(prettyGroup(g.Group)))
		b.WriteString("\n")
		header := fmt.Sprintf("  %-3s %-22s %3s %3s %3s %3s %4s %4s %4s %4s",
			"#", "Team", "P", "W", "D", "L", "GF", "GA", "GD", "Pts")
		b.WriteString(theme.Subtle.Render(header))
		b.WriteString("\n")
		for _, row := range g.Table {
			name := row.Team.Name
			if len(name) > 22 {
				name = name[:21] + "…"
			}
			line := fmt.Sprintf("%-3d %-22s %3d %3d %3d %3d %4d %4d %4d %4d",
				row.Position, name,
				row.PlayedGames, row.Won, row.Draw, row.Lost,
				row.GoalsFor, row.GoalsAgainst, row.GoalDifference, row.Points)
			if rowIdx == cursor {
				b.WriteString(theme.Selection.Render("▶ " + line))
			} else {
				b.WriteString("  " + line)
			}
			b.WriteString("\n")
			rowIdx++
		}
	}
	return b.String()
}

func VisibleStandingsRows(resp *api.StandingsResponse, groupLetter string) []api.StandingRow {
	groups := make([]api.Standing, 0, len(resp.Standings))
	for _, s := range resp.Standings {
		if s.Type != "TOTAL" {
			continue
		}
		if groupLetter != "" && !strings.EqualFold(s.Group, "Group "+groupLetter) {
			continue
		}
		groups = append(groups, s)
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].Group < groups[j].Group })

	var out []api.StandingRow
	for _, g := range groups {
		out = append(out, g.Table...)
	}
	return out
}
