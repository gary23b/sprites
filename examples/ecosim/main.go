package main

import (
	"image"
	"image/color"

	"github.com/gary23b/sprites"
)

var (
	SkyBlue    color.RGBA = color.RGBA{0x87, 0xCE, 0xEB, 0xFF}
	SandyBrown color.RGBA = color.RGBA{0xD2, 0xB4, 0x8C, 0xFF}
	LawnGreen  color.RGBA = color.RGBA{0x7C, 0xFC, 0x00, 0xFF}
)

const (
	UndefinedType int = iota
	GrassType
	BunnyType
)

func main() {
	params := sprites.ScratchParams{Width: 1000, Height: 1000, ShowFPS: true}
	sprites.Start(params, simStartFunc)
}

func simStartFunc(sim sprites.Sim) {
	// Set the background by making a single pixel image and scaling it be be the entire screen.
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, SandyBrown)
	sim.AddCostume(img, "Background")
	s := sim.AddSprite("background")
	s.Costume("Background")
	s.Z(0)
	s.Scale(1010)
	s.Visible(true)

	img = image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, LawnGreen)
	sim.AddCostume(img, "Grass")
	sim.AddCostume(sprites.DecodeCodedSprite(sprites.TurtleImage), "Bunny")

	for y := -300; y < 300; y += 10 {
		for x := -300; x < 300; x += 10 {
			go Main_Grass(sim, float64(x), float64(y))
		}
	}

	// go Main_Grass(sim, 0, 0)
	// go Main_Grass(sim, 50, 0)

	go Main_Bunny(sim, 10, 10)
}
