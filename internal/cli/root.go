package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/charliehustlr1792/fifawc26/internal/api"
	"github.com/charliehustlr1792/fifawc26/internal/cache"
	"github.com/charliehustlr1792/fifawc26/internal/config"
	"github.com/spf13/cobra"
)

type appDeps struct {
	cfg    *config.Config
	client api.Provider
	cache  cache.Cache
}

var deps *appDeps

var rootCmd = &cobra.Command{
	Use:   "fifawc26",
	Short: "FIFA World Cup 26 in your terminal",
	Long:  "fifawc26 is a TUI for live World Cup 26 scores, fixtures, standings, and stats.",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		bc, err := cache.NewBoltCache(cfg.CacheDir)
		if err != nil {
			return err
		}
		deps = &appDeps{
			cfg:    cfg,
			cache:  bc,
			client: api.NewCachedProvider(api.NewFootballDataClient(cfg.APIKey), bc, api.DefaultTTLs()),
		}
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if deps != nil && deps.cache != nil {
			_ = deps.cache.Close()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TUI mode coming in a later step. Try: fifawc26 standings")
	},
}

func Execute() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}