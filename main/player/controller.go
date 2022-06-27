package player

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRoutes(gin *gin.Engine) {
	playerGroupRoute := gin.Group("/player")
	{
		playerGroupRoute.POST("/", PostOnePlayer)

		playerGroupRoute.GET("/", GetAllPlayer)
		playerGroupRoute.GET("/email/:email", GetOnePlayerByEmail)
		playerGroupRoute.GET("/:id", GetOnePlayer)

		playerGroupRoute.PUT("/:id", PutOnePlayer)

		playerGroupRoute.DELETE("/:id", DeleteOnePlayer)
	}
}

func PostOnePlayer(c *gin.Context) {
	var newPlayer Player
	if err := c.BindJSON(&newPlayer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	var player, err = Create(newPlayer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, player)
}

func GetAllPlayer(c *gin.Context) {
	players, err := FindAll()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.IndentedJSON(http.StatusOK, players)
}

func GetOnePlayer(c *gin.Context) {
	var player, err = FindById(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, player)
}

func GetOnePlayerByEmail(c *gin.Context) {
	var player, err = FindById(c.Param("email"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, player)
}

func PutOnePlayer(c *gin.Context) {
	var player, fetchErr = FindById(c.Param("id"))
	if fetchErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fetchErr.Error(),
		})

		return
	}

	var update Player
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
	err := DeleteById(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, nil)
}
