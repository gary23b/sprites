package ebitensim

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"log"
	"sync"

	"github.com/gary23b/sprites/models"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type ebitenSprite struct {
	id         int
	z          int
	arrayIndex int

	spriteImage    []*ebiten.Image
	spriteImageMap map[string]int
	CostumeIndex   int
	width          float64
	height         float64

	x, y           float64
	angle          float64
	visible        bool
	xScale, yScale float64
	opacity        float64

	deleted bool
}

////////////////////////////////

type EbitenGame struct {
	screenWidth  int
	screenHeight int
	showFPS      bool

	controlState SavedControlState
	controls     *models.UserInput

	spritesChan   chan any
	spriteMutex   sync.Mutex
	sprites       [][]*ebitenSprite
	spriteImgs    []ebiten.Image
	spritesImgMap map[string]*ebiten.Image

	nextSpriteIndex int
	spriteMap       map[int]*ebitenSprite
}

func NewGame(width, height int, showFPS bool) *EbitenGame {
	g := &EbitenGame{
		screenWidth:   width,
		screenHeight:  height,
		showFPS:       showFPS,
		spritesChan:   make(chan any, 10000),
		sprites:       make([][]*ebitenSprite, 10),
		spritesImgMap: make(map[string]*ebiten.Image),
		spriteMap:     make(map[int]*ebitenSprite, 10000),
	}

	for i := 0; i < 10; i++ {
		g.sprites[i] = make([]*ebitenSprite, 0, 1000)
	}

	// ebiten.SetTPS(120)
	// ebiten.SetVsyncEnabled(false) // For some reason, on Windows, there is quite a bit of lag.
	// setting this to false clears it up, but also makes it run at 1000Hz...
	ebiten.SetWindowSize(g.screenWidth, g.screenHeight)
	ebiten.SetWindowTitle("Go Turtle Graphics")
	return g
}

func (g *EbitenGame) GetSpriteCmdChannel() chan any {
	return g.spritesChan
}

func (g *EbitenGame) GetNextSpriteID() int {
	g.spriteMutex.Lock()
	defer g.spriteMutex.Unlock()

	newID := g.nextSpriteIndex
	g.nextSpriteIndex++

	return newID
}

func (g *EbitenGame) addSprite(newID int) int {
	newArrayIndex := len(g.sprites[0])
	newSprite := ebitenSprite{
		id:             newID,
		z:              0,
		arrayIndex:     newArrayIndex,
		opacity:        100,
		spriteImageMap: make(map[string]int),
	}

	g.sprites[0] = append(g.sprites[0], &newSprite)
	g.spriteMap[newID] = g.sprites[0][newArrayIndex]

	return newID
}

// Hash the image so we can check if we already have it.
func hashImage(img image.Image) string {
	w := &bytes.Buffer{}
	err := png.Encode(w, img)
	if err != nil {
		panic(err)
	}

	h := sha1.New()
	h.Write(w.Bytes())
	resultBytes := h.Sum(nil)
	resultKey := base64.StdEncoding.EncodeToString(resultBytes)
	return resultKey
}

func (g *EbitenGame) addSpriteCostume(spriteIndex int, img image.Image, costumeName string) int {
	resultKey := hashImage(img)

	s, ok := g.spriteMap[spriteIndex]
	if !ok {
		log.Printf("The given sprite index is not valid: %d\n", spriteIndex)
		return -1
	}

	spriteImage, ok := g.spritesImgMap[resultKey]
	if !ok {
		newSprite := ebiten.NewImageFromImage(img)
		g.spriteImgs = append(g.spriteImgs, *newSprite)
		spriteImage = &g.spriteImgs[len(g.spriteImgs)-1]
		g.spritesImgMap[resultKey] = spriteImage

		fmt.Println("creating a new sprite")
	}

	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	s.spriteImage = append(s.spriteImage, spriteImage)
	s.spriteImageMap[costumeName] = len(s.spriteImage) - 1

	if len(s.spriteImage) == 1 {
		s.width = float64(width)
		s.height = float64(height)
	}
	return len(s.spriteImage) - 1
}

func (g *EbitenGame) deleteSprite(spriteIndex int) {
	s, ok := g.spriteMap[spriteIndex]
	if !ok {
		log.Printf("The given sprite index is not valid: %d\n", spriteIndex)
		return
	}
	delete(g.spriteMap, spriteIndex)

	s.visible = false
	s.deleted = true
}

func (g *EbitenGame) moveSpriteToNewLayer(s *ebitenSprite, newZ int) *ebitenSprite {
	g.sprites[s.z][s.arrayIndex] = nil
	g.sprites[newZ] = append(g.sprites[newZ], s)
	newIndex := len(g.sprites[newZ]) - 1
	g.spriteMap[s.id] = s
	s.arrayIndex = newIndex
	s.z = newZ
	return s
}

func (g *EbitenGame) processSpriteCommands() {
	// Pull all the items out of the command channel until it is empty.
EatSpritesCmdLoop:
	for {
		select {
		case cmd := <-g.spritesChan:
			switch v := cmd.(type) {
			case spriteUpdateMin:
				s, ok := g.spriteMap[v.SpriteIndex]
				if !ok {
					log.Printf("The given sprite index is not valid: %d\n", v.SpriteIndex)
					continue
				}
				costumeID, ok := s.spriteImageMap[v.CostumeName]
				if !ok {
					log.Printf("The given costume name is not valid: %d, %s\n", v.SpriteIndex, v.CostumeName)
					continue
				}
				s.CostumeIndex = costumeID
				s.x = v.X
				s.y = v.Y

			case spriteUpdateFull:
				s, ok := g.spriteMap[v.SpriteIndex]
				if !ok {
					log.Printf("The given sprite index is not valid: %d\n", v.SpriteIndex)
					continue
				}
				if s.z != v.Z {
					s = g.moveSpriteToNewLayer(s, v.Z)
				}

				costumeID, ok := s.spriteImageMap[v.CostumeName]
				if !ok {
					log.Printf("The given costume name is not valid: %d, %s\n", v.SpriteIndex, v.CostumeName)
					continue
				}
				s.CostumeIndex = costumeID
				s.x = v.X
				s.y = v.Y
				s.angle = v.Angle
				s.visible = v.Visible
				s.xScale = v.XScale
				s.yScale = v.YScale
				s.opacity = v.Opacity
			case spriteAddNewSprite:
				g.addSprite(v.SpriteID)
			case spriteAddCostume:
				g.addSpriteCostume(v.SpriteIndex, v.img, v.costumeName)
			case spriteUpdateDelete:
				g.deleteSprite(v.SpriteIndex)

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
	g.controls = g.controlState.GetUserInput()

	g.processSpriteCommands()

	return nil
}

func (g *EbitenGame) Draw(screen *ebiten.Image) {
	op := ebiten.DrawImageOptions{}

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
			if len(sprite.spriteImage) == 0 {
				continue
			}
			op.GeoM.Reset()
			op.ColorScale.Reset()
			op.GeoM.Translate(-sprite.width/2, -sprite.height/2) // Move the center to (0,0) so that we can rotate around the center.
			op.GeoM.Rotate(-sprite.angle)                        // This command rotates clockwise for some reason.
			op.GeoM.Scale(sprite.xScale, sprite.yScale)
			op.GeoM.Translate(float64(g.screenWidth/2), float64(g.screenHeight/2)) // (0,0) is in the center for Cartesian cordinates
			op.GeoM.Translate(sprite.x, -sprite.y)

			if sprite.opacity != 100 {
				op.ColorScale.SetA(float32(sprite.opacity) / 100)
			}

			screen.DrawImage(sprite.spriteImage[sprite.CostumeIndex], &op)

		}
	}

	if g.showFPS {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f, TPS: %0.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
	}
}

func (g *EbitenGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight
}

func (g *EbitenGame) getUserInput() models.UserInput {
	if g == nil || g.controls == nil {
		return models.UserInput{}
	}

	return *g.controls
}
