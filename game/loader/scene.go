package loader

import (
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/util"
)

func StartScreenCoordinatePositioner(screenHeight int, screenWidth int, spriteMap map[string]drawables.Drawable, fontsize float64, headerFontSize float64) {
	//theoretically programmable button parameters

	textBuffer := 103
	buttonSpacing := 20
	minSpriteSelectButtonWidth := 200
	minSpriteSelectButtonHeight := 200
	rowSpacing := 20

	_, headerHeight := util.MeasureText("Select Your Fish", 54, "nk57")

	yLocation := screenHeight/5 + rowSpacing + int(headerHeight)

	i := 0
	orderKeys := []string{"Goldfish", "Common Molly", "Back"}

	for _, key := range orderKeys {
		fish, ok := spriteMap[key].(*sprite.AnimatedSprite)
		if ok {
			width, _ := util.MeasureText(key, fontsize, "nk57")
			widthAndBuffer := int(width) + 2*textBuffer

			if widthAndBuffer < minSpriteSelectButtonWidth {
				widthAndBuffer = minSpriteSelectButtonWidth
			}

			imgWidth := fish.Img.Bounds().Dx()

			offSet := widthAndBuffer - imgWidth

			if i == 0 {
				fish.X = float32(screenWidth/2 - imgWidth - buttonSpacing/2 - offSet/2)
			} else {
				fish.X = float32(screenWidth/2 + buttonSpacing/2 + offSet/2)
			}

			fish.Y = float32(yLocation + minSpriteSelectButtonHeight/2)
			i++
		}

		sp, ok := spriteMap[key].(*sprite.Sprite)
		if ok {
			sp.Y = float32(yLocation + minSpriteSelectButtonHeight/2)
			sp.X = float32(screenWidth/2 - minSpriteSelectButtonWidth)
		}
	}
}
