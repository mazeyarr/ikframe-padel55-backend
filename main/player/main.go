package player

import (
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	router.POST("/player", PostOnePlayer)

	router.GET("/player", GetAllPlayer)
	router.GET("/player/:id", GetOnePlayer)

	router.PUT("/player/:id", PutOnePlayer)

	router.DELETE("/player/:id", DeleteOnePlayer)
}
