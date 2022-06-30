package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"padel-backend/main/auth"
	"padel-backend/main/match"
	"padel-backend/main/player"
)

func main() {
	app := gin.Default()

	middleware, err := auth.New("credentials.json", nil)
	if err != nil {
		panic(err)
	}

	//Seed()

	app.Use(CORSMiddleware())

	publicUrlRouterGroup := app.Group("/public/v1")
	publicUrlRouterGroup.GET("/ping", func(c *gin.Context) {
		claims := auth.ExtractClaims(c)
		fmt.Println(claims)
		c.JSON(http.StatusOK, gin.H{
			"message": "pong!",
		})
	})

	authUrlRouterGroup := app.Group("/api/v1")
	authUrlRouterGroup.Use(middleware.MiddlewareFunc())

	player.Init(authUrlRouterGroup)
	match.Init(authUrlRouterGroup)

	app.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().
			Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().
			Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().
			Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().
			Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
