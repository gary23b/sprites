package main

import (
	"math/rand"
	"time"

	"github.com/gary23b/sprites/ebitensim"
	"github.com/gary23b/sprites/models"
)

func main() {
	ebitensim.StartSim(ebitensim.SimParams{Width: 1000, Height: 1000, ShowFPS: true}, simStartFunc)
}

func simStartFunc(sim models.Sim) {

	s := sim.AddSprite()
	s.AddCostume(ebitensim.DecodeCodedSprite(ebitensim.TurtleImage), "t1")
	s.Costume("t1")
	s.Scale(10)
	s.Z(0)
	s.Visible(true)

	for i := 0; i < 1000; i++ {
		go turtle(sim)
	}

	time.Sleep(time.Second * 6)
	s.Z(2)
}

func turtle(sim models.Sim) {
	s := sim.AddSprite()
	s.AddCostume(ebitensim.DecodeCodedSprite(ebitensim.TurtleImage), "t1")
	s.Costume("t1")
	s.Scale(.2)
	s.Z(1)
	s.Visible(true)

	randomXStep := rand.Float64()*2 - 1
	randomYStep := rand.Float64()*2 - 1
	x := 0.0
	y := 0.0

	for i := 0; i < 500; i++ {
		time.Sleep(time.Millisecond * 10)
		x += randomXStep
		y += randomYStep
		s.Pos(x, y)

	}

	time.Sleep(time.Millisecond * time.Duration(rand.Float64()*10000))
	s.DeleteSprite()
}
