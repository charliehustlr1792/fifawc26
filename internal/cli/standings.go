package cli

import (
	"os"

	"github.com/charliehustlr1792/fifawc26/internal/render"
	"github.com/spf13/cobra"
)

var standingsCmd = &cobra.Command{
	Use:   "standings",
	Short: "Show group-stage standings",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := deps.client.GetStandings(cmd.Context(), "WC")
		if err != nil {
			return err
		}
		render.Standings(os.Stdout, resp)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(standingsCmd)
}