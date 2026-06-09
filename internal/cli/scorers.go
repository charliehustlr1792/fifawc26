package cli

import (
	"os"

	"github.com/charliehustlr1792/fifawc26/internal/render"
	"github.com/spf13/cobra"
)

var scorersLimit int

var scorersCmd = &cobra.Command{
	Use:   "scorers",
	Short: "Show top goalscorers leaderboard",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := deps.client.GetScorers(cmd.Context(), "WC", scorersLimit)
		if err != nil {
			return err
		}
		render.Scorers(os.Stdout, resp)
		return nil
	},
}

func init() {
	scorersCmd.Flags().IntVarP(&scorersLimit, "limit", "n", 10, "How many scorers to fetch")
	rootCmd.AddCommand(scorersCmd)
}