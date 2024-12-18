package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gary23b/sprites"
	"github.com/gary23b/sprites/spritestools"
)

func main() {
	params := sprites.SimParams{Width: 1000, Height: 1000, ShowFPS: true}
	sprites.Start(params, simStartFunc)
}

// func fileExists(path string) bool {
// 	_, err := os.Stat(path)
// 	return !errors.Is(err, os.ErrNotExist)
// }

func simStartFunc(sim sprites.Sim) {
	sim.AddCostume(sprites.DecodeCodedSprite(sprites.TurtleImage), "t1")

	// if fileExists("jab.wav") {
	// 	sim.AddSound("jab.wav", "jab")
	// 	sim.AddSound("jump.ogg", "jump")
	// 	sim.AddSound("ragtime.mp3", "ragtime")
	// } else {
	// 	sim.AddSound("./examples/play1/jab.wav", "jab")
	// 	sim.AddSound("./examples/play1/jump.ogg", "jump")
	// 	sim.AddSound("./examples/play1/ragtime.mp3", "ragtime")
	// }

	// sim.PlaySound("ragtime", .1)

	testScene(sim)
}

func testScene(sim sprites.Sim) {

	broker := spritestools.NewBroker[string](100)

	s := sim.AddSprite("mainTurtle")
	s.Costume("t1")
	s.Scale(10)
	s.Z(0)
	s.Visible(true)

	// time.Sleep(time.Millisecond * 5000)
	// sim.PlaySound("jab", .001)
	// time.Sleep(time.Millisecond * 200)
	// sim.PlaySound("jab", .01)
	// time.Sleep(time.Millisecond * 200)
	// sim.PlaySound("jab", .1)
	// time.Sleep(time.Millisecond * 200)
	// sim.PlaySound("jab", .4)
	// time.Sleep(time.Millisecond * 200)
	// sim.PlaySound("jab", .5)
	// time.Sleep(time.Millisecond * 200)
	// sim.PlaySound("jab", .6)
	// time.Sleep(time.Millisecond * 200)
	// sim.PlaySound("jab", .7)
	// time.Sleep(time.Millisecond * 200)
	// sim.PlaySound("jab", .8)
	// time.Sleep(time.Millisecond * 200)
	// sim.PlaySound("jab", .9)
	// time.Sleep(time.Millisecond * 200)
	// sim.PlaySound("jab", 1)
	for i := 0; i < 5000; i++ {
		go turtle(sim, broker)
	}

	time.Sleep(time.Second * 10)
	s.XYScale(-10, 10)

	time.Sleep(time.Millisecond * 5000)

	s.Z(2)
	s.Opacity(30)
	// s.Angle(45)

	// sim.PlaySound("jump", .6)
	broker.Publish("delete all")

	time.Sleep(time.Second * 10)
	sim.DeleteAllSprites()
	// sim.PlaySound("jab", .5)
	time.Sleep(time.Millisecond * 20)

	broker.Stop()

	// testScene(sim)
}

func turtle(sim sprites.Sim, broker *spritestools.Broker[string]) {
	broadcasts := broker.Subscribe()
	s := sim.AddSprite(fmt.Sprintf("turtle%d%d", rand.Uint64(), rand.Uint64()))
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

	// wait for msg
	<-broadcasts

	broker.Unsubscribe(broadcasts)

	time.Sleep(time.Millisecond * time.Duration(rand.Float64()*3000))
	s.DeleteSprite()
}
