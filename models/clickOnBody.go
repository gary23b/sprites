package models

type ClickOnBody interface {
	AddCircleBody(x, y, radius float64)
	AddRectangleBody(x1, x2, y1, y2 float64)

	IsMouseClickInBody(x, y float64) bool
	GetMousePosRelativeToOriginalSprite(x, y float64) (float64, float64)

	Clone() ClickOnBody

	// Should only be used by the sim.
	Pos(x, y float64)
	Angle(RadAngle float64)
}
