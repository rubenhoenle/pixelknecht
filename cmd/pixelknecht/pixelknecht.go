package main

import (
	"context"
	"fmt"
	"github.com/rubenhoenle/pixelknecht/fetcher"
	"github.com/rubenhoenle/pixelknecht/imgparser"
	"github.com/rubenhoenle/pixelknecht/model"
	"github.com/rubenhoenle/pixelknecht/pkg"
	"github.com/rubenhoenle/pixelknecht/tcpworker"
	"strings"
	"sync"
	"time"
)

const workerPoolSize = 15

func main() {
	var wg sync.WaitGroup
	var queue = make(chan string)

	/* init TCP worker pool */
	pixelflutConnectionString := fetcher.GetPixelflutServerStringFromCommanderer()
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go tcpworker.TcpWorker(&wg, queue, pixelflutConnectionString)
	}

	wg.Add(1)
	go commandHandler(3, &wg, queue)

	// wait for the goroutines to finish
	wg.Wait()
}

func commandHandler(pollIntervalSec int, wg *sync.WaitGroup, queue chan<- string) {
	// define the initial mode
	var mode model.FloodMode
	mode.Enabled = false

	ctx, cancel := context.WithCancel(context.Background())

	for {
		previousMode := mode
		mode = fetcher.GetModeFromCommanderer()

		// check if the mode changed in the meantime, if so, react to it
		enabledToggled := previousMode.Enabled != mode.Enabled
		posChanged := previousMode.PosY != mode.PosY || previousMode.PosX != mode.PosX
		urlChanged := previousMode.ImageUrl != mode.ImageUrl
		scaleChanged := !pkg.CompareFloat(previousMode.ScaleFactor, mode.ScaleFactor)
		if enabledToggled && mode.Enabled {
			fmt.Println("Starting flooding...")
			ctx, cancel = context.WithCancel(context.Background())
			wg.Add(1)
			go draw(ctx, wg, queue, mode.PosY, mode.PosX, mode.ScaleFactor, mode.ImageUrl)
		} else if enabledToggled {
			fmt.Print("Stopping flooding...")
			cancel()
		} else if posChanged || urlChanged || scaleChanged {
			fmt.Println("Restarting flooding with new params...")
			cancel()
			ctx, cancel = context.WithCancel(context.Background())
			wg.Add(1)
			go draw(ctx, wg, queue, mode.PosY, mode.PosX, mode.ScaleFactor, mode.ImageUrl)
		}

		time.Sleep(time.Duration(pollIntervalSec) * time.Second)
	}
}

func draw(ctx context.Context, wg *sync.WaitGroup, queue chan<- string, offsetY int, offsetX int, scaleFactor float64, imageUrl string) {
	var frames []model.ParsedFloodImage
	if strings.HasSuffix(strings.ToLower(imageUrl), ".gif") {
		frames = imgparser.ReadGif(imageUrl)
	} else {
		frames = imgparser.ReadImage(imageUrl, scaleFactor)
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
