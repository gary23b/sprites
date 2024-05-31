package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/gary23b/sprites"
	"github.com/gary23b/sprites/models"
	"github.com/gary23b/sprites/sprite"
	"github.com/gary23b/sprites/spritestools"
)

var SkyBlue color.RGBA = color.RGBA{0x87, 0xCE, 0xEB, 0xFF}

func main() {
	params := sprites.ScratchParams{Width: 1000, Height: 1000, ShowFPS: true}
	sprites.Start(params, simStartFunc)
}

func simStartFunc(sim models.Scratch) {
	go RunMouse(sim)

	// Set the background by making a single pixel image and scaling it be be the entire screen.
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, SkyBlue)
	sim.AddCostume(img, "skyBackground")
	s := sim.AddSprite("background")
	s.Costume("skyBackground")
	s.Z(0)
	s.Scale(1010)
	s.Visible(true)

	sim.AddCostume(sprite.DecodeCodedSprite(sprite.TurtleImage), "t")

	s = sim.AddSprite("mainTurtle")
	b := s.GetClickBody()
	s.Costume("t")
	s.Scale(1)
	s.Angle(-90)

	s.Pos(0, 0)
	b.AddCircleBody(0, 0, 32)
	s.Z(1)
	// s.Visible(true)

	// MainSpriteLoop:
	for {
		sClone := s.Clone(fmt.Sprintf("%d%d", rand.Uint64(), rand.Uint64()))
		go RunClone(sim, sClone)

		time.Sleep(time.Millisecond * time.Duration(rand.Float64()*3000))

		if GameState.GameOver {
			break
		}
	}

}

type GameStateStruct struct {
	GameOver bool
	Score    int
}

var GameState GameStateStruct

func RunClone(sim models.Scratch, s models.Sprite) {
	justPressedChan := sim.SubscribeToJustPressedUserInput()
	cb := s.GetClickBody()

	randomX := rand.Intn(1000) - 500
	s.Pos(float64(randomX), 550)
	s.Visible(true)
	speed := rand.Float64() * 2
	y := 550.0

	// MainSpriteLoop:
	for {
		justPressed := sprites.GetNewestJustPressedFromChan(justPressedChan)
		if justPressed != nil && justPressed.Mouse.Left {
			if cb.IsMouseClickInBody(float64(justPressed.Mouse.MouseX), float64(justPressed.Mouse.MouseY)) {
				s.DeleteSprite()
				sim.UnSubscribeToJustPressedUserInput(justPressedChan)
				GameState.Score++
				fmt.Println(GameState.Score)
				break
			}
		}

		if y <= -500 {
			GameState.GameOver = true
			fmt.Println("Game Over")
			ShowGameOver(sim)
		}

		if GameState.GameOver {
			break
		}

		y -= speed
		s.Pos(float64(randomX), y)

		time.Sleep(time.Millisecond * 10)
	}
}

// Test the WhoIsNearMe functionality.
func RunMouse(sim models.Scratch) {
	for {
		posInfo := sim.PressedUserInput()
		nearMeList := sim.WhoIsNearMe(float64(posInfo.Mouse.MouseX), float64(posInfo.Mouse.MouseY), 50)

		fmt.Printf("%v\n", nearMeList)
		time.Sleep(time.Millisecond * 100)
	}
}

func ShowGameOver(sim models.Scratch) {
	imgGameOver := spritestools.CreateTextImg(fmt.Sprintf("GAME OVER\nSCORE: %d", GameState.Score), 1000, 130, 60, spritestools.Red)
	sim.AddCostume(imgGameOver, "GameOver")
	s := sim.AddSprite("GameOver")
	s.Costume("GameOver")
	s.Z(9)

	s.Pos(0, 0)
	s.Visible(true)
}
