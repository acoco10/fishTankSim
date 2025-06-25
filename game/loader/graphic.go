package loader

import (
	"github.com/acoco10/fishTankWebGame/game/graphics"
	sprite "github.com/acoco10/fishTankWebGame/game/sprite"
)

func LoadClothGraphic() (*graphics.SpriteGraphic, error) {
	origin := [2]float32{721, 250}
	mx := [2]float32{800, 300}
	direction := "right"

	img, err := LoadImageAssetAsEbitenImage("menuAssets/cloth")
	if err != nil {
		return nil, err
	}
	gSprite := sprite.Sprite{Img: img, X: origin[0], Y: origin[1]}

	params := make(map[string]any)
	params["origin"] = origin
	params["max"] = mx
	params["direction"] = direction

	cloth := graphics.NewSpriteGraphic(gSprite, graphics.UpdateWhiteBoardCloth, params)
	return cloth, nil
}
