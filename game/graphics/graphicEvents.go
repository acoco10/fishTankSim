package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/geometry"
)

type DrawGraphic struct {
	Point   *geometry.Point
	Graphic drawables.DrawableSprite
}
