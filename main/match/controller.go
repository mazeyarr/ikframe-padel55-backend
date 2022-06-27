package match

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRoutes(router *gin.Engine) {
	matchGroupRoute := router.Group("/match")
	{
		matchGroupRoute.POST("/", PostOneMatch)

		matchGroupRoute.GET("/", GetAllMatch)
		matchGroupRoute.GET("/:id", GetOneMatch)

		matchGroupRoute.PUT("/:id", PutOneMatch)
		matchGroupRoute.PUT("/:id/join", PutJoinMatch)
		matchGroupRoute.PUT("/:id/result", PutResultMatch)

		matchGroupRoute.DELETE("/:id", DeleteOneMatch)
	}
}

func PostOneMatch(c *gin.Context) {
	var newMatch Match
	if err := c.BindJSON(&newMatch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	var match, err = Create(newMatch)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, match)
}

func GetAllMatch(c *gin.Context) {
	matches, err := FindAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.IndentedJSON(http.StatusOK, matches)
}

func GetOneMatch(c *gin.Context) {
	var match, err = FindById(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, match)
}

func PutOneMatch(c *gin.Context) {
	var match, fetchErr = FindById(c.Param("id"))
	if fetchErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fetchErr.Error(),
		})

		return
	}

	var update *Match
	if err := c.BindJSON(&update); err != nil {
		return
	}

	var updatedMatch, err = UpdateBasicFields(match, update)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	c.IndentedJSON(http.StatusOK, updatedMatch)
}

func PutJoinMatch(c *gin.Context) {
	var update JoinMatchRequest
	if parseErr := c.BindJSON(&update); parseErr != nil {
		return
	}

	var match, err = UpdatePlayerToMatchById(c.Param("id"), update)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, match)
}
func PutResultMatch(c *gin.Context) {
	var result ResultMatchRequest
	if parseErr := c.BindJSON(&result); parseErr != nil {
		return
	}

	var match, err = UpdateMatchResultById(c.Param("id"), result)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, match)
}

func DeleteOneMatch(c *gin.Context) {
	err := DeleteById(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, nil)
}
