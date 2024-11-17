package fetcher

import (
	"encoding/json"
	"fmt"
	"github.com/rubenhoenle/pixelknecht/config"
	"github.com/rubenhoenle/pixelknecht/model"
	"io"
	"log"
	"net/http"
)

func GetModeFromCommanderer() model.FloodMode {
	response, err := http.Get(config.GetCommandererUrl() + "/api/mode")
	if err != nil {
		fmt.Print(err.Error())
		// TODO: error handling
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// parse the response
	var mode model.FloodMode
	err = json.Unmarshal([]byte(string(responseData)), &mode)
	if err != nil {
		fmt.Println("Error:", err)
		// TODO: error handling
	}
	return mode
}

func GetPixelflutServerStringFromCommanderer() string {
	response, err := http.Get(config.GetCommandererUrl() + "/api/server")
	if err != nil {
		fmt.Print(err.Error())
		// TODO: error handling
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// parse the response
	var server model.PixelflutServer
	err = json.Unmarshal([]byte(string(responseData)), &server)
	if err != nil {
		fmt.Println("Error:", err)
		// TODO: error handling
	}
	println(server.Host)
	println(server.Port)
	str := fmt.Sprintf("%s:%d", server.Host, server.Port)
	println(str)
	return str
}
