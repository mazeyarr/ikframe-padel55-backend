package match

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func InitRoutes(router *gin.Engine) {
	matchGroupRoute := router.Group("/match")

	matchGroupRoute.POST("/", PostOneMatch)

	matchGroupRoute.GET("/", GetAllMatch)
	matchGroupRoute.GET("/:id", GetOneMatch)

	matchGroupRoute.PUT("/:id", PutOneMatch)
	matchGroupRoute.PUT("/:id/join", PutJoinMatch)
	matchGroupRoute.PUT("/:id/result", PutResultMatch)

	matchGroupRoute.DELETE("/:id", DeleteOneMatch)
}

func PostOneMatch(c *gin.Context) {
	var newMatch Match

	if err := c.BindJSON(&newMatch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, Create(newMatch))
}

func GetAllMatch(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, FindAll())
}

func GetOneMatch(c *gin.Context) {
	var id, _ = strconv.Atoi(c.Param("id"))
	var match, err = FindById(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusBadRequest, match)
}

func PutOneMatch(c *gin.Context) {
	var id, _ = strconv.Atoi(c.Param("id"))
	var match, fetchErr = FindById(id)
	var update Match

	if fetchErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fetchErr.Error(),
		})

		return
	}

	var updatedMatch, err = Update(match, update)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	c.IndentedJSON(http.StatusOK, updatedMatch)
}

func PutJoinMatch(c *gin.Context) {
	var id, _ = strconv.Atoi(c.Param("id"))
	var update JoinMatchRequest

	if parseErr := c.BindJSON(&update); parseErr != nil {
		return
	}

	var match, err = UpdatePlayerToMatchById(id, update)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, match)
}
func PutResultMatch(c *gin.Context) {
	var id, _ = strconv.Atoi(c.Param("id"))
	var result ResultMatchRequest

	if parseErr := c.BindJSON(&result); parseErr != nil {
		return
	}

	var match, err = UpdateMatchResultById(id, result)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, match)
}

func DeleteOneMatch(c *gin.Context) {
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
