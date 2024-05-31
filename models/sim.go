package models

type Sim interface {
	AddSprite() Sprite
	DeleteAllSprites()
}
