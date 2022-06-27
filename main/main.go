package main

import (
	"github.com/gin-gonic/gin"
	"padel-backend/main/match"
	"padel-backend/main/player"
)

func main() {
	gin := gin.Default()

	//Seed()

	player.Init(gin)
	match.Init(gin)

	gin.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
