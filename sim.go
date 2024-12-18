package sprites

import (
	"fmt"
	"image"
	"log"
	"math"
	"math/rand"
	"sync"

	"github.com/gary23b/sprites/game"
	"github.com/gary23b/sprites/spritesmodels"
	"github.com/gary23b/sprites/spritestools"
)

type Sim interface {
	GetWidth() int
	GetHeight() int

	AddCostume(img image.Image, name string)
	AddSprite(UniqueName string) Sprite // If no name is given, a random name is generated.
	DeleteSprite(Sprite)
	DeleteAllSprites()

	SpriteUpdatePosAngle(in Sprite)
	SpriteUpdateFull(in Sprite)

	AddSound(path, name string)
	PlaySound(name string, volume float64) // volume must be between 0 and 1.

	PressedUserInput() *spritesmodels.UserInput
	SubscribeToJustPressedUserInput() chan *spritesmodels.UserInput
	UnSubscribeToJustPressedUserInput(in chan *spritesmodels.UserInput)

	GetSpriteID(UniqueName string) int
	GetSpriteInfo(UniqueName string) spritesmodels.SpriteState
	GetSpriteInfoByID(id int) spritesmodels.SpriteState

	WhoIsNearMe(x, y, distance float64) []spritesmodels.NearMeInfo
	SendMsg(toSpriteID int, msg any)

	GetScreenshot() image.Image

	Exit()
}

type simState struct {
	width   int
	height  int
	g       *game.EbitenGame
	cmdChan chan any

	justPressedBroker *spritestools.Broker[*spritesmodels.UserInput]
	posBroker         *spritestools.PositionBroker

	idToSpriteMapMutex sync.RWMutex
	idToSpriteMap      map[int]Sprite
	nameToSpriteMap    map[string]Sprite
}

var _ Sim = &simState{} // Force the linter to tell us if the interface is implemented

type SimParams struct {
	Width   int  // Window Width in pixels
	Height  int  // Window Height in pixels
	ShowFPS bool // Show Frame-Rate and Update-Rate information in top left corner of window
}

// The drawFunc will be started as a go routine.
func Start(params SimParams, simStartFunc func(Sim)) {
	log.SetFlags(log.Ltime | log.Lmicroseconds | log.Lshortfile)

	ret := &simState{
		width:             params.Width,
		height:            params.Height,
		justPressedBroker: spritestools.NewBroker[*spritesmodels.UserInput](100),
		posBroker:         spritestools.NewPositionBroker(),
		idToSpriteMap:     make(map[int]Sprite),
		nameToSpriteMap:   make(map[string]Sprite),
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

func (s *simState) Exit() {
	s.g.TellGameToExit()
}

func (s *simState) AddSprite(uniqueName string) Sprite {
	spriteID := s.g.GetNextSpriteID()
	if uniqueName == "" {
		uniqueName = fmt.Sprintf("rand%X%X", rand.Uint64(), rand.Uint64())
	}
	update := spritesmodels.CmdAddNewSprite{
		SpriteID: spriteID,
	}
	s.cmdChan <- update

	s.posBroker.AddSprite(spriteID)
	ret := NewSprite(s, uniqueName, spriteID)

	s.idToSpriteMapMutex.Lock()
	s.idToSpriteMap[spriteID] = ret
	s.nameToSpriteMap[uniqueName] = ret
	s.idToSpriteMapMutex.Unlock()

	s.posBroker.UpdateSpriteInfo(spriteID, ret.GetState())
	return ret
}

func (s *simState) DeleteSprite(in Sprite) {
	spriteID := in.GetSpriteID()
	s.posBroker.RemoveSprite(spriteID)
	update := spritesmodels.CmdSpriteDelete{
		SpriteID: spriteID,
	}
	s.cmdChan <- update

	s.idToSpriteMapMutex.Lock()
	delete(s.idToSpriteMap, spriteID)
	delete(s.nameToSpriteMap, in.GetUniqueName())
	s.idToSpriteMapMutex.Unlock()
}

func (s *simState) DeleteAllSprites() {
	update := spritesmodels.CmdSpritesDeleteAll{}
	s.cmdChan <- update
	s.posBroker = spritestools.NewPositionBroker()

	s.idToSpriteMapMutex.Lock()
	s.idToSpriteMap = make(map[int]Sprite)
	s.nameToSpriteMap = make(map[string]Sprite)
	s.idToSpriteMapMutex.Unlock()
}

func (s *simState) SpriteUpdatePosAngle(in Sprite) {
	status := in.GetState()
	s.posBroker.UpdateSpriteInfo(status.SpriteID, status)
	cmd := spritesmodels.CmdSpriteUpdateMin{
		SpriteID: status.SpriteID,
		X:        status.X,
		Y:        status.Y,
		AngleRad: status.AngleDegrees * (math.Pi / 180.0),
	}

	s.cmdChan <- cmd
}

func (s *simState) SpriteUpdateFull(in Sprite) {
	status := in.GetState()
	s.posBroker.UpdateSpriteInfo(status.SpriteID, status)
	cmd := spritesmodels.CmdSpriteUpdateFull{
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

func (s *simState) GetSpriteID(uniqueName string) int {
	s.idToSpriteMapMutex.RLock()
	sprite, ok := s.nameToSpriteMap[uniqueName]
	s.idToSpriteMapMutex.RUnlock()
	if !ok {
		log.Printf("%s, doesn't exist.\n", uniqueName)
		return -1
	}
	return sprite.GetSpriteID()
}

func (s *simState) GetSpriteInfo(uniqueName string) spritesmodels.SpriteState {
	return s.posBroker.GetSpriteInfo(s.GetSpriteID(uniqueName))
}

func (s *simState) GetSpriteInfoByID(id int) spritesmodels.SpriteState {
	return s.posBroker.GetSpriteInfo(id)
}

func (s *simState) GetWidth() int {
	return s.width
}

func (s *simState) GetHeight() int {
	return s.height
}

func (s *simState) PressedUserInput() *spritesmodels.UserInput {
	ret := s.g.PressedUserInput()
	return ret
}

func (s *simState) SubscribeToJustPressedUserInput() chan *spritesmodels.UserInput {
	return s.justPressedBroker.Subscribe()
}

func (s *simState) UnSubscribeToJustPressedUserInput(in chan *spritesmodels.UserInput) {
	s.justPressedBroker.Unsubscribe(in)
}

func (sim *simState) AddCostume(img image.Image, name string) {
	update := spritesmodels.CmdAddCostume{
		Img:         img,
		CostumeName: name,
	}
	sim.cmdChan <- update
}

func (sim *simState) AddSound(path, name string) {
	cmd := spritesmodels.CmdAddSound{
		Path:      path,
		SoundName: name,
	}
	sim.cmdChan <- cmd
}

func (sim *simState) PlaySound(name string, volume float64) {
	cmd := spritesmodels.CmdPlaySound{
		SoundName: name,
		Volume:    volume,
	}
	sim.cmdChan <- cmd
}

func (sim *simState) WhoIsNearMe(x, y, distance float64) []spritesmodels.NearMeInfo {
	return sim.posBroker.GetSpritesNearMe(x, y, distance)
}

func (sim *simState) SendMsg(toSpriteID int, msg any) {
	sim.idToSpriteMapMutex.RLock()
	toSprite, ok := sim.idToSpriteMap[toSpriteID]
	sim.idToSpriteMapMutex.RUnlock()
	if !ok {
		log.Printf("Could not send msg to %d, doesn't exist.\n", toSpriteID)
		return
	}

	toSprite.AddMsg(msg)
}

func (sim *simState) GetScreenshot() image.Image {
	screenshotChan := make(chan image.Image)

	cmd := spritesmodels.CmdGetScreenshot{
		ImageChan: screenshotChan,
	}
	sim.cmdChan <- cmd

	// Now wait for the screenshot to arrive.
	screenshot := <-screenshotChan
	return screenshot
}

// This returns nil if there is no new data.
// This will throw away all but the newest set of data available. So this should be called faster that the game update rate (60Hz),
// otherwise sim.PressedUserInput() should be used instead.
func GetNewestJustPressedFromChan(justPressedChan chan *spritesmodels.UserInput) *spritesmodels.UserInput {
	var ret *spritesmodels.UserInput

ChanExtractionLoop:
	for {
		select {
		case i := <-justPressedChan:
			ret = i
		default:
			// receiving from chan would block
			break ChanExtractionLoop
		}
	}
	return ret
}
