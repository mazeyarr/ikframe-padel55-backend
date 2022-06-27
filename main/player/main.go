package player

import (
	"github.com/gin-gonic/gin"
)

const CollectionPlayer = "Player"

type Player struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func Init(router *gin.Engine) {
	router.Use(InitPlayerService)

	InitRoutes(router)
}
