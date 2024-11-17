package main

import (
	"github.com/rubenhoenle/pixelknecht/api"
	"github.com/rubenhoenle/pixelknecht/config"
)

func main() {
	router := api.SetupRouter()
	router.Run(config.GetListenerUrl())
}
