package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/gary23b/sprites"
	"github.com/gary23b/sprites/sprite"
	"github.com/gary23b/sprites/spritesmodels"
)

func main() {
	params := sprites.ScratchParams{Width: 500, Height: 500}
	sprites.Start(params, simStartFunc)
}

func simStartFunc(sim spritesmodels.Scratch) {
	sim.AddCostume(sprite.DecodeCodedSprite(sprite.TurtleImage), "t")

	// go scratch.CreateGifDithered(sim, time.Millisecond*100, time.Millisecond*100, "./examples/randomwalk/randomwalk.gif", 100)

	for {
		go turtleRandomWalk(sim)
		time.Sleep(time.Millisecond * 500)
	}
}

func turtleRandomWalk(sim spritesmodels.Scratch) {
	s := sim.AddSprite("")
	s.Costume("t")
	s.Visible(true)

	x, y := 0.0, 0.0
	velX, velY := 0.0, 0.0

	for {
		velX += (rand.Float64()*2 - 1) * .05
		velY += (rand.Float64()*2 - 1) * .05
		x += velX
		y += velY
		s.Pos(x, y)
		s.Angle(180.0 / math.Pi * math.Atan2(velY, velX))

		if max(math.Abs(x), math.Abs(y)) > 1000 {
			s.DeleteSprite()
			return
		}
		time.Sleep(time.Millisecond * 10)
	}
}
