package graphics

import (
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/hajimehoshi/ebiten/v2"
)

func UpdateWhiteBoardCloth(gs *SpriteGraphic) {
	maxPoint := gs.parameters["max"].([2]float32)
	origin := gs.parameters["origin"].([2]float32)
	direction := gs.parameters["direction"].(string)

	if gs.Sprite.X >= maxPoint[0] && direction == "right" {
		gs.parameters["direction"] = "left"
	}

	if gs.Sprite.X <= origin[0] && direction == "left" {
		gs.parameters["direction"] = "right"
	}

	if gs.Sprite.Y >= maxPoint[1] {
		gs.parameters["direction"] = "stop"
		gs.complete = true
	}

	switch direction {
	case "right":
		gs.Sprite.X += 10
		gs.Sprite.Y++
	case "left":
		gs.Sprite.X -= 10
		gs.Sprite.Y++
	}

}

func UpdateFadeInGraphic(gs *SpriteGraphic) {
	opacity := gs.parameters["opacity"].(float32)
	pulse := false
	if gs.parameters["pulse"] != nil {
		pulse = true
	}
	if opacity < 1.0 {
		opacity += 0.02 // adjust fade-in speed here
	}
	if opacity > 1.0 && pulse {
		opacity = 0.0
	}
	if opacity > 1.0 && pulse {
		opacity = 1.0
	}
	gs.parameters["opacity"] = opacity
}

func DrawFadeInSprite(gs *SpriteGraphic, screen *ebiten.Image) {

	alpha := gs.parameters["opacity"].(float32)

	opts := &ebiten.DrawImageOptions{}
	opts.ColorScale.ScaleAlpha(alpha)
	opts.GeoM.Translate(float64(gs.Sprite.X), float64(gs.Sprite.Y))

	screen.DrawImage(gs.Sprite.Img, opts)
}

func NewFadeInSprite(sp sprite.Sprite) *SpriteGraphic {
	graphSpriteParams := map[string]any{}

	graphSpriteParams["opacity"] = float32(0.1)
	graphSpriteParams["pulse"] = true

	graph := NewSpriteGraphic(sp, UpdateFadeInGraphic, graphSpriteParams)
	graph.SetDrawFunc(FadeIn)

	return graph
}
