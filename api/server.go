package api

import (
	"net/http"

	"github.com/rubenhoenle/pixelknecht/model"

	"github.com/gin-gonic/gin"
)

const (
	defaultPixelflutHost = "127.0.0.1"
	defaultPixelflutPort = 1337
)

var server model.PixelflutServer = model.PixelflutServer{
	Host: defaultPixelflutHost,
	Port: defaultPixelflutPort,
}

func getPixelflutServer(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, server)
}

func updatePixelflutServer(c *gin.Context) {
	var updatedServer model.PixelflutServer
	if err := c.BindJSON(&updatedServer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	server.Host = updatedServer.Host
	server.Port = updatedServer.Port
	c.IndentedJSON(http.StatusOK, server)
}
