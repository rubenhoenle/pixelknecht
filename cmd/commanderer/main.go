package main

import (
	"embed"
	"net/http"

	"github.com/rubenhoenle/pixelknecht/config"
	"github.com/rubenhoenle/pixelknecht/model"

	"github.com/gin-gonic/gin"
)

//go:embed static/*
var static embed.FS

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/api/mode", getMode)
	router.PUT("/api/mode", updateMode)
	router.GET("/api/server", getPixelflutServer)
	router.PUT("/api/server", updatePixelflutServer)

	trustedProxy := config.GetTrustedProxy()
	if trustedProxy != "" {
		router.SetTrustedProxies([]string{trustedProxy})
	} else {
		router.SetTrustedProxies(nil)
	}

	router.StaticFS("/static", http.FS(static))
	router.GET("/", func(c *gin.Context) {
		data, err := static.ReadFile("static/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading index.html: %s", err)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	return router
}

var mode model.FloodMode

func main() {
	mode = model.FloodMode{Enabled: true, PosY: 0, PosX: 0, ScaleFactor: 1, ImageUrl: "https://s3.sfz-aalen.space/static/hackwerk/open.png"}
	router := setupRouter()
	router.Run(config.GetListenerUrl())
}

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
