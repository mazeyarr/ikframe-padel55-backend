package match

import (
	"github.com/gin-gonic/gin"
	"padel-backend/main/player"
	"time"
)

type TeamResult struct {
	Set   int `json:"set"`
	Score int `json:"score"`
}

type Team struct {
	Player1 player.Player `json:"player1"`
	Player2 player.Player `json:"player2"`
	Results []TeamResult  `json:"results"`
}

type Match struct {
	ID       int       `json:"id"`
	Club     string    `json:"club"`
	Location string    `json:"location"`
	Time     time.Time `json:"time"`
	TeamA    Team      `json:"teamA"`
	TeamB    Team      `json:"teamB"`
	locked   bool      `json:"locked"`
}

type TeamSelection = string

const (
	None  TeamSelection = "NONE"
	TeamA               = "TEAM_A"
	TeamB               = "TEAM_B"
)

type JoinMatchRequest struct {
	PlayerId int           `json:"playerId"`
	Team     TeamSelection `json:"team"`
}

type ResultMatchRequest struct {
	PlayerId   int           `json:"playerId"`
	Team       TeamSelection `json:"team"`
	TeamResult TeamResult    `json:"result"`
}

var matches = []Match{
	{
		ID:       1,
		Club:     "B. Amsterdam",
		Location: "Johan Huizingalaan 768A",
		Time:     time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		TeamA: Team{
			Results: []TeamResult{},
		},
		TeamB: Team{
			Results: []TeamResult{},
		},
		locked: false,
	},
}

func Init(router *gin.Engine) {
	InitRoutes(router)
}
