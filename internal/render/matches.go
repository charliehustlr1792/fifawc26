package render

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/charliehustlr1792/fifawc26/internal/api"
	"github.com/jedib0t/go-pretty/v6/table"
)

func Matches(w io.Writer, resp *api.MatchesResponse, teamFilter string) {
	matches := resp.Matches
	if teamFilter != "" {
		matches = filterByTeam(matches, teamFilter)
	}
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].UTCDate.Before(matches[j].UTCDate)
	})

	if len(matches) == 0 {
		fmt.Fprintln(w, Subtle.Render("No matches match the filters."))
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(w)
	t.AppendHeader(table.Row{"Date (local)", "Stage", "Status", "Home", "Score", "Away"})
	for _, m := range matches {
		local := m.UTCDate.Local().Format("Mon Jan 02 15:04")
		t.AppendRow(table.Row{
			local,
			shortStage(m.Stage, m.Group),
			StatusBadge(m.Status),
			TeamHome.Render(m.HomeTeam.Name),
			formatScore(m),
			m.AwayTeam.Name,
		})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()
}

func filterByTeam(in []api.Match, q string) []api.Match {
	out := make([]api.Match, 0, len(in))
	q = strings.ToLower(q)
	for _, m := range in {
		if strings.Contains(strings.ToLower(m.HomeTeam.Name), q) ||
			strings.Contains(strings.ToLower(m.AwayTeam.Name), q) ||
			strings.EqualFold(m.HomeTeam.TLA, q) ||
			strings.EqualFold(m.AwayTeam.TLA, q) {
			out = append(out, m)
		}
	}
	return out
}

func formatScore(m api.Match) string {
	h := m.Score.FullTime.Home
	a := m.Score.FullTime.Away
	if h == nil || a == nil {
		return Subtle.Render("- : -")
	}
	return Score.Render(fmt.Sprintf("%d : %d", *h, *a))
}

func shortStage(stage, group string) string {
	if group != "" {
		return group
	}
	switch stage {
	case "GROUP_STAGE":
		return "Group"
	case "LAST_16":
		return "R16"
	case "QUARTER_FINALS":
		return "QF"
	case "SEMI_FINALS":
		return "SF"
	case "THIRD_PLACE":
		return "3rd"
	case "FINAL":
		return "Final"
	default:
		return stage
	}
}