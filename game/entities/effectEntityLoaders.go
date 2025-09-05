package entities

import (
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/util"
	"image"
	"log"
)

var Effects = make(map[string]*Entity)

func LoadEffects(zBounds [13]image.Rectangle) {
	spotLight, err := util.LoadImageAssetAsEbitenImage("staticEffects/spotLight")
	if err != nil {
		log.Fatal("cant load spotlight effect", err)
	}

	spotLightSprite := &sprite.Sprite{Img: spotLight, Y: float32(zBounds[0].Min.Y - 5)}

	ent := &Entity{Sprite: spotLightSprite}
	ent.Z = 10
	RegisterEntity(ent)
	Effects["spotLight"] = ent
	ent.Draw = false
}

func DrawSpotLight(x float32) {
	midPoint := Effects["spotLight"].Sprite.GetSpriteRect().Dx() / 2
	Effects["spotLight"].Draw = true
	Effects["spotLight"].Sprite.X = x - float32(midPoint)
}

func TurnOffSpotLight() {
	Effects["spotLight"].Draw = false
}
