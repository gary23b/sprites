package models

import (
	"image"
)

type Scratch interface {
	GetWidth() int
	GetHeight() int

	AddCostume(img image.Image, name string)
	AddSprite(UniqueName string) Sprite
	DeleteSprite(Sprite)
	DeleteAllSprites()

	SpriteUpdatePosAngle(in Sprite)
	SpriteUpdateFull(in Sprite)

	AddSound(path, name string)
	PlaySound(name string, volume float64) // volume must be between 0 and 1.

	PressedUserInput() *UserInput
	SubscribeToJustPressedUserInput() chan *UserInput
	UnSubscribeToJustPressedUserInput(in chan *UserInput)

	GetSpriteInfo(UniqueName string) SpriteState

	Exit()
}
