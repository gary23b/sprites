package main

import (
	"math/rand"
	"time"

	"github.com/gary23b/sprites/ebitensim"
	"github.com/gary23b/sprites/models"
)

func main() {
	params := ebitensim.SimParams{
		Width:   1000,
		Height:  1000,
		ShowFPS: true,
	}
	ebitensim.StartSim(params, simStartFunc)
}

func simStartFunc(sim models.Sim) {
	sim.AddCostume(ebitensim.DecodeCodedSprite(ebitensim.TurtleImage), "t1")
	testScene(sim)
}

func testScene(sim models.Sim) {

	s := sim.AddSprite()
	s.Costume("t1")
	s.Scale(10)
	s.Z(0)
	s.Visible(true)
	time.Sleep(time.Millisecond * 20)

	for i := 0; i < 30000; i++ {
		go turtle(sim)
	}

	time.Sleep(time.Second * 6)
	s.XYScale(-10, 10)

	time.Sleep(time.Millisecond * 500)

	s.Z(2)
	s.Opacity(30)
	// s.Angle(45)

	time.Sleep(time.Second * 10)
	sim.DeleteAllSprites()
	time.Sleep(time.Millisecond * 20)

	testScene(sim)
}

func turtle(sim models.Sim) {
	s := sim.AddSprite()
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

	//time.Sleep(time.Millisecond * time.Duration(rand.Float64()*10000))
	//s.DeleteSprite()
}
