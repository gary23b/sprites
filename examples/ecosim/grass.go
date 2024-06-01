package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/gary23b/sprites"
)

type GrassHasBeenEaten struct{}

func Main_Grass(sim sprites.Sim, x, y float64) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
	x = float64(int(x))
	y = float64(int(y))

	s := sim.AddSprite("")
	s.Costume("Grass")
	s.SetType(GrassType)

	s.Pos(x, y)
	s.Z(0)
	s.XYScale(10, 10)
	s.Visible(true)
	s.Opacity(10)

	health := 5.0
	spawnCountDown := 30

MainSpriteLoop:
	for {
		time.Sleep(time.Millisecond * 100)

		msgs := s.GetMsgs()
		for _, msg := range msgs {
			switch msg.(type) {
			case GrassHasBeenEaten:
				health -= 50
			}
		}

		if health <= 0 {
			s.DeleteSprite()
			return
		}

		health = min(100, health+1)
		s.Opacity(health)
		if health == 100 {
			spawnCountDown -= 1
			if spawnCountDown <= 0 {
				spawnCountDown = 30
				// spawn more grass
				newX := x + (float64(rand.Intn(3))-1)*10
				newY := y + (float64(rand.Intn(3))-1)*10

				nearMe := sim.WhoIsNearMe(x, y, 20)
				for _, o := range nearMe {
					if o.SpriteType != GrassType {
						continue
					}
					delta := math.Abs(o.X-newX) + math.Abs(o.Y-newY)
					if delta < 1 {
						// fmt.Printf("Found matching grass: %v\n", o)
						continue MainSpriteLoop
					}
				}

				go Main_Grass(sim, newX, newY)
			}
		} else {
			spawnCountDown = 30
		}
	}
}
