package spritestools

import (
	"log"
	"math"

	"github.com/gary23b/sprites/models"
)

type circle struct {
	x, y   float64
	radius float64
}

type rectangle struct {
	x1, x2 float64
	y1, y2 float64
}

type ClickOnBody struct {
	radiusOfCaring float64
	circles        []circle
	rectangles     []rectangle

	x, y     float64
	radAngle float64
}

var _ models.ClickOnBody = &ClickOnBody{}

func NewTouchCollisionBody() *ClickOnBody {
	ret := &ClickOnBody{}

	return ret
}

func (s *ClickOnBody) AddCircleBody(x, y, radius float64) {
	newC := circle{
		x:      x,
		y:      y,
		radius: radius,
	}
	s.circles = append(s.circles, newC)

	maxDis := math.Sqrt(x*x+y*y) + radius
	if maxDis > s.radiusOfCaring {
		s.radiusOfCaring = maxDis
	}
}

func distFromOrigin(x, y float64) float64 {
	return math.Sqrt(x*x + y*y)
}

func (s *ClickOnBody) AddRectangleBody(x1, x2, y1, y2 float64) {
	if x1 >= x2 || y1 >= y2 {
		log.Println("x2 must be greater than x1")
		return
	}

	newR := rectangle{
		x1: x1,
		x2: x2,
		y1: y1,
		y2: y2,
	}
	s.rectangles = append(s.rectangles, newR)

	s.radiusOfCaring = math.Max(s.radiusOfCaring, distFromOrigin(x1, y1))
	s.radiusOfCaring = math.Max(s.radiusOfCaring, distFromOrigin(x1, y2))
	s.radiusOfCaring = math.Max(s.radiusOfCaring, distFromOrigin(x2, y1))
	s.radiusOfCaring = math.Max(s.radiusOfCaring, distFromOrigin(x2, y2))
}

func (s *ClickOnBody) Pos(x, y float64) {
	s.x = x
	s.y = y
}

func (s *ClickOnBody) Angle(radAngle float64) {
	s.radAngle = radAngle
}

func (s *ClickOnBody) IsMouseClickInBody(x, y float64) bool {
	// get the mouse position relative to the origin centered body
	x -= s.x
	y -= s.y

	// Check if I care
	distanceSquared := x*x + y*y
	if distanceSquared >= s.radiusOfCaring*s.radiusOfCaring {
		return false
	}

	// Rotate the mouse around the origin based on the body angle
	if s.radAngle != 0 {
		sin, cos := math.Sincos(-s.radAngle)
		x0 := x
		x = cos*x0 - sin*y
		y = sin*x0 + cos*y
	}

	// Loop through the circles
	for i := range s.circles {
		c := s.circles[i]
		dx := c.x - x
		dy := c.y - y
		distanceSquared = dx*dx + dy*dy
		if distanceSquared < c.radius*c.radius {
			return true
		}
	}

	// Loop through the rectangles
	for i := range s.rectangles {
		r := s.rectangles[i]
		if x > r.x1 && x < r.x2 && y > r.y1 && y < r.y2 {
			return true
		}
	}

	return false
}

func (s *ClickOnBody) GetMousePosRelativeToOriginalSprite(x, y float64) (float64, float64) {
	// get the mouse position relative to the origin centered body
	x -= s.x
	y -= s.y

	// Rotate the mouse around the origin based on the body angle
	if s.radAngle != 0 {
		sin, cos := math.Sincos(-s.radAngle)
		x0 := x
		x = cos*x0 - sin*y
		y = sin*x0 + cos*y
	}

	return x, y
}

/*
func (s *ClickOnBody) AreWeTouchingAnotherBody(other *ClickOnBody) bool {

}

func overlapCircles(c1, c2 *circle) bool {
	dx := c1.x - c2.x
	dy := c1.y - c2.y

	distance := math.Sqrt(dx*dx + dy*dy)
	return distance < c1.radius+c2.radius
}

func overlapCircleRectangle(c *circle, r *rectangle) bool {
	// should check radius of caring first
	//


	//git the relative position of the circle compared to the rectangle
	x := c.x - r.x

	//rotate the circle center by the rotation of the rectangle


	if s.radAngle != 0 {
		sin, cos := math.Sincos(-s.radAngle)
		x0 := x
		x = cos*x0 - sin*y
		y = sin*x0 + cos*y
	}
	dx := c1.x - c2.x
	dy := c1.y - c2.y

	distance := math.Sqrt(dx*dx + dy*dy)
	return distance < c1.radius+c2.radius
}
*/

func (s *ClickOnBody) Clone() models.ClickOnBody {
	ret := *s
	copy(ret.circles, s.circles)
	copy(ret.rectangles, s.rectangles)
	return &ret
}
