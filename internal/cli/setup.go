package cli

import (
	"github.com/charliehustlr1792/fifawc26/internal/config"
	"github.com/charliehustlr1792/fifawc26/internal/tui"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Reconfigure fifawc26 (choose tier, change API key)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		cfg.Tier = ""
		cfg.APIKey = ""
		return tui.RunOnboarding(cfg)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}