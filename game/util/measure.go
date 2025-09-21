package util

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image"
)

func MeasureText(outputText string, fontsize float64, font string) (float64, float64) {
	face := registry.FontMap[font]

	return text.Measure(outputText, face, 2)
}

func Lerp64(A float64, B float64, t float64) float64 {
	return A + (B-A)*t
}

func Lerp32(A float32, B float32, t float32) float32 {
	return A + (B-A)*t
}

func GetScaledCursorPosition() (int, int) {
	x, y := ebiten.CursorPosition()

	if registry.Config.CursorSpeed != 0 {
		x = registry.Config.CursorPosition.X
		y = registry.Config.CursorPosition.Y
	}
	if registry.Config.Zoom {

		scaledX := float64(x) / registry.Config.ResolutionScalingF / registry.Config.ZoomFactor
		xOffSet := registry.Config.ZoomOffSetX / registry.Config.ZoomFactor
		xf := scaledX - xOffSet

		scaledY := float64(y-registry.Config.YOffset) / registry.Config.ResolutionScalingF / registry.Config.ZoomFactor
		yOffSet := registry.Config.ZoomOffSetY / registry.Config.ZoomFactor
		yf := scaledY - yOffSet

		return int(xf), int(yf)

	}
	return int(float64(x) / registry.Config.ResolutionScalingF), int(float64(y-registry.Config.YOffset) / registry.Config.ResolutionScalingF)

}

func ChopUpIcons(inputImage *ebiten.Image, labels []string, size int) (map[string]*ebiten.Image, map[string]*ebiten.Image) {
	imageMap := make(map[string]*ebiten.Image)
	indMap := make(map[string]*ebiten.Image)

	for i, icon := range labels {
		//horizontal slice of square images
		rect := image.Rect(i*size, 0, (i+1)*size, size)

		if i < 3 {
			indMap[icon] = ebiten.NewImageFromImage(inputImage.SubImage(rect))
		} else {
			imageMap[icon] = ebiten.NewImageFromImage(inputImage.SubImage(rect))
		}
	}

	return imageMap, indMap
}
