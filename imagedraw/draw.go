package imagedraw

import (
	"image"
	"image/color"
	"math"
)

type canvas struct {
	i       *image.RGBA
	w       int
	h       int
	halfW   float64
	halfH   float64
	penSize float64
	c       color.RGBA
	x, y    float64
	angle   float64 // radians
	penOn   bool
}

func NewImageDraw(width, height int) *canvas {
	i := image.NewRGBA(image.Rect(0, 0, width, height))

	ret := &canvas{
		i:       i,
		w:       width,
		h:       height,
		halfW:   float64(width) / 2.0,
		halfH:   float64(height) / 2.0,
		penSize: 0,
	}
	return ret
}

func (s *canvas) GetImage() *image.RGBA {
	return s.i
}

func (s *canvas) Color(c color.RGBA) {
	s.c = c
}

func (s *canvas) Angle(angleDegrees float64) {
	s.angle = angleDegrees * (math.Pi / 180.0)
}

func (s *canvas) Fill() {
	for y := 0; y < s.h; y++ {
		for x := 0; x < s.w; x++ {
			s.i.Set(x, y, s.c)
		}
	}
}

func (s *canvas) Copy() *canvas {
	// might be able to do the following
	ret := *s
	copy(ret.i.Pix, s.i.Pix)

	// ret := NewImageDraw(s.w, s.h)
	// for y := 0; y < s.h; y++ {
	// 	for x := 0; x < s.w; x++ {
	// 		ret.i.Set(x, y, s.i.At(x, y))
	// 	}
	// }
	return &ret
}

func (s *canvas) SetPixel(x, y int) {
	s.i.SetRGBA(x, y, s.c)
}

// 0,0 is the center of the screen. positive X is right, positive y is up.
func (s *canvas) CartesianPixelCord(x, y float64) (x2, y2 int) {
	x2, y2 = floatPosToPixel(x+s.halfW, -y+s.halfH)
	return x2, y2
}

// 0,0 is the center of the screen. positive X is right, positive y is up.
func (s *canvas) SetCartesianPixel(x, y float64) {
	pixX, pixY := s.CartesianPixelCord(x, y)
	s.SetPixel(pixX, pixY)
}

// This is what splits cartesian space into discrete pixels.
// This includes moving (0,0) be be centered in the middle of the (0,0) pixel. The center of the (0,0) pixel is at (.5, .5)
func floatPosToPixel(x, y float64) (int, int) {
	retX := int(math.Floor(x + .5))
	retY := int(math.Floor(y + .5))
	return retX, retY
}

// The concept of this line draw function is to determine if X or Y have a larger number of pixels to cover,
// and the larger one is chosen. Then we step
func (s *canvas) GoTo(x2, y2 float64) {
	if s.penOn {
		s.DrawLine(s.x, s.y, x2, y2)
	}

	s.x = x2
	s.y = y2
}

func (s *canvas) DrawLine(x1, y1, x2, y2 float64) {
	xDelta := x2 - x1
	yDelta := y2 - y1
	largerDelta := math.Max(math.Abs(xDelta), math.Abs(yDelta))

	loopSteps := int(math.Ceil(largerDelta))
	xStep := xDelta / float64(loopSteps)
	yStep := yDelta / float64(loopSteps)

	x := x1
	y := y1
	for i := 0; i <= loopSteps; i++ {
		s.SetCartesianPixel(x, y)

		if s.penSize > 0 {
			s.DrawFilledCircle(x, y, s.penSize)
		}
		x += xStep
		y += yStep
	}
}

func (s *canvas) DrawFilledCircle(x, y, size float64) {
	if !s.penOn {
		return
	}

	halfSize := size / 2
	halfSizeSquared := halfSize * halfSize
	xMax := int(math.Floor(x + halfSize))
	xMin := int(math.Floor(x - halfSize))
	yMax := int(math.Floor(y + halfSize))
	yMin := int(math.Floor(y - halfSize))

	for yInt := yMin; yInt <= yMax; yInt++ {
		yFlt := float64(yInt)
		for xInt := xMin; xInt <= xMax; xInt++ {
			xFlt := float64(xInt)
			deltaX := (float64(xInt) - x)
			deltaY := (float64(yInt) - y)
			distanceSquared := deltaX*deltaX + deltaY*deltaY
			if distanceSquared <= halfSizeSquared {
				s.SetCartesianPixel(xFlt, yFlt)
			}
		}
	}
}

func (s *canvas) CircleFromStartPos(radius, angleAmountToDraw float64, steps int) {
	// Convert to radians
	angleAmountToDraw *= (math.Pi / 180.0)

	angleAmountToDraw = max(angleAmountToDraw, -math.Pi*2.0)
	angleAmountToDraw = min(angleAmountToDraw, math.Pi*2.0)

	if radius < 0 {
		angleAmountToDraw *= -1
	}
	angleStepSize := angleAmountToDraw / float64(steps)
	endTurtleAngle := s.angle + angleAmountToDraw

	// Get center of Circle
	sin, cos := math.Sincos(s.angle + math.Pi/2.0)
	xCenter := s.x + radius*cos
	yCenter := s.y + radius*sin
	radius = math.Abs(radius)

	// Get the start of the circle
	deltaX := s.x - xCenter
	deltaY := s.y - yCenter
	startAngle := math.Atan2(deltaY, deltaX)

	for step := 1; step <= steps; step++ {
		currentAngle := startAngle + float64(step)*angleStepSize
		sin, cos := math.Sincos(currentAngle)
		x := xCenter + radius*cos
		y := yCenter + radius*sin
		s.angle += angleStepSize
		s.GoTo(x, y)
	}
	s.angle = endTurtleAngle
}

func (s *canvas) CircleAroundPoint(xCenter, yCenter, radius, startAngle, angleAmountToDraw float64, steps int) {
	if radius <= 0 {
		return
	}

	// Convert to radians
	startAngle *= (math.Pi / 180.0)
	angleAmountToDraw *= (math.Pi / 180.0)
	angleAmountToDraw = max(angleAmountToDraw, -math.Pi*2.0)
	angleAmountToDraw = min(angleAmountToDraw, math.Pi*2.0)

	angleStepSize := angleAmountToDraw / float64(steps)

	sin, cos := math.Sincos(startAngle)
	x1 := xCenter + radius*cos
	y1 := yCenter + radius*sin

	for step := 1; step <= steps; step++ {
		currentAngle := startAngle + float64(step)*angleStepSize
		sin, cos := math.Sincos(currentAngle)
		x2 := xCenter + radius*cos
		y2 := yCenter + radius*sin
		s.angle += angleStepSize
		s.DrawLine(x1, y1, x2, y2)
		x1 = x2
		y1 = y2
	}
}

func (s *canvas) BucketFill() {
	s.BucketFillPoint(s.x, s.y, s.c)
}

func colorMatches(a, b color.RGBA) bool {
	return a.R == b.R && a.G == b.G && a.B == b.B && a.A == b.A
}

func (s *canvas) BucketFillPoint(xIn, yIn float64, c color.RGBA) {
	type upNextStruct struct {
		x, y int
	}
	upNextStack := []upNextStruct{}

	xMin := s.i.Rect.Min.X
	yMin := s.i.Rect.Min.Y
	xMax := s.i.Rect.Max.X - 1
	yMax := s.i.Rect.Max.Y - 1

	xInt, yInt := s.CartesianPixelCord(xIn, yIn)

	srcColor := s.i.RGBAAt(xInt, yInt)
	if colorMatches(srcColor, c) {
		// The selected pixes is already the correct color
		return
	}
	upNextStack = append(upNextStack, upNextStruct{x: xInt, y: yInt})

	for len(upNextStack) > 0 {
		xy := upNextStack[len(upNextStack)-1]
		x := xy.x
		y := xy.y
		upNextStack = upNextStack[:len(upNextStack)-1]

		s.i.SetRGBA(x, y, c)
		if x > xMin && colorMatches(s.i.RGBAAt(x-1, y), srcColor) {
			upNextStack = append(upNextStack, upNextStruct{x: x - 1, y: y})
		}
		if x < xMax && colorMatches(s.i.RGBAAt(x+1, y), srcColor) {
			upNextStack = append(upNextStack, upNextStruct{x: x + 1, y: y})
		}
		if y > yMin && colorMatches(s.i.RGBAAt(x, y-1), srcColor) {
			upNextStack = append(upNextStack, upNextStruct{x: x, y: y - 1})
		}
		if y < yMax && colorMatches(s.i.RGBAAt(x, y+1), srcColor) {
			upNextStack = append(upNextStack, upNextStruct{x: x, y: y + 1})
		}
	}
}

func rotatePointAroundOrigin(x, y, angle float64) (x2, y2 float64) {
	sin, cos := math.Sincos(angle)
	x2 = cos*x - sin*y
	y2 = sin*x + cos*y

	return x2, y2
}

func (s *canvas) FakeEllipseAroundPoint(x, y, semiMajorAxis, semiMinorAxis, rotationAngle float64, steps int) {
	if semiMajorAxis <= 0 || semiMinorAxis <= 0 {
		return
	}

	// Convert to radians
	rotationAngle *= (math.Pi / 180.0)

	angleStepSize := math.Pi * 2.0 / float64(steps)

	x1 := semiMajorAxis
	y1 := 0.0
	x1, y1 = rotatePointAroundOrigin(x1, y1, rotationAngle)
	x1 += x
	y1 += y

	for step := 1; step <= steps; step++ {
		currentAngle := float64(step) * angleStepSize
		sin, cos := math.Sincos(currentAngle)
		x2 := semiMajorAxis * cos
		y2 := semiMinorAxis * sin
		x2, y2 = rotatePointAroundOrigin(x2, y2, rotationAngle)
		x2 += x
		y2 += y

		s.DrawLine(x1, y1, x2, y2)
		x1 = x2
		y1 = y2
	}
}
