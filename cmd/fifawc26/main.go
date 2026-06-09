package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charliehustlr1792/fifawc26/internal/api"
	"github.com/charliehustlr1792/fifawc26/internal/cache"
	"github.com/charliehustlr1792/fifawc26/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "config error:", err)
		os.Exit(1)
	}

	bc, err := cache.NewBoltCache(cfg.CacheDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "cache error:", err)
		os.Exit(1)
	}
	defer bc.Close()

	client := api.NewCachedProvider(
		api.NewFootballDataClient(cfg.APIKey),
		bc,
		api.DefaultTTLs(),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	standings, err := client.GetStandings(ctx, "WC")
	if err != nil {
		fmt.Fprintln(os.Stderr, "standings:", err)
	} else {
		fmt.Printf("Standings groups: %d\n", len(standings.Standings))
		for _, s := range standings.Standings {
			fmt.Printf("  %s (%s) — %d teams\n", s.Group, s.Type, len(s.Table))
		}
	}

	matches, err := client.GetMatches(ctx, "WC", api.MatchFilter{Status: "SCHEDULED"})
	if err != nil {
		fmt.Fprintln(os.Stderr, "matches:", err)
	} else {
		fmt.Printf("\nScheduled matches: %d\n", matches.ResultSet.Count)
		for i, m := range matches.Matches {
			if i >= 3 {
				break
			}
			fmt.Printf("  %s — %s vs %s (%s)\n", m.UTCDate.Format("Jan 02 15:04"), m.HomeTeam.Name, m.AwayTeam.Name, m.Stage)
		}
	}

	scorers, err := client.GetScorers(ctx, "WC", 10)
	if err != nil {
		fmt.Fprintln(os.Stderr, "scorers:", err)
	} else {
		fmt.Printf("\nTop scorers: %d\n", scorers.Count)
		for i, s := range scorers.Scorers {
			if i >= 5 {
				break
			}
			fmt.Printf("  %s (%s) — %d goals\n", s.Player.Name, s.Team.TLA, s.Goals)
		}
	}
}