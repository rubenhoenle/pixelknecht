package model

type FloodMode struct {
	Enabled bool `json:"enabled"`
	// x and y offset
	PosY int `json:"posY"`
	PosX int `json:"posX"`
	// the scale factor for the image
	ScaleFactor float64 `json:"scaleFactor"`
	// the url of the image to paint
	ImageUrl string `json:"imageUrl"`
}

type PixelflutServer struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type ParsedFloodImage struct {
	Pixels   []string
	HeightPX int
	WidthPX  int
}
