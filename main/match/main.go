package match

import (
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	router.POST("/match", PostOneMatch)

	router.GET("/match", GetAllMatch)
	router.GET("/match/:id", GetOneMatch)

	router.PUT("/match/:id", PutOneMatch)

	router.DELETE("/match/:id", DeleteOneMatch)
}
