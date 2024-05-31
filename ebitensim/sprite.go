package ebitensim

import (
	"bytes"
	"fmt"
	"image"
	"log"
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
	// width   int
	// height  int
	x, y    float64
	z       int
	angle   float64
	visible bool
	opacity float64
	scaleX  float64
	scaleY  float64

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

func (s *sprite) AddCostume(img image.Image, name string) {
	_, ok := s.costumeMap[name]
	if ok {
		return
	}

	if s.deleted {
		log.Printf("Error: sprite %d is deleted but being updated\n", s.spriteID)
		return
	}

	update := spriteAddCostume{
		SpriteIndex: s.spriteID,
		img:         img,
		costumeName: name,
	}
	s.sim.cmdChan <- update
}

// Updates
func (s *sprite) Costume(name string) {
	s.constumeName = name
	s.minUpdate()
}

func (s *sprite) Angle(radianAngle float64) {
	s.angle = radianAngle
	s.fullUpdate()
}

func (s *sprite) Pos(cartX, cartY float64) {
	s.x = cartX
	s.y = cartY
	s.minUpdate()
}

func (s *sprite) Z(z int) {
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
	s.constumeName = in.CostumeName
	s.x = in.X
	s.y = in.Y
	s.z = in.Z
	s.angle = in.Angle
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
		Angle:       s.angle,
		Visible:     s.visible,
		ScaleX:      s.scaleX,
		ScaleY:      s.scaleY,
		Opacity:     s.opacity,
	}
}

// exit
func (s *sprite) DeleteSprite() {
	if s.deleted {
		log.Printf("Error: sprite %d is deleted but being deleted again\n", s.spriteID)
		return
	}

	update := spriteUpdateDelete{
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

type spriteUpdateDelete struct {
	SpriteIndex int
}

type spriteAddCostume struct {
	SpriteIndex int
	costumeName string
	img         image.Image
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
		Angle:       s.angle,
		Visible:     s.visible,
		XScale:      s.scaleX,
		YScale:      s.scaleY,
		Opacity:     s.opacity,
	}
	s.sim.cmdChan <- update
}
