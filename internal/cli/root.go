package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/charliehustlr1792/fifawc26/internal/api"
	"github.com/charliehustlr1792/fifawc26/internal/cache"
	"github.com/charliehustlr1792/fifawc26/internal/config"
	"github.com/spf13/cobra"
	"github.com/charliehustlr1792/fifawc26/internal/tui"
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
		if cmd.Name() == "setup" {
    return nil
}
    cfg, err := config.Load()
    if err != nil {
        return err
    }

    if cfg.NeedsOnboarding() {
        if err := tui.RunOnboarding(cfg); err != nil {
            return err
        }
        cfg, err = config.Load()
        if err != nil {
            return err
        }
    }

    bc, err := cache.NewBoltCache(cfg.CacheDir)
    if err != nil {
        return err
    }

    var provider api.Provider
    if cfg.Tier == config.TierKeyed && cfg.APIKey != "" {
        provider = api.NewCachedProvider(
            api.NewFootballDataClient(cfg.APIKey),
            bc,
            api.DefaultTTLs(),
        )
    } else {
        provider = api.NewCachedProvider(
            api.NewFootballDataClient(""),
            bc,
            api.DefaultTTLs(),
        )
    }

    deps = &appDeps{cfg: cfg, cache: bc, client: provider}
    return nil
},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if deps != nil && deps.cache != nil {
			_ = deps.cache.Close()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
    if err := tui.Run(deps.client); err != nil {
        fmt.Fprintln(os.Stderr, "tui error:", err)
        os.Exit(1)
    }
},
}

func Execute() {
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}