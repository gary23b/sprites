package ebitensim

import (
	"image"
	"log"

	"github.com/gary23b/sprites/models"
	"github.com/gary23b/sprites/tools"
)

type simStruct struct {
	width   int
	height  int
	g       *EbitenGame
	cmdChan chan any

	justPressedBroker *tools.Broker[*models.UserInput]
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
		width:             params.Width,
		height:            params.Height,
		justPressedBroker: tools.NewBroker[*models.UserInput](),
	}

	gameInit := GameInitStruct{
		width:             params.Width,
		height:            params.Height,
		showFPS:           params.ShowFPS,
		justPressedBroker: ret.justPressedBroker,
	}
	ret.g = NewGame(gameInit)
	ret.cmdChan = ret.g.GetSpriteCmdChannel()
	go simStartFunc(ret)
	ret.g.RunGame()
}

func (s *simStruct) Exit() {
	s.g.exitFlag = true
}

func (s *simStruct) AddSprite() models.Sprite {
	spriteID := s.g.GetNextSpriteID()
	update := models.CmdAddNewSprite{
		SpriteID: spriteID,
	}
	s.cmdChan <- update

	ret := newSprite(s, spriteID)
	return ret
}

func (s *simStruct) DeleteSprite(in models.Sprite) {
	update := models.CmdSpriteDelete{
		SpriteIndex: in.GetSpriteID(),
	}
	s.cmdChan <- update
}

func (s *simStruct) DeleteAllSprites() {
	update := models.CmdSpritesDeleteAll{}
	s.cmdChan <- update
}

func (s *simStruct) SpriteMinUpdate(in *models.CmdSpriteUpdateMin) {
	s.cmdChan <- *in
}

func (s *simStruct) SpriteFullUpdate(in *models.CmdSpriteUpdateFull) {
	s.cmdChan <- *in
}

func (s *simStruct) GetWidth() int {
	return s.width
}

func (s *simStruct) GetHeight() int {
	return s.height
}

func (s *simStruct) PressedUserInput() *models.UserInput {
	ret := s.g.PressedUserInput()
	return ret
}

func (s *simStruct) SubscribeToJustPressedUserInput() chan *models.UserInput {
	return s.justPressedBroker.Subscribe()
}

func (s *simStruct) UnSubscribeToJustPressedUserInput(in chan *models.UserInput) {
	s.justPressedBroker.Unsubscribe(in)
}

func (sim *simStruct) AddCostume(img image.Image, name string) {
	update := models.CmdAddCostume{
		Img:         img,
		CostumeName: name,
	}
	sim.cmdChan <- update
}

func (sim *simStruct) AddSound(path, name string) {
	cmd := models.CmdAddSound{
		Path:      path,
		SoundName: name,
	}
	sim.cmdChan <- cmd
}

func (sim *simStruct) PlaySound(name string, volume float64) {
	cmd := models.CmdPlaySound{
		SoundName: name,
		Volume:    volume,
	}
	sim.cmdChan <- cmd
}
