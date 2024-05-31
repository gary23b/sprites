package main

import (
	"time"

	"github.com/gary23b/sprites/models"
	"github.com/gary23b/sprites/sim"
	"github.com/gary23b/sprites/tools"
)

func main() {
	params := sim.SimParams{
		Width:   1000,
		Height:  1000,
		ShowFPS: true,
	}
	sim.StartSim(params, simStartFunc)
}

func simStartFunc(sim models.Sim) {
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

		select {
		case i := <-justPressedChan:
			// use i
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
		default:
			// receiving from chan would block
		}
		time.Sleep(time.Millisecond * 10)
	}

	sim.UnSubscribeToJustPressedUserInput(justPressedChan)
	s.DeleteSprite()
	sim.Exit()
}
