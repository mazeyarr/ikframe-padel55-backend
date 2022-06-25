package player

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func InitRoutes(router *gin.Engine) {
	playerGroupRoute := router.Group("/player")
	playerGroupRoute.POST("/", PostOnePlayer)

	playerGroupRoute.GET("/", GetAllPlayer)
	playerGroupRoute.GET("/:id", GetOnePlayer)

	playerGroupRoute.PUT("/:id", PutOnePlayer)

	playerGroupRoute.DELETE("/:id", DeleteOnePlayer)
}

func PostOnePlayer(c *gin.Context) {
	var newPlayer Player

	if err := c.BindJSON(&newPlayer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, Create(newPlayer))
}

func GetAllPlayer(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, FindAll())
}

func GetOnePlayer(c *gin.Context) {
	var id, _ = strconv.Atoi(c.Param("id"))
	var player, err = FindById(id)

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
	var player, fetchErr = FindById(id)
	var update Player

	if fetchErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fetchErr.Error(),
		})

		return
	}

	if err := c.BindJSON(&update); err != nil {
		return
	}

	var updatedPlayer, err = Update(player, update)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	c.IndentedJSON(http.StatusOK, updatedPlayer)
}

func DeleteOnePlayer(c *gin.Context) {
	var id, _ = strconv.Atoi(c.Param("id"))

	err := DeleteById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, nil)
}
