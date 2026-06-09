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

type StandingRow struct {
	Position       int    `json:"position"`
	Team           Team   `json:"team"`
	PlayedGames    int    `json:"playedGames"`
	Form           string `json:"form"`
	Won            int    `json:"won"`
	Draw           int    `json:"draw"`
	Lost           int    `json:"lost"`
	Points         int    `json:"points"`
	GoalsFor       int    `json:"goalsFor"`
	GoalsAgainst   int    `json:"goalsAgainst"`
	GoalDifference int    `json:"goalDifference"`
}

type Standing struct {
	Stage string        `json:"stage"`
	Type  string        `json:"type"`
	Group string        `json:"group"`
	Table []StandingRow `json:"table"`
}

type StandingsResponse struct {
	Competition Competition `json:"competition"`
	Season      Season      `json:"season"`
	Standings   []Standing  `json:"standings"`
}

type Score struct {
	Winner   string `json:"winner"`
	Duration string `json:"duration"`
	FullTime struct {
		Home *int `json:"home"`
		Away *int `json:"away"`
	} `json:"fullTime"`
	HalfTime struct {
		Home *int `json:"home"`
		Away *int `json:"away"`
	} `json:"halfTime"`
}

type Match struct {
	ID          int       `json:"id"`
	UTCDate     time.Time `json:"utcDate"`
	Status      string    `json:"status"`
	Matchday    int       `json:"matchday"`
	Stage       string    `json:"stage"`
	Group       string    `json:"group"`
	HomeTeam    Team      `json:"homeTeam"`
	AwayTeam    Team      `json:"awayTeam"`
	Score       Score     `json:"score"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type MatchesResponse struct {
	Competition Competition `json:"competition"`
	Matches     []Match     `json:"matches"`
	ResultSet   struct {
		Count  int    `json:"count"`
		First  string `json:"first"`
		Last   string `json:"last"`
		Played int    `json:"played"`
	} `json:"resultSet"`
}

type Player struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Nationality string `json:"nationality"`
	Position    string `json:"position"`
}

type Scorer struct {
	Player    Player `json:"player"`
	Team      Team   `json:"team"`
	Goals     int    `json:"goals"`
	Assists   int    `json:"assists"`
	Penalties int    `json:"penalties"`
}

type ScorersResponse struct {
	Count       int         `json:"count"`
	Competition Competition `json:"competition"`
	Season      Season      `json:"season"`
	Scorers     []Scorer    `json:"scorers"`
}