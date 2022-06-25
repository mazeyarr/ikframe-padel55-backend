package player

import (
	"github.com/gin-gonic/gin"
)

var players = []Player{
	{ID: 1, Name: "Mazeyar Rezaei", Email: "mazeyarr@padel55.nl"},
	{ID: 2, Name: "John Doe", Email: "johndoe@padel55.nl"},
	{ID: 3, Name: "Pete Johnson", Email: "petejohnson@padel55.nl"},
	{ID: 4, Name: "Rober James", Email: "roberjames@padel55.nl"},
}

type Player struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func Init(router *gin.Engine) {
	InitRoutes(router)
}
