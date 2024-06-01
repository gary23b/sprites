package main

import (
	"image"
	"image/color"
	"math"
	"time"

	"github.com/fogleman/gg"
	"github.com/gary23b/sprites"
)

// https://stackoverflow.com/questions/68339204/double-pendulum-rk4
// https://stackoverflow.com/questions/65224923/i-want-to-have-the-pendulum-blob-in-my-double-pendulum
func main() {
	params := sprites.SimParams{Width: 1000, Height: 1000, ShowFPS: true}
	sprites.Start(params, simStartFunc)
}

func simStartFunc(sim sprites.Sim) {
	// Starting conditions
	theta1 := (math.Pi / 180.0) * 180.1
	theta2 := (math.Pi / 180.0) * 180.0
	w1 := 0.0
	w2 := 0.0
	scaleUpBy := 200.0

	sim.AddCostume(createCircle(20, sprites.White), "White Circle")
	sim.AddCostume(createCircle(18, sprites.Green), "Green Circle")
	sim.AddCostume(createCircle(16, sprites.Red), "Red Circle")
	sim.AddCostume(createRectangle(16, scaleUpBy, sprites.SkyBlue), "connecting bar 1")
	sim.AddCostume(createRectangle(12, scaleUpBy, sprites.Yellow), "connecting bar 2")

	pivot := sim.AddSprite("pivot")
	pivot.Costume("White Circle")
	pivot.Pos(0, 0)
	pivot.Z(1)

	bar1 := sim.AddSprite("bar1")
	bar1.Costume("connecting bar 1")
	bar1.Z(0)

	mass1 := sim.AddSprite("mass1")
	mass1.Costume("Green Circle")
	mass1.Z(3)

	bar2 := sim.AddSprite("bar2")
	bar2.Costume("connecting bar 2")
	bar2.Z(2)

	mass2 := sim.AddSprite("mass2")
	mass2.Costume("Red Circle")
	mass2.Z(4)

	pivot.Visible(true)
	mass1.Visible(true)
	mass2.Visible(true)
	bar1.Visible(true)
	bar2.Visible(true)

	// Loop forever
	for {
		// perform x steps per image update.
		for i := 0; i < 1; i++ {
			w1, w2, theta1, theta2 = step(w1, w2, theta1, theta2)
			x1, y1, x2, y2 := GetPos(w1, w2, theta1, theta2, scaleUpBy)
			mass1.Pos(x1, y1)
			mass2.Pos(x2, y2)

			bar1.Pos(x1/2, y1/2)
			bar1.Angle(theta1 * 180.0 / math.Pi)

			bar2.Pos((x1+x2)/2, (y1+y2)/2)
			bar2.Angle(theta2 * 180.0 / math.Pi)
		}
		time.Sleep(time.Millisecond)
	}
}

func createCircle(radius float64, c color.Color) image.Image {
	dc := gg.NewContext(int(radius*3), int(radius*3))
	dc.DrawCircle(radius*1.5, radius*1.5, radius)
	dc.SetColor(c)
	dc.Fill()
	return dc.Image()
}

func createRectangle(width, height float64, c color.Color) image.Image {
	sideSize := max(width, height)
	dc := gg.NewContext(int(sideSize), int(sideSize))
	dc.DrawRectangle(sideSize/2-width/2, sideSize/2-height/2, width, height)
	dc.SetColor(c)
	dc.Fill()
	return dc.Image()
}
