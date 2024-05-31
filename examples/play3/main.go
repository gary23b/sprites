package main

import (
	"fmt"
	"time"

	"github.com/gary23b/sprites"
	"github.com/gary23b/sprites/models"
	"github.com/gary23b/sprites/tools"
)

func main() {
	params := sprites.ScratchParams{Width: 1000, Height: 1000, ShowFPS: true}
	sprites.Start(params, simStartFunc)
}

func simStartFunc(sim models.Scratch) {
	// sim.AddCostume(ebitensim.DecodeCodedSprite(ebitensim.TurtleImage), "t1")

	textImage := tools.CreateTextBubble(200, 100, "abasd sdf sdfsdfsd fs dfsdfsd fsdf sf\n    c234", 20)
	sim.AddCostume(textImage, "textBubble")

	a := 0.0

	s := sim.AddSprite("mainTurtle")
	b := s.GetClickBody()
	s.Costume("textBubble")
	s.Scale(1)

	s.Pos(0, 0)
	b.AddCircleBody(32, 0, 32)
	b.AddRectangleBody(-300, 0, -5, 5)
	// body.AddCircleBody(-32, 0, 32)
	s.Z(0)
	s.Visible(true)

	justPressedChan := sim.SubscribeToJustPressedUserInput()

MainSpriteLoop:
	for {
		i := sprites.GetNewestJustPressedFromChan(justPressedChan)
		if i != nil {
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
			if i.Keys.N {
				textImage := tools.CreateTextBubble(200, 100, fmt.Sprintf("angle=%.0f", a), 30)
				sim.AddCostume(textImage, "textBubble")
			}
		}
		time.Sleep(time.Millisecond * 10)
	}

	sim.UnSubscribeToJustPressedUserInput(justPressedChan)
	s.DeleteSprite()
	sim.Exit()
}
