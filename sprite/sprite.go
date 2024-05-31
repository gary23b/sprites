package sprite

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"math"
	"os"

	"github.com/gary23b/sprites/models"
	"github.com/gary23b/sprites/tools"
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
	sim models.Scratch

	spriteID    int
	spriteType  int
	UniqueName  string
	costumeName string
	x, y        float64
	z           int
	angleRad    float64
	visible     bool
	opacity     float64
	scaleX      float64
	scaleY      float64

	deleted bool

	clickBody     models.ClickOnBody
	userInputChan chan *models.UserInput
	receivedMsgs  chan any
}

var _ models.Sprite = &sprite{}

func NewSprite(sim models.Scratch, UniqueName string, spriteID int) *sprite {
	ret := &sprite{
		spriteID:     spriteID,
		UniqueName:   UniqueName,
		opacity:      100,
		scaleX:       1,
		scaleY:       1,
		sim:          sim,
		clickBody:    tools.NewTouchCollisionBody(),
		receivedMsgs: make(chan any, 10),
	}
	return ret
}

func (s *sprite) Clone(UniqueName string) models.Sprite {
	sClone := s.sim.AddSprite(UniqueName)
	sClone.All(s.GetState())
	if s.clickBody != nil {
		sClone.ReplaceClickBody(s.clickBody.Clone())
	}
	return sClone
}

func (s *sprite) GetSpriteID() int {
	return s.spriteID
}

func (s *sprite) GetUniqueName() string {
	return s.UniqueName
}

// Updates
func (s *sprite) Costume(name string) {
	s.costumeName = name
	s.fullUpdate()
}

func (s *sprite) SetType(newType int) {
	s.spriteType = newType
	s.minUpdate()
}

func (s *sprite) Angle(angleDegrees float64) {
	s.angleRad = angleDegrees * (math.Pi / 180.0)
	s.minUpdate()

	s.clickBody.Angle(s.angleRad)
}

func (s *sprite) Pos(cartX, cartY float64) {
	s.x = cartX
	s.y = cartY
	s.minUpdate()

	s.clickBody.Pos(s.x, s.y)
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

	s.spriteType = in.SpriteType
	s.costumeName = in.CostumeName
	s.x = in.X
	s.y = in.Y
	s.z = in.Z
	s.angleRad = in.AngleDegrees * (math.Pi / 180.0)
	s.visible = in.Visible
	s.opacity = in.Opacity
	s.scaleX = in.ScaleX
	s.scaleY = in.ScaleY

	s.fullUpdate()

	s.clickBody.Pos(s.x, s.y)
	s.clickBody.Angle(s.angleRad)
}

func (s *sprite) GetState() models.SpriteState {
	return models.SpriteState{
		SpriteID:     s.spriteID,
		SpriteType:   s.spriteType,
		UniqueName:   s.UniqueName,
		CostumeName:  s.costumeName,
		X:            s.x,
		Y:            s.y,
		Z:            s.z,
		AngleDegrees: s.angleRad * (180.0 / math.Pi),
		Visible:      s.visible,
		ScaleX:       s.scaleX,
		ScaleY:       s.scaleY,
		Opacity:      s.opacity,
		Deleted:      s.deleted,
	}
}

func (s *sprite) DeleteSprite() {
	if s.deleted {
		log.Printf("Error: sprite %d is deleted but being deleted again\n", s.spriteID)
		return
	}

	s.sim.DeleteSprite(s)
}

func (s *sprite) GetClickBody() models.ClickOnBody {
	return s.clickBody
}

func (s *sprite) ReplaceClickBody(in models.ClickOnBody) {
	s.clickBody = in
	s.clickBody.Pos(s.x, s.y)
	s.clickBody.Angle(s.angleRad)
}

func (s *sprite) PressedUserInput() *models.UserInput {
	return s.sim.PressedUserInput()
}

func (s *sprite) JustPressedUserInput() *models.UserInput {
	if s.userInputChan == nil {
		s.userInputChan = s.sim.SubscribeToJustPressedUserInput()
	}

	select {
	case i := <-s.userInputChan:
		return i
	default:
		// receiving from chan would block without this
	}

	return nil
}

func (s *sprite) WhoIsNearMe(distance float64) []models.NearMeInfo {
	return s.sim.WhoIsNearMe(s.x, s.y, distance)
}

func (s *sprite) SendMsg(toSpriteID int, msg any) {
	s.sim.SendMsg(toSpriteID, msg)
}

func (s *sprite) GetMsgs() []any {
	msgs := []any{}

GetAllReceivedMsgs:
	for {
		select {
		case i := <-s.receivedMsgs:
			msgs = append(msgs, i)
		default:
			// receiving from chan would block without this
			break GetAllReceivedMsgs
		}
	}
	return msgs
}

func (s *sprite) AddMsg(msg any) {
	s.receivedMsgs <- msg
}

func (s *sprite) minUpdate() {
	if s.deleted {
		log.Printf("Error: sprite %d is deleted but being updated\n", s.spriteID)
		return
	}

	s.sim.SpriteUpdatePosAngle(s)
}

func (s *sprite) fullUpdate() {
	if s.deleted {
		log.Printf("Error: sprite %d is deleted but being updated\n", s.spriteID)
		return
	}

	s.sim.SpriteUpdateFull(s)
}
