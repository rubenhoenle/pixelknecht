package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rubenhoenle/pixelknecht/model"
	"net/http"
)

var mode model.FloodMode = model.FloodMode{Enabled: true, PosY: 0, PosX: 0, ScaleFactor: 1, ImageUrl: "https://s3.sfz-aalen.space/static/hackwerk/open.png"}

func getMode(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, mode)
}

func updateMode(c *gin.Context) {
	var updatedMode model.FloodMode
	if err := c.BindJSON(&updatedMode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	mode.Enabled = updatedMode.Enabled
	mode.PosY = updatedMode.PosY
	mode.PosX = updatedMode.PosX
	mode.ScaleFactor = updatedMode.ScaleFactor
	mode.ImageUrl = updatedMode.ImageUrl
	c.IndentedJSON(http.StatusOK, mode)
}
