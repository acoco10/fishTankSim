package gameEntities

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spriteSheet"
	"github.com/hajimehoshi/ebiten/v2"
)

func LoadFishSprite(creatureType FishList, creatureLvl int) *AnimatedSprite {
	var c AnimatedSprite
	c.Sprite = &Sprite{}
	img := LoadFishImg(creatureType, creatureLvl)
	switch creatureLvl {
	case 1:
		c.Img = img
		c.Animation = animations.NewAnimation(0, 3, 1, 4)
		c.SpriteSheet = spritesheet.NewSpritesheet(4, 2, 19, 7)
	case 2:
		c.Img = img
		c.Animation = animations.NewAnimation(0, 3, 1, 4)
		c.SpriteSheet = spritesheet.NewSpritesheet(4, 1, 24, 11)
	case 3:
		c.Img = img
		c.Animation = animations.NewAnimation(0, 5, 1, 4)
		c.SpriteSheet = spritesheet.NewSpritesheet(6, 1, 40, 24)
	default:
		c.Img = img
		c.Animation = animations.NewAnimation(0, 3, 1, 4)
		c.SpriteSheet = spritesheet.NewSpritesheet(4, 1, 24, 11)
	}
	return &c
}

func LoadFishImg(fType FishList, level int) *ebiten.Image {
	var fishImgName string
	switch fType {
	case fish:
		fishImgName = fmt.Sprintf("fish%dSpriteSheet", level)
	case mollyFish:
		fishImgName = fmt.Sprintf("mollyFish%dSpriteSheet", level)
	}
	img := LoadImageAssetAsEbitenImage(fishImgName)
	return img
}
