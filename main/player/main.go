package player

import (
	"github.com/gin-gonic/gin"
)

var players = []Player{
	{ID: 1, Name: "Mazeyar Rezaei", Email: "mazeyarr@padel55.nl"},
	{ID: 2, Name: "John Doe", Email: "johndoe@padel55.nl"},
	{ID: 3, Name: "Pete Johnson", Email: "petejohnson@padel55.nl"},
}

type Player struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func Init(router *gin.Engine) {
	router.POST("/player", PostOnePlayer)

	router.GET("/player", GetAllPlayer)
	router.GET("/player/:id", GetOnePlayer)

	router.PUT("/player/:id", PutOnePlayer)

	router.DELETE("/player/:id", DeleteOnePlayer)
}
