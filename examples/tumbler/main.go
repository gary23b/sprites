package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
	"github.com/gary23b/sprites"
	"github.com/gary23b/sprites/models"
	"github.com/gary23b/sprites/tools"
	"github.com/jakecoffman/cp"
)

// https://github.com/jakecoffman/cp-ebiten/tree/gh-pages
func main() {
	params := sprites.ScratchParams{Width: 600, Height: 600, ShowFPS: true}
	sprites.Start(params, simStartFunc)
}

type object struct {
	sprite models.Sprite
	body   *cp.Body
}

/*
A Body appears to be what has mass, momentum, and moment of inertia.
A Shape is what touches other things, and is attached to a body.
A shape can be made out of segments, a circle, or a rectangle.

A segment seems to be a wide line with rounded ends.
*/

func simStartFunc(sim models.Scratch) {
	oList := []*object{}

	space := cp.NewSpace()
	space.SetGravity(cp.Vector{X: 0, Y: -600})

	oList = append(oList, AddContainer(sim, space))

	mass := 1.0
	width := 30.0
	height := width * 2

	for i := 0; i < 7; i++ {
		for j := 0; j < 3; j++ {
			pos := cp.Vector{X: float64(i) * width, Y: float64(j) * height}

			switch rand.Intn(3) {
			case 0:
				a := NewBox(sim, space, pos, mass, width, height)
				oList = append(oList, a)
			case 1:
				a := AddSegment(sim, space, pos, mass, width, height)
				oList = append(oList, a)
			case 2:
				a := NewCircle(sim, space, pos.Add(cp.Vector{X: 0, Y: (height - width) / 2}), mass, width/2)
				b := NewCircle(sim, space, pos.Add(cp.Vector{X: 0, Y: (width - height) / 2}), mass, width/2)
				oList = append(oList, a, b)
			}
		}
	}

	// go scratch.CreateGif(sim, time.Millisecond*100, time.Millisecond*100, "./examples/tumbler/tumbler.gif", 100)

	// Run the processing loop forever.
	for {
		time.Sleep(time.Millisecond * 10)
		space.Step(.010)

		for _, o := range oList {
			p := o.body.Position()
			o.sprite.Pos(p.X, p.Y)
			o.sprite.Angle(180.0 / math.Pi * o.body.Angle())
		}
	}
}

func AddContainer(sim models.Scratch, space *cp.Space) *object {
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
	dc.SetColor(tools.Aqua)
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

func NewBox(sim models.Scratch, space *cp.Space, pos cp.Vector, mass, width, height float64) *object {
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
	}
	return ret
}

func NewCircle(sim models.Scratch, space *cp.Space, pos cp.Vector, mass, radius float64) *object {
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
	}
	return ret
}

func AddSegment(sim models.Scratch, space *cp.Space, pos cp.Vector, mass, width, height float64) *object {
	body := cp.NewBody(mass, cp.MomentForBox(mass, width, height))
	_ = space.AddBody(body)
	body.SetPosition(pos)

	a, b := cp.Vector{X: 0, Y: (height - width) / 2.0}, cp.Vector{X: 0, Y: (width - height) / 2.0}
	shape := cp.NewSegment(body, a, b, width/2.0)
	_ = space.AddShape(shape)
	shape.SetElasticity(0)
	shape.SetFriction(0.7)

	c := color.RGBA{R: uint8(rand.Intn(256)), G: uint8(rand.Intn(256)), B: uint8(rand.Intn(256)), A: 0xFF}
	costumeName := fmt.Sprintf("%X", rand.Uint64())
	sim.AddCostume(createSegmentImage(width, height, c), costumeName)
	sprite := sim.AddSprite("")
	sprite.Costume(costumeName)
	sprite.Pos(pos.X, pos.Y)
	sprite.Visible(true)

	ret := &object{
		sprite: sprite,
		body:   body,
	}
	return ret
}

func createCircleImage(radius float64, c color.Color) image.Image {
	dc := gg.NewContext(int(radius*3), int(radius*3))
	dc.DrawCircle(radius*1.5, radius*1.5, radius)
	dc.SetColor(c)
	dc.Fill()
	dc.SetColor(tools.White)
	dc.SetLineWidth(3)
	dc.DrawLine(radius*1.5, radius*1.5, radius*2.0, radius*1.5)
	dc.Stroke()

	return dc.Image()
}

func createRectangleImage(width, height float64, c color.Color) image.Image {
	dc := gg.NewContext(int(width), int(height))
	dc.SetColor(c)
	dc.DrawRectangle(0, 0, width, height)
	dc.Fill()
	return dc.Image()
}

func createSegmentImage(width, height float64, c color.Color) image.Image {
	hw := width / 2

	dc := gg.NewContext(int(width), int(height))
	dc.SetColor(c)

	dc.MoveTo(0, hw)
	dc.DrawArc(hw, hw, hw, 1.0*math.Pi, 2*math.Pi)
	dc.LineTo(width, height-hw)
	dc.DrawArc(hw, height-hw, hw, 0*math.Pi, 1*math.Pi)
	dc.LineTo(0, hw)
	dc.FillPreserve()
	dc.Stroke()

	return dc.Image()
}
