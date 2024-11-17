package api

import (
	"embed"
	"net/http"

	"github.com/rubenhoenle/pixelknecht/config"

	"github.com/gin-gonic/gin"
)

//go:embed static/*
var static embed.FS

func SetupRouter() *gin.Engine {
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
