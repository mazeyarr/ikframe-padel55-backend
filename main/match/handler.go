package match

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"padel-backend/main/player"
	"padel-backend/main/util"
	"strconv"
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

func fetchOneById(ID int) (*Match, error) {
	for i, m := range matches {
		if m.ID == ID {
			return &matches[i], nil
		}
	}

	return nil, errors.New("match not found")
}

func PostOneMatch(c *gin.Context) {
	var newMatch Match

	if err := c.BindJSON(&newMatch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	matches = append(matches, newMatch)

	c.JSON(http.StatusCreated, newMatch)
}

func GetAllMatch(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, matches)
}

func GetOneMatch(c *gin.Context) {
	var id, _ = strconv.Atoi(c.Param("id"))
	var match, err = fetchOneById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusBadRequest, match)
}

func PutOneMatch(c *gin.Context) {
	var id, _ = strconv.Atoi(c.Param("id"))
	var match, fetchErr = fetchOneById(id)
	var updatedMatch Match

	if fetchErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fetchErr.Error(),
		})

		return
	}

	if err := c.BindJSON(&updatedMatch); err != nil {
		return
	}

	match.Club = updatedMatch.Club
	match.Location = updatedMatch.Location
	match.Time = updatedMatch.Time

	c.IndentedJSON(http.StatusOK, match)
}

func DeleteOneMatch(c *gin.Context) {
	var id, _ = strconv.Atoi(c.Param("id"))
	var _, err = fetchOneById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	// TODO: Remove from database
	// TODO: Remove lines below when we have a database
	for i, m := range matches {
		if m.ID == id {
			matches = util.RemoveElementByIndex(matches, i)
		}
	}

	c.JSON(http.StatusOK, nil)
}
