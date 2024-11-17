package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rubenhoenle/pixelknecht/model"
	"github.com/rubenhoenle/pixelknecht/pkg"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

var queue = make(chan string)

func main() {
	initTcpWorkerPool()

	wg.Add(1)
	go commandHandler(3)

	// wait for the goroutines to finish
	wg.Wait()
}

func commandHandler(pollIntervalSec int) {
	// define the initial mode
	var mode model.FloodMode
	mode.Enabled = false

	ctx, cancel := context.WithCancel(context.Background())

	for {
		previousMode := mode
		mode = getModeFromCommanderer()

		// check if the mode changed in the meantime, if so, react to it
		enabledToggled := previousMode.Enabled != mode.Enabled
		posChanged := previousMode.PosY != mode.PosY || previousMode.PosX != mode.PosX
		urlChanged := previousMode.ImageUrl != mode.ImageUrl
		scaleChanged := !pkg.CompareFloat(previousMode.ScaleFactor, mode.ScaleFactor)
		if enabledToggled {
			if mode.Enabled {
				fmt.Println("Starting flooding...")
				ctx, cancel = context.WithCancel(context.Background())
				wg.Add(1)
				go draw(ctx, mode.PosY, mode.PosX, mode.ScaleFactor, mode.ImageUrl)
			} else {
				fmt.Print("Stopping flooding...")
				cancel()
			}
		} else if posChanged || urlChanged || scaleChanged {
			fmt.Println("Restarting flooding with new params...")
			cancel()
			ctx, cancel = context.WithCancel(context.Background())
			wg.Add(1)
			go draw(ctx, mode.PosY, mode.PosX, mode.ScaleFactor, mode.ImageUrl)
		}

		time.Sleep(time.Duration(pollIntervalSec) * time.Second)
	}
}

func ReadEnvWithFallback(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getCommandererUrl() string {
	return ReadEnvWithFallback("COMMANDERER_URL", "http://commanderer.hoenle.xyz:9000")
}

func getModeFromCommanderer() model.FloodMode {
	response, err := http.Get(getCommandererUrl() + "/api/mode")
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

func draw(ctx context.Context, offsetY int, offsetX int, scaleFactor float64, imageUrl string) {
	var frames []model.ParsedFloodImage
	if strings.HasSuffix(strings.ToLower(imageUrl), ".gif") {
		frames = readGif(imageUrl)
	} else {
		frames = readImage(imageUrl, scaleFactor)
	}

	idx, img := 0, frames[0]
	for {
		select {
		case <-ctx.Done(): // if cancel() execute
			wg.Done()
			return
		default:
			for y := 0; y < img.HeightPX; y++ {
				for x := 0; x < img.WidthPX; x++ {
					cmd := fmt.Sprintf("PX %d %d %s\n", x+offsetX, y+offsetY, img.Pixels[y*img.WidthPX+x])
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
