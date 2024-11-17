package api

import (
	"net/http"

	"github.com/rubenhoenle/pixelknecht/model"

	"github.com/gin-gonic/gin"
)

var server model.PixelflutServer = model.PixelflutServer{Host: "127.0.0.1", Port: 1337}

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
