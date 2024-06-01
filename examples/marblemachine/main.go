package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
	"github.com/gary23b/sprites"
	"github.com/gary23b/sprites/spritesmodels"
	"github.com/gary23b/sprites/spritestools"
	"github.com/jakecoffman/cp"
)

// https://github.com/jakecoffman/cp-ebiten/tree/gh-pages
func main() {
	params := sprites.ScratchParams{Width: 1000, Height: 1000, ShowFPS: true}
	sprites.Start(params, simStartFunc)
}

type object struct {
	sprite spritesmodels.Sprite
	body   *cp.Body
	shape  *cp.Shape
}

/*
A Body appears to be what has mass, momentum, and moment of inertia.
A Shape is what touches other things, and is attached to a body.
A shape can be made out of segments, a circle, or a rectangle.

A segment seems to be a wide line with rounded ends.
*/

func simStartFunc(sim spritesmodels.Sim) {
	oList := []*object{}

	space := cp.NewSpace()
	space.SetGravity(cp.Vector{X: 0, Y: -600})

	_ = AddMarbleTrack(sim, space)

	// oList = append(oList, AddContainer(sim, space))

	mass := 1.0
	width := 30.0
	// height := width * 2

	for i := 0; i < 1000; i++ {
		a := NewCircle(sim, space, cp.Vector{X: (rand.Float64()*2 - 1) * 300, Y: rand.Float64() * 300}, mass, width/2)
		oList = append(oList, a)
	}

	// Run the processing loop forever.
	for {
		time.Sleep(time.Millisecond)
		space.Step(.010)

		for _, o := range oList {
			p := o.body.Position()
			o.sprite.Pos(p.X, p.Y)
			o.sprite.Angle(180.0 / math.Pi * o.body.Angle())
		}
	}
}

func AddMarbleElevator(sim spritesmodels.Sim, space *cp.Space) *object {
	container := space.AddBody(cp.NewKinematicBody())
	container.SetAngularVelocity(0.4)
	container.SetPosition(cp.Vector{X: 0, Y: 0})

	a := cp.Vector{X: -200, Y: -200}
	b := cp.Vector{X: -200, Y: 200}
	c := cp.Vector{X: 200, Y: 200}
	d := cp.Vector{X: 200, Y: -200}

	AddWall := func(space *cp.Space, body *cp.Body, a, b cp.Vector, radius float64) {
		shape := cp.NewSegment(body, a, b, radius)
		_ = space.AddShape(shape)
		shape.SetElasticity(1)
		shape.SetFriction(1)
	}

	AddWall(space, container, a, b, 1)
	AddWall(space, container, b, c, 1)
	AddWall(space, container, c, d, 1)
	AddWall(space, container, d, a, 1)

	dc := gg.NewContext(400, 400)
	dc.DrawRectangle(0, 0, 400, 400)
	dc.SetColor(spritestools.Aqua)
	dc.Stroke()
	sim.AddCostume(dc.Image(), "container")
	sprite := sim.AddSprite("container")
	sprite.Costume("container")
	sprite.Pos(0, 0)
	sprite.Visible(true)

	ret := &object{
		sprite: sprite,
		body:   container,
	}
	return ret
}

func AddMarbleTrack(sim spritesmodels.Sim, space *cp.Space) []*object {
	oList := []*object{}

	o := NewStaticBox(sim, space, cp.Vector{X: 0, Y: -490}, 990, 4)
	oList = append(oList, o)
	o = NewStaticBox(sim, space, cp.Vector{X: -490, Y: 0}, 4, 990)
	oList = append(oList, o)
	o = NewStaticBox(sim, space, cp.Vector{X: 490, Y: 0}, 4, 990)
	oList = append(oList, o)
	o = NewStaticBox(sim, space, cp.Vector{X: 0, Y: 490}, 990, 4)
	oList = append(oList, o)

	return oList
}

func NewBox(sim spritesmodels.Sim, space *cp.Space, pos cp.Vector, mass, width, height float64) *object {
	body := cp.NewBody(mass, cp.MomentForBox(mass, width, height))
	_ = space.AddBody(body)
	body.SetPosition(pos)

	shape := cp.NewBox(body, width, height, 0)
	_ = space.AddShape(shape)
	shape.SetElasticity(0)
	shape.SetFriction(0.7)

	c := color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 0xFF}
	costumeName := fmt.Sprintf("%X", rand.Uint64())
	sim.AddCostume(createRectangleImage(width, height, c), costumeName)
	sprite := sim.AddSprite("")
	sprite.Costume(costumeName)
	sprite.Pos(pos.X, pos.Y)
	sprite.Visible(true)

	ret := &object{
		sprite: sprite,
		body:   body,
		shape:  shape,
	}
	return ret
}

func NewStaticBox(sim spritesmodels.Sim, space *cp.Space, pos cp.Vector, width, height float64) *object {
	body := cp.NewStaticBody()
	_ = space.AddBody(body)
	body.SetPosition(pos)

	shape := cp.NewBox(body, width, height, 0)
	_ = space.AddShape(shape)
	shape.SetElasticity(1.0)
	shape.SetFriction(1.0)

	c := color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 0xFF}
	costumeName := fmt.Sprintf("NewStaticBox:%X", rand.Uint64())
	sim.AddCostume(createRectangleImage(width, height, c), costumeName)
	sprite := sim.AddSprite("")
	sprite.Costume(costumeName)
	sprite.Pos(pos.X, pos.Y)
	sprite.Visible(true)

	ret := &object{
		sprite: sprite,
		body:   body,
		shape:  shape,
	}
	return ret
}

func NewCircle(sim spritesmodels.Sim, space *cp.Space, pos cp.Vector, mass, radius float64) *object {
	body := cp.NewBody(mass, cp.MomentForCircle(mass, 0, radius, cp.Vector{}))
	_ = space.AddBody(body)
	body.SetPosition(pos)

	shape := cp.NewCircle(body, radius, cp.Vector{})
	_ = space.AddShape(shape)
	shape.SetElasticity(0)
	shape.SetFriction(0.7)

	c := color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 0xFF}
	costumeName := fmt.Sprintf("%X", rand.Uint64())
	sim.AddCostume(createCircleImage(radius, c), costumeName)
	sprite := sim.AddSprite("")
	sprite.Costume(costumeName)
	sprite.Pos(pos.X, pos.Y)
	sprite.Visible(true)

	ret := &object{
		sprite: sprite,
		body:   body,
		shape:  shape,
	}
	return ret
}
