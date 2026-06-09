package main

import (
	"fmt"
	"os"

	"github.com/charliehustlr1792/fifawc26/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "config error:", err)
		os.Exit(1)
	}

	fmt.Println("fifawc26 starting...")
	fmt.Println("API key loaded:", maskKey(cfg.APIKey))
	fmt.Println("Cache dir:", cfg.CacheDir)
	fmt.Println("Poll interval:", cfg.PollInterval)
}

func maskKey(k string) string {
	if len(k) <= 6 {
		return "****"
	}
	return k[:3] + "..." + k[len(k)-3:]
}