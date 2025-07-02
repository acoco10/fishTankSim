package shaders

import (
	"github.com/acoco10/fishTankWebGame/game/geometry"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func UpdatePulseWithText(params map[string]any) map[string]any {
	opacity := params["Opacity"].(float32)

	cs := ebiten.ColorScale{}

	cs.SetR(0.1)
	cs.SetB(0.2)
	cs.SetG(1.0)
	cs.SetA(1.0)

	cs.ScaleAlpha(opacity)

	params["OutlineColor"] = [4]float32{cs.R(), cs.G(), cs.B(), cs.A()}

	if opacity < 1.0 {
		opacity += 0.02 // adjust fade-in speed here
	}

	if opacity >= 1.0 {
		opacity = 0.0
	}

	params["Opacity"] = opacity
	return params
}

func UpdateCounter(params map[string]any) map[string]any {
	if params["Counter"] == nil {
		//log.Printf("nil counter value inside shader update parameters")
		return nil
	}

	maxCounter := 1000

	if params["MaxCounter"] != nil {
		if params["MaxCounter"].(int) > 0 {
			maxCounter = params["MaxCounter"].(int)
		}
	}

	counter := params["Counter"].(int)

	counter++
	if counter > maxCounter {
		//arbitrary shut off at 1000 or at defined point passed through max counter param
		counter = 0
	}

	params["Counter"] = counter
	return params
}

func UpdateCounterOneShot(params map[string]any) map[string]any {

	if params["Counter"] == nil {
		//log.Printf("nil counter value inside shader update counter one shot parameters")
		return nil
	}

	maxCounter := 0

	if params["MaxCounter"] != nil {
		if params["MaxCounter"].(int) > 0 {
			maxCounter = params["MaxCounter"].(int)
		}
	} else {
		log.Printf("No max counter value recieved in one shot time update func: recheck code")
		return nil
	}

	counter := params["Counter"].(int)

	if counter >= maxCounter {
		params["Counter"] = counter
		return params
	}

	counter++
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
