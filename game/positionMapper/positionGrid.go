package positionMapper

import (
	"github.com/acoco10/fishTankWebGame/game/geometry"
)

type SpatialGrid struct {
	positions map[int]*geometry.Rect
}
