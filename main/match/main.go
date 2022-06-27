package match

import (
	"github.com/gin-gonic/gin"
	"padel-backend/main/player"
	"time"
)

const CollectionMatch = "Match"

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
	ID       string    `json:"id"`
	Club     string    `json:"club"`
	Location string    `json:"location"`
	Time     time.Time `json:"time"`
	TeamA    Team      `json:"teamA"`
	TeamB    Team      `json:"teamB"`
	Locked   bool      `json:"locked"`
}

type TeamSelection = string

const (
	None  TeamSelection = "NONE"
	TeamA               = "TEAM_A"
	TeamB               = "TEAM_B"
)

type JoinMatchRequest struct {
	PlayerId string        `json:"playerId"`
	Team     TeamSelection `json:"team"`
}

type ResultMatchRequest struct {
	PlayerId   string        `json:"playerId"`
	Team       TeamSelection `json:"team"`
	TeamResult TeamResult    `json:"result"`
}

func Init(router *gin.RouterGroup) {
	router.Use(InitMatchService)

	InitRoutes(router)
}
