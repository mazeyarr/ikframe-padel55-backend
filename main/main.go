package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"padel-backend/main/match"
	"padel-backend/main/player"
)

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	player.Init(router)
	match.Init(router)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
