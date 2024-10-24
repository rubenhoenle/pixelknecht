package main

import (
	"encoding/hex"
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
)

type floodImage struct {
	Bytes    []string
	HeightPX int
	WidthPX  int
}

func readImage(imageURL string) []floodImage {
	resp, err := http.Get(imageURL)
	if err != nil {
		fmt.Println("Failed to download the image:", err)
		// TODO: clarify how to do proper error handling here
		return []floodImage{}
	}
	defer resp.Body.Close()

	img, format, err := image.Decode(resp.Body)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		// TODO: clarify how to do proper error handling here
		return []floodImage{}
	}
	fmt.Printf("Image format: %s\n", format)

	return []floodImage{parseFrame(img)}
}

func readGif(gifURL string) []floodImage {
	resp, err := http.Get(gifURL)
	if err != nil {
		fmt.Println("Failed to download the image:", err)
		// TODO: clarify how to do proper error handling here
		return []floodImage{}
	}
	defer resp.Body.Close()

	// Decode the GIF
	gifImg, err := gif.DecodeAll(resp.Body)
	if err != nil {
		fmt.Println("Error decoding GIF:", err)
		return []floodImage{}
	}

	// Get the number of frames
	numFrames := len(gifImg.Image)
	fmt.Printf("GIF has %d frames\n", numFrames)

	var frames = []floodImage{}
	for _, frame := range gifImg.Image {
		frames = append(frames, parseFrame(frame))
	}
	return frames
}

func parseFrame(img image.Image) floodImage {
	var rgbValues []string
	widthPX := img.Bounds().Dx()
	heightPX := img.Bounds().Dy()
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
