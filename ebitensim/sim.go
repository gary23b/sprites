package ebitensim

import (
	"log"

	"github.com/gary23b/sprites/models"
)

type simStruct struct {
	width   int
	height  int
	g       *EbitenGame
	cmdChan chan any
}

var _ models.Sim = &simStruct{} // Force the linter to tell us if the interface is implemented

type SimParams struct {
	Width   int
	Height  int
	ShowFPS bool
}

// The drawFunc will be started as a go routine.
func StartSim(params SimParams, simStartFunc func(models.Sim)) {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)

	ret := &simStruct{
		width:  params.Width,
		height: params.Height,
	}

	ret.g = NewGame(params.Width, params.Height, params.ShowFPS)
	ret.cmdChan = ret.g.GetSpriteCmdChannel()
	go simStartFunc(ret)
	ret.g.RunGame()
}

func (s *simStruct) AddSprite() models.Sprite {
	ret := newSprite(s)
	return ret
}

func (sim *simStruct) DeleteAllSprites() {
	update := spriteCmdDeleteAll{}
	sim.cmdChan <- update
}

func (s *simStruct) GetWidth() int {
	return s.width
}

func (s *simStruct) GetHeight() int {
	return s.height
}

func (s *simStruct) GetUserInput() models.UserInput {
	ret := s.g.getUserInput()

	// translate game space to turtle space
	ret.MouseX -= s.width / 2
	ret.MouseY = -ret.MouseY + s.height/2

	return ret
}
