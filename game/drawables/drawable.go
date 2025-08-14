package drawables

import (
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"reflect"
)

type DrawableSaveAbleSprite interface {
	Draw(screen *ebiten.Image)
	Update()
	SavePosition() SavePositionData
}

type SavePositionData struct {
	X    float32
	Y    float32
	Name string
}

type Drawable interface {
	Draw(screen *ebiten.Image)
	Update()
	SpriteHovered() bool
	Coord() (float32, float32)
	ShouldRemove() bool
	Highlighted() bool
}

func ClosestDrawableOfTypeToCursor(cursorX int, cursorY int, drawSlice []Drawable, sortFor string) (Drawable, float64) {
	//takes sort for Parameter that is string of the base type of the drawable and will only sort for those
	//cursor x and y taking as inputs in case cursor position changes

	//filters for type
	toBeSorted := FilterDrawables(drawSlice, sortFor)

	//returns closest from filtered list
	drawable, distance := ClosestDrawableToCursor(cursorX, cursorY, toBeSorted)

	return drawable, distance

}

func FilterDrawables(drawSlice []Drawable, sortFor string) []Drawable {
	var toBeSorted []Drawable

	for _, draw := range drawSlice {
		if reflect.TypeOf(draw).String() == sortFor {
			toBeSorted = append(toBeSorted, draw)
		}
	}

	return toBeSorted
}

func ClosestDrawableToCursor(cursorX int, cursorY int, drawSlice []Drawable) (Drawable, float64) {
	//cursor x and y taking as inputs in case cursor position changes

	cursorPoint := image.Point{X: cursorX, Y: cursorY}

	var distMap = make(map[float64]Drawable)
	var closestDistance float64

	closestDistance = 1000

	for _, draw := range drawSlice {
		x, y := draw.Coord()
		pt := image.Point{X: int(x), Y: int(y)}
		dis := util.DistanceBetweenPoints(cursorPoint, pt)
		distMap[dis] = draw
		if dis < closestDistance {
			closestDistance = dis
		}
	}
	return distMap[closestDistance], closestDistance
}
