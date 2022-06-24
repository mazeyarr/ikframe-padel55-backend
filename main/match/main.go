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
}

type JoinTeam = int64

const (
	None JoinTeam = iota
	TeamA
	TeamB
)

type JoinMatchRequest struct {
	PlayerId int      `json:"playerId"`
	Team     JoinTeam `json:"team"`
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
	},
}

func Init(router *gin.Engine) {
	router.POST("/match", PostOneMatch)

	router.GET("/match", GetAllMatch)
	router.GET("/match/:id", GetOneMatch)

	router.PUT("/match/:id", PutOneMatch)

	router.DELETE("/match/:id", DeleteOneMatch)
}
