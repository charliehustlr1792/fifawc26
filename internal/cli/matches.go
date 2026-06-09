package cli

import (
	"os"

	"github.com/charliehustlr1792/fifawc26/internal/api"
	"github.com/charliehustlr1792/fifawc26/internal/render"
	"github.com/spf13/cobra"
)

var (
	matchesStatus   string
	matchesMatchday int
	matchesTeam     string
	matchesDateFrom string
	matchesDateTo   string
)

var matchesCmd = &cobra.Command{
	Use:   "matches",
	Short: "Show World Cup matches (upcoming, live, or finished)",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := deps.client.GetMatches(cmd.Context(), "WC", api.MatchFilter{
			Status:   matchesStatus,
			Matchday: matchesMatchday,
			DateFrom: matchesDateFrom,
			DateTo:   matchesDateTo,
		})
		if err != nil {
			return err
		}
		render.Matches(os.Stdout, resp, matchesTeam)
		return nil
	},
}

func init() {
	matchesCmd.Flags().StringVarP(&matchesStatus, "status", "s", "", "Filter by status: SCHEDULED, IN_PLAY, FINISHED, etc.")
	matchesCmd.Flags().IntVarP(&matchesMatchday, "matchday", "m", 0, "Filter by matchday number")
	matchesCmd.Flags().StringVarP(&matchesTeam, "team", "t", "", "Filter by team name or TLA (e.g. Brazil, BRA)")
	matchesCmd.Flags().StringVar(&matchesDateFrom, "from", "", "Filter from date (yyyy-MM-dd)")
	matchesCmd.Flags().StringVar(&matchesDateTo, "to", "", "Filter to date (yyyy-MM-dd)")
	rootCmd.AddCommand(matchesCmd)
}