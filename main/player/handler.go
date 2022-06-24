package player

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"padel-backend/main/util"
	"strconv"
)

type Player struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var players = []Player{
	{ID: 1, Name: "Mazeyar Rezaei", Email: "mazeyarr@padel55.nl"},
	{ID: 2, Name: "John Doe", Email: "johndoe@padel55.nl"},
	{ID: 3, Name: "Pete Johnson", Email: "petejohnson@padel55.nl"},
}

func fetchOneById(ID int) (*Player, error) {
	for i, p := range players {
		if p.ID == ID {
			return &players[i], nil
		}
	}

	return nil, errors.New("player not found")
}

func PostOnePlayer(c *gin.Context) {
	var newPlayer Player

	if err := c.BindJSON(&newPlayer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	players = append(players, newPlayer)

	c.JSON(http.StatusCreated, newPlayer)
}

func GetAllPlayer(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, players)
}

func GetOnePlayer(c *gin.Context) {
	var id, _ = strconv.Atoi(c.Param("id"))
	var player, err = fetchOneById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusBadRequest, player)
}

func PutOnePlayer(c *gin.Context) {
	var id, _ = strconv.Atoi(c.Param("id"))
	var player, fetchErr = fetchOneById(id)
	var updatedPlayer Player

	if fetchErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fetchErr.Error(),
		})

		return
	}

	if err := c.BindJSON(&updatedPlayer); err != nil {
		return
	}

	player.Name = updatedPlayer.Name
	player.Email = updatedPlayer.Email

	c.IndentedJSON(http.StatusOK, player)
}

func DeleteOnePlayer(c *gin.Context) {
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
	for i, p := range players {
		if p.ID == id {
			players = util.RemoveElementByIndex(players, i)
		}
	}

	c.JSON(http.StatusOK, nil)
}
