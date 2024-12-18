package game

import (
	"fmt"
	"image"
	"log"
	"sync"

	"github.com/gary23b/sprites/spritesmodels"
	"github.com/gary23b/sprites/spritestools"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	sampleRate = 48000
)

type ebitenSprite struct {
	id         int // used in idToSpriteMap as the key to point to this struct
	z          int // The current layer. 0-9 allowed
	arrayIndex int // Used for moving a sprite to a new layer

	CostumeIndex int // the index to use to get the current sprite bitmap costume from g.costumes[]

	x, y           float64
	angleRad       float64
	visible        bool
	xScale, yScale float64
	opacity        float64
}

////////////////////////////////

type EbitenGame struct {
	screenWidth       int
	screenHeight      int
	showFPS           bool
	justPressedBroker *spritestools.Broker[*spritesmodels.UserInput]
	exitFlag          bool

	controlState        SavedControlState
	controlsPressed     *spritesmodels.UserInput
	controlsJustPressed *spritesmodels.UserInput

	cmdChan      chan any
	spriteMutex  sync.Mutex // only for protecting nextSpriteID
	nextSpriteID int
	idToSprite   []*ebitenSprite
	sprites      [][]*ebitenSprite // The sprites separated into layers 0 through 9

	costumes           []ebiten.Image
	nameToCostumeIDMap map[string]int

	// Sounds:
	audioContext *audio.Context
	sounds       map[string][]byte

	screenShotRequests []chan image.Image
}

type GameInitStruct struct {
	Width             int
	Height            int
	ShowFPS           bool
	JustPressedBroker *spritestools.Broker[*spritesmodels.UserInput]
}

func NewGame(init GameInitStruct) *EbitenGame {
	g := &EbitenGame{
		screenWidth:       init.Width,
		screenHeight:      init.Height,
		showFPS:           init.ShowFPS,
		justPressedBroker: init.JustPressedBroker,

		cmdChan:      make(chan any, 100000),
		nextSpriteID: 0,
		sprites:      make([][]*ebitenSprite, 10),
		idToSprite:   make([]*ebitenSprite, 0, 31000), // Not sure if this should be an list or map...

		costumes:           make([]ebiten.Image, 0, 1000),
		nameToCostumeIDMap: make(map[string]int),

		audioContext: audio.NewContext(sampleRate),
		sounds:       make(map[string][]byte),
	}

	for i := 0; i < 10; i++ {
		g.sprites[i] = make([]*ebitenSprite, 0, 31000)
	}

	ebiten.SetTPS(120)
	// ebiten.SetVsyncEnabled(false) // For some reason, on Windows, there is quite a bit of lag.
	// setting this to false clears it up, but also makes it run at 1000Hz...
	ebiten.SetWindowSize(g.screenWidth, g.screenHeight)
	ebiten.SetWindowTitle("Go Turtle Graphics")
	return g
}

func (g *EbitenGame) deleteAllSprite() {
	// Deleting everything means just allocating new arrays.
	g.idToSprite = make([]*ebitenSprite, 0, 31000)
	g.sprites = make([][]*ebitenSprite, 10)
	for i := 0; i < 10; i++ {
		g.sprites[i] = make([]*ebitenSprite, 0, 31000)
	}
}

func (g *EbitenGame) GetSpriteCmdChannel() chan any {
	return g.cmdChan
}

func (g *EbitenGame) TellGameToExit() {
	g.exitFlag = true
}

func (g *EbitenGame) GetNextSpriteID() int {
	g.spriteMutex.Lock()
	defer g.spriteMutex.Unlock()

	newID := g.nextSpriteID
	g.nextSpriteID++
	g.idToSprite = append(g.idToSprite, nil)

	return newID
}

func (g *EbitenGame) addSprite(newID int) {
	newArrayIndex := len(g.sprites[0])
	newSprite := ebitenSprite{
		id:           newID,
		z:            0,
		arrayIndex:   newArrayIndex,
		opacity:      100,
		CostumeIndex: -1,
	}

	g.sprites[0] = append(g.sprites[0], &newSprite)
	g.idToSprite[newSprite.id] = g.sprites[0][newSprite.arrayIndex]
}

func (g *EbitenGame) addSpriteCostume(img image.Image, costumeName string) {
	newSprite := ebiten.NewImageFromImage(img)

	// check if we should replace an existing costume:
	id, ok := g.nameToCostumeIDMap[costumeName]
	if ok {
		g.costumes[id] = *newSprite
		// fmt.Printf("Replacing a sprite: %s\n", costumeName)
		return
	}

	g.costumes = append(g.costumes, *newSprite)
	g.nameToCostumeIDMap[costumeName] = len(g.costumes) - 1
	// fmt.Printf("creating a new sprite: %s\n", costumeName)
}

func (g *EbitenGame) deleteSprite(spriteIndex int) {
	s := g.idToSprite[spriteIndex]
	g.idToSprite[spriteIndex] = nil
	g.sprites[s.z][s.arrayIndex] = nil

	s.visible = false
	// Ideally when this function returns, there will be no more refs to the struct, so it will be garbage collected.
}

func (g *EbitenGame) moveSpriteToNewLayer(s *ebitenSprite, newZ int) *ebitenSprite {
	g.sprites[s.z][s.arrayIndex] = nil
	g.sprites[newZ] = append(g.sprites[newZ], s)
	newIndex := len(g.sprites[newZ]) - 1
	g.idToSprite[s.id] = s
	s.arrayIndex = newIndex
	s.z = newZ
	return s
}

func (g *EbitenGame) processSpriteCommands() {
	// Pull all the items out of the command channel until it is empty.
EatSpritesCmdLoop:
	for {
		select {
		case cmd := <-g.cmdChan:
			switch v := cmd.(type) {
			case spritesmodels.CmdSpriteUpdateMin:
				s := g.idToSprite[v.SpriteID]
				s.x = v.X
				s.y = v.Y
				s.angleRad = v.AngleRad

			case spritesmodels.CmdSpriteUpdateFull:
				s := g.idToSprite[v.SpriteID]
				if s.z != v.Z {
					s = g.moveSpriteToNewLayer(s, v.Z)
				}

				costumeID, ok := g.nameToCostumeIDMap[v.CostumeName]
				if !ok {
					log.Printf("The given costume name is not valid: %d, %s\n", v.SpriteID, v.CostumeName)
					continue
				}
				s.CostumeIndex = costumeID
				s.x = v.X
				s.y = v.Y
				s.angleRad = v.Angle
				s.visible = v.Visible
				s.xScale = v.XScale
				s.yScale = v.YScale
				s.opacity = v.Opacity
			case spritesmodels.CmdAddNewSprite:
				g.addSprite(v.SpriteID)
			case spritesmodels.CmdAddCostume:
				g.addSpriteCostume(v.Img, v.CostumeName)
			case spritesmodels.CmdSpriteDelete:
				g.deleteSprite(v.SpriteID)
			case spritesmodels.CmdSpritesDeleteAll:
				g.deleteAllSprite()
			// Sounds
			case spritesmodels.CmdAddSound:
				g.addSound(v.Path, v.SoundName)

			case spritesmodels.CmdPlaySound:
				g.playSound(v.SoundName, v.Volume)

			case spritesmodels.CmdGetScreenshot:
				g.screenShotRequests = append(g.screenShotRequests, v.ImageChan)

			default:
				log.Printf("I don't know about type %T!\n", v)
			}
		default:
			break EatSpritesCmdLoop
		}
	}
}

////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////

// This function will not return. It must be run on the main thread.
func (g *EbitenGame) RunGame() {
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *EbitenGame) Update() error {
	if g.exitFlag {
		return ebiten.Termination
	}

	g.controlsPressed, g.controlsJustPressed = g.controlState.GetUserInput(g.screenWidth, g.screenHeight)
	if g.controlsJustPressed.AnyPressed {
		g.justPressedBroker.Publish(g.controlsJustPressed)
	}

	g.processSpriteCommands()

	return nil
}

func (g *EbitenGame) Draw(screen *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	count := 0
	for i := range g.sprites {
		a := g.sprites[i]
		for j := range a {
			sprite := a[j]
			if sprite == nil {
				continue
			}
			if !sprite.visible {
				continue
			}
			if sprite.CostumeIndex < 0 {
				continue
			}
			op.GeoM.Reset()
			op.ColorScale.Reset()
			costume := g.costumes[sprite.CostumeIndex]
			w, h := costume.Bounds().Dx(), costume.Bounds().Dy()
			op.GeoM.Translate(-float64(w)/2, -float64(h)/2) // Move the center to (0,0) so that we can rotate around the center.
			op.GeoM.Scale(sprite.xScale, sprite.yScale)
			op.GeoM.Rotate(-sprite.angleRad) // This command rotates clockwise for some reason.

			op.GeoM.Translate(float64(g.screenWidth/2), float64(g.screenHeight/2)) // (0,0) is in the center for Cartesian coordinates
			op.GeoM.Translate(sprite.x, -sprite.y)

			if sprite.opacity != 100 {
				op.ColorScale.SetA(float32(sprite.opacity) / 100)
			}

			screen.DrawImage(&costume, &op)
			count++
		}
	}

	if g.showFPS {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f, TPS: %0.2f, Cnt: %d", ebiten.ActualFPS(), ebiten.ActualTPS(), count))
	}

	if len(g.screenShotRequests) > 0 {
		screenshot := image.NewRGBA(screen.Bounds())
		screen.ReadPixels(screenshot.Pix)
		for i := range g.screenShotRequests {
			g.screenShotRequests[i] <- screenshot
		}
		g.screenShotRequests = []chan image.Image{}
	}
}

func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

func (g *EbitenGame) PressedUserInput() *spritesmodels.UserInput {
	if g == nil || g.controlsPressed == nil {
		return &spritesmodels.UserInput{}
	}

	return g.controlsPressed
}
