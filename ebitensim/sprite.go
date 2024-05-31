package ebitensim

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"math"
	"os"

	"github.com/gary23b/sprites/models"
)

func LoadSpriteFile(path string) (image.Image, error) {
	spriteFileData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read sprite file: %s, %w", path, err)
	}
	img, _, err := image.Decode(bytes.NewReader(spriteFileData))
	if err != nil {
		return nil, fmt.Errorf("Failed to decode image data: %s, %w", path, err)
	}
	return img, nil
}

type sprite struct {
	sim *simStruct

	spriteID     int
	constumeName string
	x, y         float64
	z            int
	angleRad     float64
	visible      bool
	opacity      float64
	scaleX       float64
	scaleY       float64

	deleted bool

	costumeMap map[string]int
}

var _ models.Sprite = &sprite{}

func newSprite(sim *simStruct) *sprite {
	spriteID := sim.g.GetNextSpriteID()

	ret := &sprite{
		spriteID:   spriteID,
		opacity:    100,
		scaleX:     1,
		scaleY:     1,
		sim:        sim,
		costumeMap: make(map[string]int),
	}

	update := spriteAddNewSprite{
		SpriteID: ret.spriteID,
	}
	sim.cmdChan <- update

	return ret
}

// Updates
func (s *sprite) Costume(name string) {
	s.constumeName = name
	s.minUpdate()
}

func (s *sprite) Angle(angleDegrees float64) {
	s.angleRad = angleDegrees * (math.Pi / 180.0)
	s.fullUpdate()
}

func (s *sprite) Pos(cartX, cartY float64) {
	s.x = cartX
	s.y = cartY
	s.minUpdate()
}

func (s *sprite) Z(z int) {
	if z < 0 || z > 9 {
		log.Println("Z must be from 0 to 9")
		return
	}

	s.z = z
	s.fullUpdate()
}

func (s *sprite) Visible(visible bool) {
	s.visible = visible
	s.fullUpdate()
}

func (s *sprite) Scale(scale float64) {
	s.scaleX = scale
	s.scaleY = scale
	s.fullUpdate()
}

func (s *sprite) XYScale(xScale, yScale float64) {
	s.scaleX = xScale
	s.scaleY = yScale
	s.fullUpdate()
}

func (s *sprite) Opacity(opacityPercent float64) {
	s.opacity = opacityPercent
	s.fullUpdate()
}

func (s *sprite) All(in models.SpriteState) {
	if in.Z < 0 || in.Z > 9 {
		log.Println("Z must be from 0 to 9")
		return
	}

	s.constumeName = in.CostumeName
	s.x = in.X
	s.y = in.Y
	s.z = in.Z
	s.angleRad = in.Angle
	s.visible = in.Visible
	s.opacity = in.Opacity
	s.scaleX = in.ScaleX
	s.scaleY = in.ScaleY

	s.fullUpdate()
}

func (s *sprite) GetState() models.SpriteState {
	return models.SpriteState{
		CostumeName: s.constumeName,
		X:           s.x,
		Y:           s.y,
		Z:           s.z,
		Angle:       s.angleRad,
		Visible:     s.visible,
		ScaleX:      s.scaleX,
		ScaleY:      s.scaleY,
		Opacity:     s.opacity,
	}
}

func (s *sprite) DeleteSprite() {
	if s.deleted {
		log.Printf("Error: sprite %d is deleted but being deleted again\n", s.spriteID)
		return
	}

	update := spriteCmdDelete{
		SpriteIndex: s.spriteID,
	}
	s.sim.cmdChan <- update
}

///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////

type spriteAddNewSprite struct {
	SpriteID int
}

type spriteUpdateMin struct {
	SpriteIndex int
	CostumeName string
	X           float64
	Y           float64
}

type spriteUpdateFull struct {
	SpriteIndex int
	CostumeName string
	X           float64
	Y           float64
	Z           int
	Angle       float64
	Visible     bool
	XScale      float64
	YScale      float64
	Opacity     float64
}

type spriteCmdDelete struct {
	SpriteIndex int
}

type spriteCmdDeleteAll struct {
	SpriteIndex int
}

type spriteAddCostume struct {
	costumeName string
	img         image.Image
}

type cmdAddSound struct {
	path      string
	soundName string
}

type cmdPlaySound struct {
	soundName string
	volume    float64 // between 0 and 1.
}

func (s *sprite) minUpdate() {
	if s.deleted {
		log.Printf("Error: sprite %d is deleted but being updated\n", s.spriteID)
		return
	}

	update := spriteUpdateMin{
		SpriteIndex: s.spriteID,
		CostumeName: s.constumeName,
		X:           s.x,
		Y:           s.y,
	}
	s.sim.cmdChan <- update
}

func (s *sprite) fullUpdate() {
	if s.deleted {
		log.Printf("Error: sprite %d is deleted but being updated\n", s.spriteID)
		return
	}

	update := spriteUpdateFull{
		SpriteIndex: s.spriteID,
		CostumeName: s.constumeName,
		X:           s.x,
		Y:           s.y,
		Z:           s.z,
		Angle:       s.angleRad,
		Visible:     s.visible,
		XScale:      s.scaleX,
		YScale:      s.scaleY,
		Opacity:     s.opacity,
	}
	s.sim.cmdChan <- update
}
