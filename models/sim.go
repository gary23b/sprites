package models

import "image"

type Sim interface {
	AddCostume(img image.Image, name string)
	AddSprite() Sprite
	DeleteAllSprites()
}
