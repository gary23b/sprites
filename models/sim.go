package models

import "image"

type Sim interface {
	AddCostume(img image.Image, name string)
	AddSprite() Sprite
	DeleteAllSprites()

	AddSound(path, name string)
	PlaySound(name string, volume float64) // volume must be between 0 and 1.
}
