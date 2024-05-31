package models

import (
	"image"
)

type CmdAddNewSprite struct {
	SpriteID int
}

type CmdSpriteUpdateMin struct {
	SpriteID int
	X        float64
	Y        float64
	AngleRad float64
}

type CmdSpriteUpdateFull struct {
	SpriteID    int
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
	SpriteID int
}

type CmdSpritesDeleteAll struct {
	SpriteID int
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

type CmdGetScreenshot struct {
	ImageChan chan image.Image
}
