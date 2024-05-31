package models

import (
	"image"
)

type CmdAddNewSprite struct {
	SpriteID int
}

type CmdSpriteUpdateMin struct {
	SpriteIndex int
	CostumeName string
	X           float64
	Y           float64
}

type CmdSpriteUpdateFull struct {
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

type CmdSpriteDelete struct {
	SpriteIndex int
}

type CmdSpritesDeleteAll struct {
	SpriteIndex int
}

type CmdAddCostume struct {
	CostumeName string
	Img         image.Image
}

type CmdAddSound struct {
	Path      string
	SoundName string
}

type CmdPlaySound struct {
	SoundName string
	Volume    float64 // between 0 and 1.
}
