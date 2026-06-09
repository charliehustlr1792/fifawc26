package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/charliehustlr1792/fifawc26/internal/api"
	"github.com/charliehustlr1792/fifawc26/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "config error:", err)
		os.Exit(1)
	}

	client := api.NewFootballDataClient(cfg.APIKey)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	comp, err := client.GetCompetition(ctx, "WC")
	if err != nil {
		fmt.Fprintln(os.Stderr, "api error:", err)
		os.Exit(1)
	}

	fmt.Printf("Competition: %s (%s)\n", comp.Name, comp.Code)
	fmt.Printf("Area: %s\n", comp.Area.Name)
	fmt.Printf("Current season: %s → %s\n", comp.CurrentSeason.StartDate, comp.CurrentSeason.EndDate)
	fmt.Printf("Stages: %v\n", comp.CurrentSeason.Stages)
}