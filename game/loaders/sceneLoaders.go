package loaders

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/sprite"
)

func LoadAnimatedSelectSprites() (map[string]*sprite.AnimatedSprite, error) {

	fishOptions := make(map[string]*sprite.AnimatedSprite)

	fishSprite := LoadFishSprite(entities.Fish, 1)

	//scaledImg := ebiten.NewImage(fishSprite.Img.Bounds().Dx()*4, fishSprite.Img.Bounds().Dy()*4)

	//dopts := &ebiten.DrawImageOptions{}
	//dopts.GeoM.Scale(4, 4)
	//scaledImg.DrawImage(fishSprite.Img, dopts)
	//fishSprite.Img = scaledImg

	//fishSprite.SpriteWidth = fishSprite.SpriteWidth * 4
	//fishSprite.SpriteHeight = fishSprite.SpriteHeight * 4

	fishSprite.X = 355
	fishSprite.Y = 260

	mollyFishSprite := LoadFishSprite(entities.MollyFish, 1)
	mollyFishSprite.Animation.SpeedInTPS = 10

	mollyFishSprite.X = 490
	mollyFishSprite.Y = 260

	fishOptions["Common Molly"] = mollyFishSprite
	fishOptions["Goldfish"] = fishSprite

	return fishOptions, nil
}
