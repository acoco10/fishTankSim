package shaders

import (
	"github.com/acoco10/fishTankWebGame/game/geometry"
)

func UpdatePulse(params map[string]any) map[string]any {
	counter := params["Counter"].(int)
	counter++
	if counter > 40 {
		counter = 0
	}
	params["Counter"] = counter
	return params
}

func CpuUpdateRotatingOutlineHighlight(params map[string]any, cpuParams map[string]any) (map[string]any, map[string]any) {
	//variables used in game loop to be mutated by update func
	cpuDirection := cpuParams["direction"].(geometry.Direction)
	imageRect := cpuParams["imageRect"].([2]float64)
	origin := cpuParams["origin"].([2]float64)
	highLightRect := cpuParams["hlRectPoint"].([2]float64)

	//variables passed to shader to be mutated by update func
	rectangleSize := params["HLRectSize"].(float64)

	updateSpeed := 0.7

	X0 := origin[0]
	Y0 := origin[1]

	switch cpuDirection {
	case geometry.Right:
		highLightRect[0] += updateSpeed
		if highLightRect[0]+rectangleSize >= imageRect[0] {
			cpuParams["direction"] = geometry.Down
		}
	case geometry.Left:
		highLightRect[0] -= updateSpeed
		if highLightRect[0] <= 0 {
			cpuParams["direction"] = geometry.Up
		}
	case geometry.Down:
		highLightRect[1] += updateSpeed
		if highLightRect[1]+rectangleSize >= imageRect[1] {
			cpuParams["direction"] = geometry.Left
		}
	case geometry.Up:
		highLightRect[1] -= updateSpeed
		if highLightRect[1] <= 0 {
			cpuParams["direction"] = geometry.Right
		}
	}

	params["HighLightRect"] = [2]float64{highLightRect[0] + X0, highLightRect[1] + Y0}
	cpuParams["hlRectPoint"] = highLightRect
	return params, cpuParams
}
