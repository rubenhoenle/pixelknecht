package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/rubenhoenle/pixelknecht/commanderer/config"
	"net/http"
	"strconv"
	"time"
)

//go:embed static/*
var static embed.FS

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

func uploadFile(c *gin.Context) {
	// single file
	file, _ := c.FormFile("file")
	// Upload the file to specific dst.
	c.SaveUploadedFile(file, "temp-files/upload-"+strconv.FormatInt(time.Now().Unix(), 10)+".png")

}
func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/mode", getMode)
	router.StaticFS("/pictures", http.Dir("./temp-files"))

	router.MaxMultipartMemory = 10 << 20 // 8 MiB
	router.POST("/upload", uploadFile)
	router.PUT("/mode", updateMode)

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

var mode floodMode

func main() {
	mode = floodMode{Enabled: true, PosY: 0, PosX: 0, ServerHost: "127.0.0.1", ServerPort: 1234, ImageUrl: "https://s3.sfz-aalen.space/static/hackwerk/open.png"}
	router := setupRouter()
	router.Run(config.ReadEnvWithFallback("COMMANDERER_LISTEN_HOST", "localhost") + ":9000")
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
