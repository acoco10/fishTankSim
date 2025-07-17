package loader

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"log"
)

func getPath(s string) string {
	return fmt.Sprintf("data/animationData/%sAnimation.json", s)
}
func LoadEffect(eff string) *sprite.AnimatedSprite {

	path := getPath(eff)
	ani, ss, err := LoadAnimation(path)
	if err != nil {
		log.Fatal("cant load animation", err)
	}
	image, err := LoadImageAssetAsEbitenImage(fmt.Sprintf("effectSpriteSheets/%sSpriteSheet", eff))
	if err != nil {
		log.Fatal(err)
	}
	se := &sprite.Sprite{Img: image}
	ase := &sprite.AnimatedSprite{Sprite: se, SpriteSheet: ss, Animation: ani}

	return ase
}

func LoadStaticEffect(eff string, x float32, y float32) *sprite.Sprite {

	image, err := LoadImageAssetAsEbitenImage("staticEffects/" + eff)
	if err != nil {
		log.Fatal(err)
	}
	se := &sprite.Sprite{Img: image, X: x, Y: y}
	return se

}
