package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

// TODO: check if we can use the struct from the "commanderer" module here instead of duplicating it
type floodMode struct {
	Enabled bool `json:"enabled"`
}

type floodImage struct {
	Bytes    []string
	HeightPX int
	WidthPX  int
}

var wg sync.WaitGroup

func main() {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	wg.Add(1)
	go commandHandler(conn, 3)

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

func commandHandler(conn net.Conn, pollIntervalSec int) {
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
				go draw(ctx, conn)
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

func readImage(filename string) floodImage {
	imgFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error: ", err)
		// TODO: clarify how to do proper error handling here
		return floodImage{}
	}
	defer imgFile.Close()

	// decode the image
	img, _, err := image.Decode(imgFile)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		// TODO: clarify how to do proper error handling here
		return floodImage{}
	}

	widthPX := img.Bounds().Dx()
	heightPX := img.Bounds().Dy()

	var rgbValues []string

	for y := 0; y < heightPX; y++ {
		for x := 0; x < widthPX; x++ {
			// Get the color of the pixel at (x, y)
			r, g, b, _ := img.At(x, y).RGBA()

			// Convert to 8-bit RGB
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			buf := []uint8{r8, g8, b8}

			rgb := hex.EncodeToString(buf)
			rgbValues = append(rgbValues, rgb)
		}
	}
	return floodImage{HeightPX: heightPX, WidthPX: widthPX, Bytes: rgbValues}
}

func draw(ctx context.Context, conn net.Conn) {
	for {
		img := readImage("image.png")
		select {
		case <-ctx.Done(): // if cancel() execute
			wg.Done()
			return
		default:
			for y := 0; y < img.HeightPX; y++ {
				for x := 0; x < img.WidthPX; x++ {
					cmd := fmt.Sprintf("PX %d %d %s\n", y, x, img.Bytes[y*img.WidthPX+x])
					_, err := conn.Write([]byte(cmd))
					if err != nil {
						fmt.Println(err)
						wg.Done()
						return
					}
				}
			}
		}
	}
}
