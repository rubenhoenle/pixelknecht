package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type floodMode struct {
	Enabled bool `json:"enabled"`
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/mode", getMode)
	router.PUT("/mode", updateMode)

	router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
		template := Page()
		template.Render(c.Request.Context(), c.Writer)
	})
    router.Static("/assets", "./assets")
	return router
}

var mode = floodMode{Enabled: true}

func main() {
	router := setupRouter()
	router.Run("localhost:9000")
}

func getMode(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, mode)
}

func updateMode(c *gin.Context) {
	var updatedMode floodMode
	if err := c.BindJSON(&updatedMode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	mode.Enabled = updatedMode.Enabled
	c.IndentedJSON(http.StatusOK, mode)
}
