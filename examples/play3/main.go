package main

import (
	"time"

	"github.com/gary23b/sprites"
	"github.com/gary23b/sprites/models"
	"github.com/gary23b/sprites/tools"
)

func main() {
	params := sprites.ScratchParams{
		Width:   1000,
		Height:  1000,
		ShowFPS: true,
	}
	sprites.Start(params, simStartFunc)
}

func simStartFunc(sim models.Scratch) {
	// sim.AddCostume(ebitensim.DecodeCodedSprite(ebitensim.TurtleImage), "t1")

	textImage := tools.CreateTextBubble("abasdfasdfadsfasdfadsf\nc234", 16)
	sim.AddCostume(textImage, "t1")

	a := 0.0

	s := sim.AddSprite("mainTurtle")
	b := s.GetClickBody()
	s.Costume("t1")
	s.Scale(1)

	s.Pos(0, 0)
	b.AddCirleBody(32, 0, 32)
	b.AddRectangleBody(-300, 0, -5, 5)
	// body.AddCirleBody(-32, 0, 32)
	s.Z(0)
	s.Visible(true)

	justPressedChan := sim.SubscribeToJustPressedUserInput()

MainSpriteLoop:
	for {
		i := sprites.GetNewestJustPressedFromChan(justPressedChan)
		if i.Mouse.Left {
			if b.IsMouseClickInBody(float64(i.Mouse.MouseX), float64(i.Mouse.MouseY)) {
				a += 10
				s.Angle(a)
				s.Pos(a, 0)
			}

		}
		if i.Mouse.Right {
			break MainSpriteLoop
		}
		time.Sleep(time.Millisecond * 10)
	}

	sim.UnSubscribeToJustPressedUserInput(justPressedChan)
	s.DeleteSprite()
	sim.Exit()
}
