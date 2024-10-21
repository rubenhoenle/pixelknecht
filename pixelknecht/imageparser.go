package main

import (
	"encoding/hex"
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type floodImage struct {
	Bytes    []string
	HeightPX int
	WidthPX  int
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

func readGif(filename string) []floodImage {
	// Open the GIF file
	gifFile, err := os.Open(filename) // Replace with your GIF file path
	if err != nil {
		fmt.Println("Error: ", err)
		return []floodImage{}
	}
	defer gifFile.Close()

	// Decode the GIF
	gifImg, err := gif.DecodeAll(gifFile)
	if err != nil {
		fmt.Println("Error decoding GIF:", err)
		return []floodImage{}
	}

	// Get the number of frames
	numFrames := len(gifImg.Image)
	fmt.Printf("GIF has %d frames\n", numFrames)

	//frames := make([]floodImage, numFrames)
	var frames = []floodImage{}

	// Loop over each frame
	for frameIndex, frame := range gifImg.Image {
		width := frame.Bounds().Dx()
		height := frame.Bounds().Dy()

		fmt.Printf("Frame %d dimensions: %d x %d\n", frameIndex, width, height)

		var rgbValues []string

		// Loop over each pixel in the frame
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				// Get the color of the pixel at (x, y)
				r, g, b, _ := frame.At(x, y).RGBA()

				// Convert to 8-bit RGB
				r8 := uint8(r >> 8)
				g8 := uint8(g >> 8)
				b8 := uint8(b >> 8)

				buf := []uint8{r8, g8, b8}

				rgb := hex.EncodeToString(buf)
				rgbValues = append(rgbValues, rgb)
			}
		}

		// For demonstration, print the first few RGB values of the frame
		//fmt.Printf("Frame %d RGB values: %v...\n", frameIndex, rgbValues[:9])
		frames = append(frames, floodImage{HeightPX: height, WidthPX: width, Bytes: rgbValues})
	}
	return frames
}
