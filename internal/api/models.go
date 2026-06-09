package api

import "time"

type Area struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
	Flag string `json:"flag"`
}

type Season struct {
	ID              int       `json:"id"`
	StartDate       string    `json:"startDate"`
	EndDate         string    `json:"endDate"`
	CurrentMatchday int       `json:"currentMatchday"`
	Winner          *Team     `json:"winner"`
	Stages          []string  `json:"stages"`
}

type Team struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	TLA       string `json:"tla"`
	Crest     string `json:"crest"`
}

type Competition struct {
	Area          Area     `json:"area"`
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Code          string   `json:"code"`
	Type          string   `json:"type"`
	Emblem        string   `json:"emblem"`
	CurrentSeason Season   `json:"currentSeason"`
	LastUpdated   time.Time `json:"lastUpdated"`
}