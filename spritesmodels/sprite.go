package spritesmodels

type SpriteState struct {
	SpriteID       int
	SpriteType     int
	UniqueName     string
	CostumeName    string
	X, Y           float64
	Z              int // Effectively the layer index. 0 through 9 with 9 being the top.
	AngleDegrees   float64
	Visible        bool
	ScaleX, ScaleY float64
	Opacity        float64
	Deleted        bool
}

type NearMeInfo struct {
	SpriteID   int
	SpriteType int
	X, Y       float64
}
