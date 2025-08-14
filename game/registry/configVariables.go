package registry

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"math"
)

var Config ConfigManager

type ConfigValue int

const (
	ResolutionScaling ConfigValue = iota
	ResolutionHeight
	ResolutionWidth
	ScreenWidth
	ScreenHeight
	ScreenOffsetY
	Zoom
	ZoomFactor
	ZoomOffset
)

type GameMode uint8

const (
	Position GameMode = iota
	Normal
)

type ConfigManager struct {
	ResolutionScalingi int
	ResolutionScalingf float32
	ResolutionScalingF float64
	ResolutionHeight   int
	ResolutionWidth    int
	ScreenWidth        int
	ScreenHeight       int
	YOffset            int
	YOffsetF           float64
	ScaledYOffsetF     float64
	Zoom               bool
	ZoomFactor         float64
	ZoomOffSetY        float64
	ZoomOffSetX        float64
}

func (c *ConfigManager) Set(value ConfigValue, input any) {
	switch value {
	case ResolutionScaling:
		i, ok := input.(float64)
		if !ok {
			log.Fatal("Tried to set non integer resolution scalar")
		}
		Config.ResolutionScalingi = int(i)
		Config.ResolutionScalingf = float32(i)
		Config.ResolutionScalingF = i

		targetAspectRatio := 16.0 / 9.0
		x, y := ebiten.WindowSize()
		xf, yf := float64(x), float64(y)
		windowAspectRatio := xf / yf

		tolerance := 0.01
		if math.Abs(windowAspectRatio-targetAspectRatio) > tolerance {
			scaledHeight := c.ScreenHeight * c.ResolutionScalingi
			dif := c.ResolutionHeight - scaledHeight
			c.YOffset = dif / 8
			c.YOffsetF = float64(dif / 8)
			c.ScaledYOffsetF = c.ResolutionScalingF * c.YOffsetF
		}

	case ResolutionWidth:
		i, ok := input.(int)
		if !ok {
			log.Fatal("Tried to set non integer ResolutionWidth")
		}
		Config.ResolutionWidth = i
	case ResolutionHeight:
		i, ok := input.(int)
		if !ok {
			log.Fatal("Tried to set non integer ResolutionHeight")
		}
		Config.ResolutionHeight = i
	case ScreenWidth:
		i, ok := input.(int)
		if !ok {
			log.Fatal("Tried to set non integer screen width")
		}
		Config.ScreenWidth = i
	case ScreenHeight:
		i, ok := input.(int)
		if !ok {
			log.Fatal("Tried to set non integer screen height")
		}
		Config.ScreenHeight = i
	case Zoom:
		i, ok := input.(bool)
		if !ok {
			log.Fatal("Tried to set zoom to non bool value")
		}
		Config.Zoom = i
	case ZoomFactor:
		i, ok := input.(float64)
		if !ok {
			log.Fatal("Tried to set zoom to non bool value")
		}
		Config.ZoomFactor = i
	case ZoomOffset:
		i, ok := input.(image.Point)
		if !ok {
			log.Fatal("Tried to set zoom to non bool value")
		}
		Config.ZoomOffSetX = float64(i.X)
		Config.ZoomOffSetY = float64(i.Y)
	}

}
