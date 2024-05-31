package models

import (
	"image"
)

type NearMeInfo struct {
	SpriteID   int
	SpriteType int
	X, Y       float64
}

type Scratch interface {
	GetWidth() int
	GetHeight() int

	AddCostume(img image.Image, name string)
	AddSprite(UniqueName string) Sprite // If no name is given, a random name is generated.
	DeleteSprite(Sprite)
	DeleteAllSprites()

	SpriteUpdatePosAngle(in Sprite)
	SpriteUpdateFull(in Sprite)

	AddSound(path, name string)
	PlaySound(name string, volume float64) // volume must be between 0 and 1.

	PressedUserInput() *UserInput
	SubscribeToJustPressedUserInput() chan *UserInput
	UnSubscribeToJustPressedUserInput(in chan *UserInput)

	GetSpriteID(UniqueName string) int
	GetSpriteInfo(UniqueName string) SpriteState
	GetSpriteInfoByID(id int) SpriteState

	WhoIsNearMe(x, y, distance float64) []NearMeInfo
	SendMsg(toSpriteID int, msg any)

	GetScreenshot() image.Image

	Exit()
}
