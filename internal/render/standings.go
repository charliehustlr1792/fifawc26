package render

import (
	"io"
	"sort"
	"strings"

	"github.com/charliehustlr1792/fifawc26/internal/api"
	"github.com/jedib0t/go-pretty/v6/table"
)

func Standings(w io.Writer, resp *api.StandingsResponse) {
	StandingsFiltered(w, resp, "")
}

func StandingsFiltered(w io.Writer, resp *api.StandingsResponse, groupLetter string) {
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
		w.Write([]byte(Subtle.Render("No groups match the filter.")))
		return
	}

	for _, g := range groups {
		t := table.NewWriter()
		t.SetOutputMirror(w)
		t.SetTitle(prettyGroup(g.Group))
		t.AppendHeader(table.Row{"#", "Team", "P", "W", "D", "L", "GF", "GA", "GD", "Pts", "Form"})
		for _, row := range g.Table {
			t.AppendRow(table.Row{
				row.Position,
				row.Team.Name,
				row.PlayedGames,
				row.Won,
				row.Draw,
				row.Lost,
				row.GoalsFor,
				row.GoalsAgainst,
				row.GoalDifference,
				row.Points,
				row.Form,
			})
		}
		t.SetStyle(table.StyleRounded)
		t.Render()
	}
}

func prettyGroup(g string) string {
	if strings.HasPrefix(g, "GROUP_") {
		return "Group " + strings.TrimPrefix(g, "GROUP_")
	}
	return g
}