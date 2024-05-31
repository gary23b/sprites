package sim

import (
	"image"
	"log"
	"math"

	"github.com/gary23b/sprites/game"
	"github.com/gary23b/sprites/models"
	"github.com/gary23b/sprites/sprite"
	"github.com/gary23b/sprites/tools"
)

type simStruct struct {
	width   int
	height  int
	g       *game.EbitenGame
	cmdChan chan any

	justPressedBroker *tools.Broker[*models.UserInput]
	posBroker         *tools.PositionBroker
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
		posBroker:         tools.NewPositionBroker(),
	}

	gameInit := game.GameInitStruct{
		Width:             params.Width,
		Height:            params.Height,
		ShowFPS:           params.ShowFPS,
		JustPressedBroker: ret.justPressedBroker,
	}
	ret.g = game.NewGame(gameInit)
	ret.cmdChan = ret.g.GetSpriteCmdChannel()
	go simStartFunc(ret)
	ret.g.RunGame()
}

func (s *simStruct) Exit() {
	s.g.TellGameToExit()
}

func (s *simStruct) AddSprite(UniqueName string) models.Sprite {
	spriteID := s.g.GetNextSpriteID()
	update := models.CmdAddNewSprite{
		SpriteID: spriteID,
	}
	s.cmdChan <- update

	s.posBroker.AddSprite(UniqueName)
	ret := sprite.NewSprite(s, UniqueName, spriteID)
	s.posBroker.UpdateSpriteInfo(UniqueName, ret.GetState())
	return ret
}

func (s *simStruct) DeleteSprite(in models.Sprite) {
	s.posBroker.RemoveSprite(in.GetUniqueName())
	update := models.CmdSpriteDelete{
		SpriteID: in.GetSpriteID(),
	}
	s.cmdChan <- update
}

func (s *simStruct) DeleteAllSprites() {
	update := models.CmdSpritesDeleteAll{}
	s.cmdChan <- update
	s.posBroker = tools.NewPositionBroker()
}

func (s *simStruct) SpriteUpdatePosAngle(in models.Sprite) {
	status := in.GetState()
	s.posBroker.UpdateSpriteInfo(status.UniqueName, status)
	cmd := models.CmdSpriteUpdateMin{
		SpriteID: status.SpriteID,
		X:        status.X,
		Y:        status.Y,
		AngleRad: status.AngleDegrees * (math.Pi / 180.0),
	}

	s.cmdChan <- cmd
}

func (s *simStruct) SpriteUpdateFull(in models.Sprite) {
	status := in.GetState()
	s.posBroker.UpdateSpriteInfo(status.UniqueName, status)
	cmd := models.CmdSpriteUpdateFull{
		SpriteID:    status.SpriteID,
		CostumeName: status.CostumeName,
		X:           status.X,
		Y:           status.Y,
		Z:           status.Z,
		Angle:       status.AngleDegrees * (math.Pi / 180.0),
		Visible:     status.Visible,
		XScale:      status.ScaleX,
		YScale:      status.ScaleY,
		Opacity:     status.Opacity,
	}

	s.cmdChan <- cmd
}

func (s *simStruct) GetSpriteInfo(UniqueName string) models.SpriteState {
	return s.posBroker.GetSpriteInfo(UniqueName)
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
