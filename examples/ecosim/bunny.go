package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/gary23b/sprites"
	"github.com/gary23b/sprites/spritesmodels"
)

type bunny struct {
	sim    sprites.Sim
	sprite sprites.Sprite
	food   *spritesmodels.NearMeInfo
	x, y   float64
	health float64

	wanderX, wanderY float64

	spawnCounter int
}

func Main_Bunny(sim sprites.Sim, x, y float64) {

	s := sim.AddSprite("")
	s.Costume("Bunny")
	s.SetType(BunnyType)

	s.Pos(x, y)
	s.Z(1)
	s.XYScale(.2, .2)
	s.Visible(true)
	s.Opacity(100)

	b := bunny{
		sim:          sim,
		sprite:       s,
		x:            x,
		y:            y,
		health:       50,
		spawnCounter: 200,
	}
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))

	b.main()
}

func (s *bunny) main() {
	// MainSpriteLoop:
	for {
		time.Sleep(time.Millisecond * 20)
		s.health -= .1

		switch {
		case s.health <= 0:
			s.sprite.DeleteSprite()
			return
		case s.health < 75 && s.findFood():
			s.feed()
		case s.breed():
			//
		default:
			s.wander()
		}

		s.sprite.Angle(s.health)
	}
}

func (s *bunny) findFood() bool {

	if s.food == nil {
		aroundMe := s.sprite.WhoIsNearMe(100)
		minDis := math.MaxFloat64
		for i, o := range aroundMe {
			if o.SpriteType == GrassType {
				deltaX := o.X - s.x
				deltaY := o.Y - s.y
				dist := deltaX*deltaX + deltaY*deltaY
				if dist < minDis {
					minDis = dist
					s.food = &aroundMe[i]
				}
			}
		}
	}

	return s.food != nil
}

func (s *bunny) feed() {
	if s.food == nil {
		return
	}

	if s.food.SpriteType != GrassType {
		return
	}

	deltaX := s.food.X - s.x
	deltaY := s.food.Y - s.y
	dist := math.Sqrt(deltaX*deltaX + deltaY*deltaY)
	if dist < 5 {
		foodStatus := s.sim.GetSpriteInfoByID(s.food.SpriteID)
		if foodStatus.Deleted {
			s.food = nil
			return
		}

		s.health = min(100, s.health+25)
		s.sim.SendMsg(s.food.SpriteID, GrassHasBeenEaten{})
		s.food = nil
		return
	}

	speed := 1.0
	deltaX /= dist * speed
	deltaY /= dist * speed
	s.x += deltaX
	s.y += deltaY
	s.sprite.Pos(s.x, s.y)
}

func (s *bunny) breed() bool {
	if s.health < 90 {
		return false
	}
	s.spawnCounter -= 1
	if s.spawnCounter <= 0 {
		s.spawnCounter = 200
		s.health = 25
		go Main_Bunny(s.sim, s.x, s.y)
		return true
	}
	return false
}

func (s *bunny) wander() {
	s.wanderX += rand.NormFloat64() * .01
	s.wanderY += rand.NormFloat64() * .01

	s.wanderX = max(-.5, min(.5, s.wanderX))
	s.wanderY = max(-.5, min(.5, s.wanderY))

	s.x += s.wanderX
	s.y += s.wanderY
	s.sprite.Pos(s.x, s.y)
}
