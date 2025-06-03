package loaders

import (
	"encoding/json"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	interactableUIObjects2 "github.com/acoco10/fishTankWebGame/game/interactableUIObjects"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func loadUiSpritesImgs(label interactableUIObjects2.UISpriteLabel) ([]*ebiten.Image, error) {
	var imgs []*ebiten.Image
	tags := []string{"Main", "Outline", "Alt"}

	for _, tag := range tags {
		assetName := string(label) + tag
		img, err := LoadImageAssetAsEbitenImage("uiSprites/" + assetName)
		if err != nil {
			log.Printf("%s not found for loading UiSprite %s, proceeding with loading other files. error msg: %s", assetName, string(label), err)
		} else {
			imgs = append(imgs, img)
		}
	}
	return imgs, nil
}

func LoadUISprites(spritesToLoad []interactableUIObjects2.UISpriteLabel, hub *events.EventHub, screenWidth, screenHeight int) ([]drawables.DrawableSprite, error) {
	var sprites []drawables.DrawableSprite

	spritePositions, err := loadSpritePositionData()
	if err != nil {
		return nil, err
	}

	for _, elem := range spritesToLoad {
		x := spritePositions[string(elem)].X
		y := spritePositions[string(elem)].Y
		imgs, err := loadUiSpritesImgs(elem)
		if err != nil {
			return sprites, err
		}
		sprite := interactableUIObjects2.NewUiSprite(imgs, hub, x, y, string(elem), screenWidth, screenHeight)
		if elem == interactableUIObjects2.FishFood {
			ffSprite := interactableUIObjects2.FishFoodSprite{sprite}
			ffSprite.Subscribe()
			lightingShader := shaders.LoadOnePointLightingNeutral()
			ffSprite.Shader = lightingShader
			LoadSpriteLightingParams(sprite.Sprite)
			sprites = append(sprites, &ffSprite)
			continue
		}
		if elem == interactableUIObjects2.WhiteBoard {
			wbSprite := interactableUIObjects2.WhiteBoardSprite{UiSprite: sprite}
			wbSprite.Init()
			wbSprite.Subscribe(hub)
			lightingShader := shaders.LoadOnePointLightingNeutral()
			wbSprite.Sprite.Shader = lightingShader
			LoadSpriteLightingParams(wbSprite.Sprite)
			sprites = append(sprites, &wbSprite)
			continue
		}
		//lightingShader := shaders.LoadOnePointLightingNeutral()
		//sprite.Shader = lightingShader
		//LoadSpriteLightingParams(sprite.Sprite)
		sprites = append(sprites, sprite)

	}

	return sprites, nil

}

func loadSpritePositionData() (map[string]*drawables.SavePositionData, error) {
	var positions = make(map[string]*drawables.SavePositionData)
	spritePosition, err := assets.DataDir.ReadFile("data/spritePosition.json")
	if err != nil {
		return positions, err
	}
	json.Unmarshal(spritePosition, &positions)
	return positions, nil
}

func LoadSelectedAnimations() (map[string]*sprite.AnimatedSprite, error) {
	fishOptions := make(map[string]*sprite.AnimatedSprite)

	mollyFishSprite, err := LoadFishSpriteAltAnimations(entities.MollyFish)
	if err != nil {
		log.Fatal(err)
	}
	mollyFishSprite.Animation.SpeedInTPS = 4

	mollyFishSprite.X = 420
	mollyFishSprite.Y = 260

	fishOptions["mollyFish"] = mollyFishSprite

	return fishOptions, nil
}

//imgLoader() []*ebitenImgs
//jsonLoader
