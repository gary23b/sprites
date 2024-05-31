package models

type ClickOnBody interface {
	AddCirleBody(x, y, radius float64)
	AddRectangleBody(x1, x2, y1, y2 float64)

	IsMouseClickInBody(x, y float64) bool
	GetMousePosRelativeToOriginalSprite(x, y float64) (float64, float64)
}
