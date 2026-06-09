package render

import (
	"fmt"
	"io"

	"github.com/charliehustlr1792/fifawc26/internal/api"
	"github.com/jedib0t/go-pretty/v6/table"
)

func Scorers(w io.Writer, resp *api.ScorersResponse) {
	if len(resp.Scorers) == 0 {
		fmt.Fprintln(w, Subtle.Render("No scorers yet. Tournament hasn't started or no goals scored."))
		return
	}
	t := table.NewWriter()
	t.SetOutputMirror(w)
	t.AppendHeader(table.Row{"#", "Player", "Team", "Goals", "Assists", "Pen"})
	for i, s := range resp.Scorers {
		t.AppendRow(table.Row{
			i + 1,
			s.Player.Name,
			s.Team.TLA,
			s.Goals,
			s.Assists,
			s.Penalties,
		})
	}
	t.SetStyle(table.StyleRounded)
	t.Render()
}