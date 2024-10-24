package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type floodMode struct {
	Enabled  bool   `json:"enabled"`
	PosY     int    `json:"posY"`
	PosX     int    `json:"posX"`
	ImageUrl string `json:"imageUrl"`
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/mode", getMode)
	router.PUT("/mode", updateMode)
	return router
}

var mode floodMode

func main() {
	mode = floodMode{Enabled: true, PosY: 0, PosX: 0, ImageUrl: "https://s3.sfz-aalen.space/static/hackwerk/open.png"}
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
	mode.PosY = updatedMode.PosY
	mode.PosX = updatedMode.PosX
	mode.ImageUrl = updatedMode.ImageUrl
	c.IndentedJSON(http.StatusOK, mode)
}
