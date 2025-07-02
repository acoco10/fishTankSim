package entity

import "github.com/acoco10/fishTankWebGame/game/sprite"

type Entity struct {
	Id       int
	Movement bool
	Sprite   sprite.Sprite
}
