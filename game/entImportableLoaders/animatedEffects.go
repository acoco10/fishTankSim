package entImportableLoaders

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/util"
	"log"
)

func getAnimationPath(s string) string {
	return fmt.Sprintf("data/animationData/%sAnimation.json", s)
}
func LoadEffect(eff string) *sprite.Sprite {

	path := getAnimationPath(eff)
	ani, ss, err := LoadAnimation(path)
	if err != nil {
		log.Fatal("cant load animation", err)
	}
	image, err := util.LoadImageAssetAsEbitenImage(fmt.Sprintf("effectSpriteSheets/%sSpriteSheet", eff))
	if err != nil {
		log.Fatal(err)
	}
	se := &sprite.Sprite{Img: image, SpriteSheet: ss, Animation: ani}

	return se
}

func LoadStaticEffect(eff string, x, y float32, location string) *sprite.Sprite {
	image, err := util.LoadImageAssetAsEbitenImage("staticEffects/" + eff)
	if err != nil {
		log.Fatal(err)
	}
	if location == "LM" {
		xM := float32(registry.Config.ScreenWidth/2 + int(float64(image.Bounds().Dx()))/2)
		yL := float32(registry.Config.ScreenHeight-image.Bounds().Dy()) * 1.5
		se := &sprite.Sprite{Img: image, X: xM, Y: yL}
		return se
	} else {
		se := &sprite.Sprite{Img: image, X: x, Y: y}
		return se
	}
}
