package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type pixelflutServer struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

var server pixelflutServer = pixelflutServer{Host: "127.0.0.1", Port: 1234}

func getPixelflutServer(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, server)
}

func updatePixelflutServer(c *gin.Context) {
	var updatedServer pixelflutServer
	if err := c.BindJSON(&updatedServer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	server.Host = updatedServer.Host
	server.Port = updatedServer.Port
	c.IndentedJSON(http.StatusOK, server)
}
