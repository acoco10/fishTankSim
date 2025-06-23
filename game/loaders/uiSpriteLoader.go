package loaders

import (
	"encoding/json"
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spriteSheet"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/acoco10/fishTankWebGame/game/drawables"
	iObj "github.com/acoco10/fishTankWebGame/game/interactableUIObjects"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func loadUiSpritesImgs(label iObj.Label) ([]*ebiten.Image, error) {
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

func LoadUISprites(spritesToLoad []iObj.Label, hub *tasks.EventHub, screenWidth, screenHeight int) ([]drawables.DrawableSaveAbleSprite, error) {
	var sprites []drawables.DrawableSaveAbleSprite

	spritePositions, err := loadSpritePositionData()
	if err != nil {
		return nil, err
	}

	for _, elem := range spritesToLoad {
		x := spritePositions[string(elem)].X
		y := spritePositions[string(elem)].Y

		imgs, err2 := loadUiSpritesImgs(elem)
		if err2 != nil {
			return sprites, err
		}

		uSprite := iObj.NewUiSprite(imgs, hub, x, y, string(elem), screenWidth, screenHeight)

		iObj.Pubs(hub, *uSprite)
		switch elem {
		case iObj.FishFood:
			ffSprite := iObj.FishFoodSprite{UiSprite: uSprite}
			ffSprite.Subscribe()
			sprites = append(sprites, &ffSprite)
			continue
		case iObj.WhiteBoard:
			wbSprite := iObj.WhiteBoardSprite{UiSprite: uSprite}
			wbSprite.Init()
			wbSprite.Subscribe(hub)
			sprites = append(sprites, &wbSprite)
			continue
		case iObj.PiggyBank:
			pbSprite := iObj.PiggyBankSprite{UiSprite: uSprite}
			//pbSprite.Init()
			pbSprite.Subscribe(hub)
			sprites = append(sprites, &pbSprite)
			aniMap := LoadPiggyBankAnimationMap(x, y, float32(pbSprite.Img.Bounds().Dy()))
			pbSprite.AnimationMap = aniMap
			continue
		case iObj.Pillow:
			pillowSprite := iObj.PillowUI{UiSprite: uSprite, Triggered: false}
			pillowSprite.Subscribe(hub)
			sprites = append(sprites, &pillowSprite)
			continue
		}

		//lightingShader := shaders.LoadOnePointLightingNeutral()
		//sprite.Shader = lightingShader
		//LoadSpriteLightingParams(sprite.Sprite)
		sprites = append(sprites, uSprite)

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

func LoadPiggyBankAnimationMap(x, y, srcImageHeight float32) map[string]drawables.Drawable {

	aniMap := make(map[string]drawables.Drawable)
	img, err := LoadImageAssetAsEbitenImage("uiSprites/allowanceCollectedAni")
	if err != nil {
		log.Fatal(err, "cant load piggy bank animation thing")
	}
	animation := animations.NewAnimation(0, 7, 1, 5)
	spriteSheet := spritesheet.NewSpritesheet(8, 1, 149, 202)

	animatedSprite := sprite.NewAnimatedSprite()
	animatedSprite.Img = img
	animatedSprite.Animation = animation
	animatedSprite.SpriteSheet = spriteSheet

	animatedSprite.X = x
	yOffSet := float32(animatedSprite.Img.Bounds().Dy()) - srcImageHeight
	animatedSprite.Y = y - yOffSet

	aniMap["allowance"] = animatedSprite
	return aniMap
}

//imgLoader() []*ebitenImgs
//jsonLoader
