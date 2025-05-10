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

func GetModeFromCommanderer() (model.FloodMode, error) {
	response, err := http.Get(config.GetCommandererUrl() + "/api/mode")
	if err != nil {
		return model.FloodMode{}, err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// parse the response
	var mode model.FloodMode
	err = json.Unmarshal([]byte(string(responseData)), &mode)
	if err != nil {
		return model.FloodMode{}, err
	}
	return mode, nil
}

func GetPixelflutServerStringFromCommanderer() (string, error) {
	response, err := http.Get(config.GetCommandererUrl() + "/api/server")
	if err != nil {
		return "", err
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// parse the response
	var server model.PixelflutServer
	err = json.Unmarshal([]byte(string(responseData)), &server)
	if err != nil {
		return "", err
	}
	println(server.Host)
	println(server.Port)
	str := fmt.Sprintf("%s:%d", server.Host, server.Port)
	println(str)
	return str, nil
}
