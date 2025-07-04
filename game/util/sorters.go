package util

import (
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math"
	"reflect"
)

func ClosestPoint(point1 image.Point, candidates []image.Point) image.Point {
	distMap := make(map[float64]image.Point)
	closestDistance := Distance(candidates[0], point1)

	for _, pt := range candidates {

		dis := Distance(point1, pt)

		distMap[dis] = pt

		if dis < closestDistance {
			closestDistance = dis
		}
	}
	return distMap[closestDistance]
}

func Distance(p1, p2 image.Point) float64 {
	dx := float64(p2.X - p1.X)
	dy := float64(p2.Y - p1.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

func SpriteToPoint(sp sprite.Sprite) image.Point {
	pt := image.Point{int(sp.X), int(sp.Y)}
	return pt
}

func ClosestSpriteToCursor(sps []*sprite.Sprite) *sprite.Sprite {
	cursorX, cursorY := ebiten.CursorPosition()
	cursorPoint := image.Point{X: cursorX, Y: cursorY}
	var distMap = make(map[float64]*sprite.Sprite)
	var closestDistance float64

	for _, sp := range sps {
		pt := SpriteToPoint(*sp)
		dis := Distance(cursorPoint, pt)
		distMap[dis] = sp
		if dis < closestDistance {
			closestDistance = dis
		}
	}

	return distMap[closestDistance]
}

func ClosestDrawableToCursor(cursorX int, cursorY int, drawSlice []drawables.Drawable, filter func(any) bool, sortFor string) drawables.Drawable {
	//cursor x and y taking as inputs in case cursor position changes?
	var toBeSorted []drawables.Drawable

	for _, draw := range drawSlice {
		println(reflect.TypeOf(draw).String())

		if reflect.TypeOf(draw).String() == sortFor {
			toBeSorted = append(toBeSorted, draw)
		}
	}

	cursorPoint := image.Point{X: cursorX, Y: cursorY}

	var distMap = make(map[float64]drawables.Drawable)
	var closestDistance float64
	closestDistance = 1000

	for _, draw := range toBeSorted {
		x, y := draw.Coord()
		pt := image.Point{X: int(x), Y: int(y)}
		dis := Distance(cursorPoint, pt)
		distMap[dis] = draw
		if dis < closestDistance {
			closestDistance = dis
		}
	}

	if filter(closestDistance) {
		println("Returning closest creature in closest creature func")
		return distMap[closestDistance]
	}

	println("no creature close enough to return in closest creature func")

	return nil

}
