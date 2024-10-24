package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// TODO: check if we can use the struct from the "commanderer" module here instead of duplicating it
type floodMode struct {
	Enabled  bool   `json:"enabled"`
	PosY     int    `json:"posY"`
	PosX     int    `json:"posX"`
	ImageUrl string `json:"imageUrl"`
}

var wg sync.WaitGroup

var queue = make(chan string)

func main() {
	initTcpWorkerPool()

	wg.Add(1)
	go commandHandler(3)

	// wait for the goroutines to finish
	wg.Wait()
}

func generateColor() string {
	buf := make([]byte, 3)
	_, err := rand.Read(buf)
	if err != nil {
		fmt.Println(err)
	}
	return hex.EncodeToString(buf)
}

func commandHandler(pollIntervalSec int) {
	// define the initial mode
	var mode floodMode
	mode.Enabled = false

	ctx, cancel := context.WithCancel(context.Background())

	for {
		previousMode := mode
		mode = getModeFromCommanderer()

		// check if the mode changed in the meantime, if so, react to it
		if previousMode.Enabled != mode.Enabled {
			if mode.Enabled {
				fmt.Print("Starting flooding...")
				ctx, cancel = context.WithCancel(context.Background())
				wg.Add(1)
				go draw(ctx, mode.PosY, mode.PosX, mode.ImageUrl)
			} else {
				fmt.Print("Stopping flooding...")
				cancel()
			}
		}

		time.Sleep(time.Duration(pollIntervalSec) * time.Second)
	}
}

func getModeFromCommanderer() floodMode {
	response, err := http.Get("http://localhost:9000/mode")
	if err != nil {
		fmt.Print(err.Error())
		// TODO: error handling
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// parse the response
	var mode floodMode
	err = json.Unmarshal([]byte(string(responseData)), &mode)
	if err != nil {
		fmt.Println("Error:", err)
		// TODO: error handling
	}
	return mode
}

func draw(ctx context.Context, offsetY int, offsetX int, imageUrl string) {
	//frames := readGif("https://www.icegif.com/wp-content/uploads/2022/04/icegif-626.gif")
	//frames := []floodImage{readImage("https://wiki.hackerspaces.org/images/8/85/Hackwerk.png")}
	frames := readImage(imageUrl)
	idx, img := 0, frames[0]
	for {
		select {
		case <-ctx.Done(): // if cancel() execute
			wg.Done()
			return
		default:
			for y := 0; y < img.HeightPX; y++ {
				for x := 0; x < img.WidthPX; x++ {
					cmd := fmt.Sprintf("PX %d %d %s\n", x+offsetX, y+offsetY, img.Bytes[y*img.WidthPX+x])
					queue <- cmd
				}
			}

			// go to next frame
			idx++
			if idx >= len(frames) {
				idx = 0
			}
			img = frames[idx]
		}
	}
}
