package main

import (
	"github.com/rubenhoenle/pixelknecht/commanderer/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type floodMode struct {
	Enabled bool `json:"enabled"`
	// x and y offset
	PosY int `json:"posY"`
	PosX int `json:"posX"`
	// the url of the image to paint
	ImageUrl string `json:"imageUrl"`
	// the IP/hostname of the pixelflut server
	ServerHost string `json:"serverHost"`
	// the port of the pixelflut server
	ServerPort int `json:"serverPort"`
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/mode", getMode)
	router.PUT("/mode", updateMode)

	trustedProxy := config.GetTrustedProxy()
	if trustedProxy != "" {
		router.SetTrustedProxies([]string{trustedProxy})
	} else {
		router.SetTrustedProxies(nil)
	}

	return router
}

var mode floodMode

func main() {
	mode = floodMode{Enabled: true, PosY: 0, PosX: 0, ServerHost: "127.0.0.1", ServerPort: 1234, ImageUrl: "https://s3.sfz-aalen.space/static/hackwerk/open.png"}
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
	mode.ServerHost = updatedMode.ServerHost
	mode.ServerPort = updatedMode.ServerPort
	c.IndentedJSON(http.StatusOK, mode)
}
