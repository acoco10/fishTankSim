package gameEntities

import (
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spriteSheet"
)

func LoadFishSprite(creatureType string, creatureLvl int) *AnimatedSprite {
	var c AnimatedSprite
	c.Sprite = &Sprite{}
	switch creatureLvl {
	case 1:
		img := LoadImageAssetAsEbitenImage("fishSpriteSheet")
		c.Img = img
		c.Animation = animations.NewAnimation(0, 3, 1, 4)
		c.SpriteSheet = spritesheet.NewSpritesheet(4, 2, 19, 7)
	case 2:
		img := LoadImageAssetAsEbitenImage("fish2SpriteSheet")
		c.Img = img
		c.Animation = animations.NewAnimation(0, 3, 1, 4)
		c.SpriteSheet = spritesheet.NewSpritesheet(4, 1, 24, 11)
	default:
		img := LoadImageAssetAsEbitenImage("fish2SpriteSheet")
		c.Img = img
		c.Animation = animations.NewAnimation(0, 3, 1, 4)
		c.SpriteSheet = spritesheet.NewSpritesheet(4, 1, 24, 11)
	}
	return &c
}
