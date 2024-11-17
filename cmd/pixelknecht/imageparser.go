package main

import (
	"encoding/hex"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/rubenhoenle/pixelknecht/model"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
)

func readImage(imageURL string, scaleFactor float64) []model.ParsedFloodImage {
	resp, err := http.Get(imageURL)
	if err != nil {
		fmt.Println("Failed to download the image:", err)
		// TODO: clarify how to do proper error handling here
		return []model.ParsedFloodImage{}
	}
	defer resp.Body.Close()

	img, format, err := image.Decode(resp.Body)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		// TODO: clarify how to do proper error handling here
		return []model.ParsedFloodImage{}
	}

	fmt.Printf("Image format: %s\n", format)

	if !CompareFloat(1, scaleFactor) {
		// scale the image
		newWidth := uint(float64(img.Bounds().Dx()) * scaleFactor)
		newHeight := uint(float64(img.Bounds().Dy()) * scaleFactor)
		img = resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	}

	return []model.ParsedFloodImage{parseFrame(img)}
}

func readGif(gifURL string) []model.ParsedFloodImage {
	resp, err := http.Get(gifURL)
	if err != nil {
		fmt.Println("Failed to download the image:", err)
		// TODO: clarify how to do proper error handling here
		return []model.ParsedFloodImage{}
	}
	defer resp.Body.Close()

	// Decode the GIF
	gifImg, err := gif.DecodeAll(resp.Body)
	if err != nil {
		fmt.Println("Error decoding GIF:", err)
		return []model.ParsedFloodImage{}
	}

	// Get the number of frames
	numFrames := len(gifImg.Image)
	fmt.Printf("GIF has %d frames\n", numFrames)

	var frames = []model.ParsedFloodImage{}
	for _, frame := range gifImg.Image {
		frames = append(frames, parseFrame(frame))
	}
	return frames
}

func parseFrame(img image.Image) model.ParsedFloodImage {
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
	return model.ParsedFloodImage{HeightPX: heightPX, WidthPX: widthPX, Pixels: rgbValues}
}
