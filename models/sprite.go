package models

import "image"

///////////////////////

type SpriteState struct {
	CostumeName    string
	X, Y           float64
	Z              int // Effetively the layer index. 0 through 9 with 9 being the top.
	Angle          float64
	Visible        bool
	ScaleX, ScaleY float64
	Opacity        float64
}

type Sprite interface {
	// Init
	AddCostume(img image.Image, name string)

	// Updates
	Costume(name string)
	Angle(angleDegrees float64)
	Pos(cartX, cartY float64) // Cartesian (x,y). Center in the middle of the window
	Z(int)                    //
	Visible(visible bool)
	Scale(scale float64) // Sets xScale and yScale together
	XYScale(xScale, yScale float64)
	Opacity(opacityPercent float64) // 0 is completely transparent and 100 is completely opaque
	All(in SpriteState)

	// Info
	GetState() SpriteState

	// exit
	DeleteSprite()
}
