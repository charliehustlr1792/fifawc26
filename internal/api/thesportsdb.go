package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const sportsDBBase = "https://www.thesportsdb.com/api/v1/json/123"

type TheSportsDBClient struct {
	http     *http.Client
	leagueID string
	season   string
}

func NewTheSportsDBClient() *TheSportsDBClient {
	return &TheSportsDBClient{
		http:     &http.Client{Timeout: 15 * time.Second},
		leagueID: "4429",
		season:   "2026",
	}
}

type sdbEvent struct {
	IdEvent        string `json:"idEvent"`
	StrEvent       string `json:"strEvent"`
	StrHomeTeam    string `json:"strHomeTeam"`
	StrAwayTeam    string `json:"strAwayTeam"`
	IntHomeScore   string `json:"intHomeScore"`
	IntAwayScore   string `json:"intAwayScore"`
	DateEvent      string `json:"dateEvent"`
	StrTime        string `json:"strTime"`
	StrTimestamp   string `json:"strTimestamp"`
	StrStatus      string `json:"strStatus"`
	IntRound       string `json:"intRound"`
	StrGroup       string `json:"strGroup"`
	StrHomeTeamTLA string `json:"strHomeTeamShort"`
	StrAwayTeamTLA string `json:"strAwayTeamShort"`
}

type sdbEventsResponse struct {
	Events []sdbEvent `json:"events"`
}


func (c *TheSportsDBClient) GetCompetition(ctx context.Context, code string) (*Competition, error) {
	return &Competition{
		ID:   2000,
		Name: "FIFA World Cup",
		Code: "WC",
		Type: "CUP",
		Area: Area{Name: "World"},
		CurrentSeason: Season{
			StartDate: "2026-06-11",
			EndDate:   "2026-07-19",
		},
	}, nil
}

func (c *TheSportsDBClient) GetStandings(ctx context.Context, code string) (*StandingsResponse, error) {
	return &StandingsResponse{
		Competition: Competition{Name: "FIFA World Cup", Code: "WC"},
		Standings:   []Standing{},
	}, nil
}

func (c *TheSportsDBClient) GetMatches(ctx context.Context, code string, f MatchFilter) (*MatchesResponse, error) {
	url := fmt.Sprintf("%s/eventsseason.php?id=%s&s=%s", sportsDBBase, c.leagueID, c.season)
	var raw sdbEventsResponse
	if err := c.getJSON(ctx, url, &raw); err != nil {
		return nil, err
	}

	matches := make([]Match, 0, len(raw.Events))
	for _, e := range raw.Events {
		m := sdbToMatch(e)
		if f.Status != "" && !strings.EqualFold(m.Status, f.Status) {
			continue
		}
		if f.Matchday > 0 && m.Matchday != f.Matchday {
			continue
		}
		matches = append(matches, m)
	}

	return &MatchesResponse{
		Competition: Competition{Name: "FIFA World Cup", Code: "WC"},
		Matches:     matches,
		ResultSet: struct {
			Count  int    `json:"count"`
			First  string `json:"first"`
			Last   string `json:"last"`
			Played int    `json:"played"`
		}{Count: len(matches)},
	}, nil
}

func (c *TheSportsDBClient) GetTeam(ctx context.Context, id int) (*TeamDetail, error) {
	return &TeamDetail{ID: id}, nil
}

func (c *TheSportsDBClient) GetScorers(ctx context.Context, code string, limit int) (*ScorersResponse, error) {
	return &ScorersResponse{
		Competition: Competition{Name: "FIFA World Cup", Code: "WC"},
		Scorers:     []Scorer{},
	}, nil
}

func (c *TheSportsDBClient) getJSON(ctx context.Context, url string, out any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("thesportsdb: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("thesportsdb %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func sdbToMatch(e sdbEvent) Match {
	m := Match{
		Status:   normalizeSDBStatus(e.StrStatus),
		Matchday: atoi(e.IntRound),
		Stage:    "GROUP_STAGE",
		Group:    e.StrGroup,
		HomeTeam: Team{Name: e.StrHomeTeam, TLA: e.StrHomeTeamTLA},
		AwayTeam: Team{Name: e.StrAwayTeam, TLA: e.StrAwayTeamTLA},
	}
	if id, err := strconv.Atoi(e.IdEvent); err == nil {
		m.ID = id
	}
	if e.StrTimestamp != "" {
		if t, err := time.Parse(time.RFC3339, e.StrTimestamp); err == nil {
			m.UTCDate = t
		}
	}
	if m.UTCDate.IsZero() && e.DateEvent != "" {
		layout := "2006-01-02 15:04:05"
		s := e.DateEvent
		if e.StrTime != "" {
			s = e.DateEvent + " " + e.StrTime
		}
		if t, err := time.Parse(layout, s); err == nil {
			m.UTCDate = t.UTC()
		}
	}
	if h, err := strconv.Atoi(e.IntHomeScore); err == nil {
		m.Score.FullTime.Home = &h
	}
	if a, err := strconv.Atoi(e.IntAwayScore); err == nil {
		m.Score.FullTime.Away = &a
	}
	return m
}

func normalizeSDBStatus(s string) string {
	switch strings.ToUpper(s) {
	case "MATCH FINISHED", "FT", "FINISHED":
		return "FINISHED"
	case "NOT STARTED", "SCHEDULED", "NS":
		return "SCHEDULED"
	case "IN PROGRESS", "1H", "2H", "HT":
		return "IN_PLAY"
	default:
		if s == "" {
			return "SCHEDULED"
		}
		return strings.ToUpper(s)
	}
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}