package models

import (
	"image"
)

type Sim interface {
	GetWidth() int
	GetHeight() int

	AddCostume(img image.Image, name string)
	AddSprite() Sprite
	DeleteSprite(Sprite)
	DeleteAllSprites()

	SpriteMinUpdate(in *CmdSpriteUpdateMin)
	SpriteFullUpdate(in *CmdSpriteUpdateFull)

	AddSound(path, name string)
	PlaySound(name string, volume float64) // volume must be between 0 and 1.

	PressedUserInput() *UserInput
	SubscribeToJustPressedUserInput() chan *UserInput
	UnSubscribeToJustPressedUserInput(in chan *UserInput)

	Exit()
}
