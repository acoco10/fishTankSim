package util

import (
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"math"
)

var DebugZColors = []color.Color{
	colornames.Blue,            // 1
	colornames.Royalblue,       // 2
	colornames.Dodgerblue,      // 3
	colornames.Deepskyblue,     // 4
	colornames.Skyblue,         // 5
	colornames.Lightblue,       // 6
	colornames.Cyan,            // 7
	colornames.Turquoise,       // 8
	colornames.Mediumturquoise, // 9
	colornames.Lightseagreen,   // 10
	colornames.Mediumseagreen,  // 11
	colornames.Seagreen,        // 12
	colornames.Green,           // 13
}

func ClosestPoint(point1 image.Point, candidates []image.Point) image.Point {
	distMap := make(map[float64]image.Point)
	closestDistance := DistanceBetweenPoints(candidates[0], point1)

	for _, pt := range candidates {

		dis := DistanceBetweenPoints(point1, pt)

		distMap[dis] = pt

		if dis < closestDistance {
			closestDistance = dis
		}
	}
	return distMap[closestDistance]
}

func DistanceBetweenPoints(p1, p2 image.Point) float64 {
	dx := float64(p2.X - p1.X)
	dy := float64(p2.Y - p1.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

func ClosestPointToCursor(sps []image.Point) image.Point {
	cursorX, cursorY := GetScaledCursorPosition()
	cursorPoint := image.Point{X: int(cursorX), Y: int(cursorY)}

	var distMap = make(map[float64]image.Point)
	var closestDistance float64

	for _, sp := range sps {
		dis := DistanceBetweenPoints(cursorPoint, sp)
		distMap[dis] = sp
		if dis < closestDistance {
			closestDistance = dis
		}
	}

	return distMap[closestDistance]
}
