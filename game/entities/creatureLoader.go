package entities

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spriteSheet"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func LoadFishSprite(creatureType FishList, creatureLvl int) *sprite.AnimatedSprite {
	var c sprite.AnimatedSprite
	c.Sprite = &sprite.Sprite{}

	shaderParams := make(map[string]any)
	shaderParams["OutlineColor"] = [4]float64{255, 255, 0, 255}
	c.ShaderParams = shaderParams

	img, err := LoadFishImg(creatureType, creatureLvl)
	if err != nil {
		log.Fatal(err)
	}

	switch creatureLvl {
	case 1:
		c.Img = img
		c.Animation = animations.NewAnimation(0, 3, 1, 15)
		c.SpriteSheet = spritesheet.NewSpritesheet(4, 1, 21, 9)
	case 2:
		c.Img = img
		c.Animation = animations.NewAnimation(0, 3, 1, 15)
		c.SpriteSheet = spritesheet.NewSpritesheet(4, 1, 26, 13)
	case 3:
		c.Img = img
		c.Animation = animations.NewAnimation(0, 5, 1, 15)
		c.SpriteSheet = spritesheet.NewSpritesheet(6, 1, 40, 24)
	default:
		c.Img = img
		c.Animation = animations.NewAnimation(0, 3, 1, 15)
		c.SpriteSheet = spritesheet.NewSpritesheet(4, 1, 24, 11)
	}
	return &c
}

func LoadFishImg(fType FishList, level int) (*ebiten.Image, error) {
	var fishImgName string
	switch fType {
	case Fish:
		fishImgName = fmt.Sprintf("fish%dSpriteSheet", level)
	case MollyFish:
		fishImgName = fmt.Sprintf("mollyFish%dSpriteSheet", level)
	}
	img, err := loader.LoadImageAssetAsEbitenImage("fishSpriteSheets/" + fishImgName)
	if err != nil {
		return &ebiten.Image{}, err
	}
	return img, nil
}

func LoadFishSpriteAltAnimations(fType FishList) (*sprite.AnimatedSprite, error) {
	c := sprite.AnimatedSprite{}
	c.Sprite = &sprite.Sprite{}

	switch fType {
	case MollyFish:
		println("Loading Molly Fish Animation")
		img, err := loader.LoadImageAssetAsEbitenImage("fishSpriteSheets/mollyFishSpinAnimation")
		if err != nil {
			return &c, err
		}

		c.Img = img
	}

	c.Animation = animations.NewAnimation(0, 3, 1, 15)
	c.SpriteSheet = spritesheet.NewSpritesheet(7, 1, 21, 9)

	return &c, nil
}
