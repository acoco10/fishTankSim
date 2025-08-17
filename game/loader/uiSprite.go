package loader

import (
	"encoding/json"
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spriteSheet"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/system"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

func loadUiSpritesImgs(label entities.Label) ([]*ebiten.Image, error) {
	var imgs []*ebiten.Image
	tags := []string{"Main", "Outline", "Alt"}

	for _, tag := range tags {
		assetName := string(label) + tag
		img, err := util.LoadImageAssetAsEbitenImage("uiSprites/" + assetName)
		if err != nil {
			log.Printf("%s not found for loading UiSpriteData %s, proceeding with loading other files. error msg: %s", assetName, string(label), err)
		} else {
			imgs = append(imgs, img)
		}
	}
	return imgs, nil
}

func LoadUISprites(spritesToLoad []entities.Label, environment *system.Environment, tankBounds image.Rectangle, hub *tasks.EventHub) ([]*entities.Entity, *entities.WhiteBoardSprite, error) {
	var uiEntities []*entities.Entity
	var wbSprite *entities.WhiteBoardSprite

	spritePositions, err := LoadSpritePositionData()
	if err != nil {
		return nil, nil, err
	}

	for i, elem := range spritesToLoad {

		if spritePositions[string(elem)] == nil {
			log.Fatal("No position data for sprite", elem)
		}

		x := spritePositions[string(elem)].X
		y := spritePositions[string(elem)].Y

		imgs, err2 := loadUiSpritesImgs(elem)
		if err2 != nil {
			return uiEntities, wbSprite, err
		}

		uSprite := entities.NewUiSprite(environment, imgs, hub, x, y, string(elem))

		uSprite.ActivationRect = image.Rect(tankBounds.Min.X, tankBounds.Min.Y+50, tankBounds.Max.X, tankBounds.Min.Y-200)

		entity := &entities.Entity{UiData: uSprite, Sprite: uSprite.Sprite}

		uSprite.LayerIndex = i //well just order by the order we load these cus we lazy af
		entity.EventHub = hub
		entities.RegisterEntity(entity)
		entities.UiSpriteSubs(hub, entity)
		switch elem {

		case entities.WhiteBoard:
			uSprite.Unfocusable = true
			wbSprite = &entities.WhiteBoardSprite{UiSprite: entity}
			wbSprite.Subscribe(hub)
			wbSprite.Init(hub)

		case entities.FishFood:
			uSprite.ActivationRect = image.Rect(tankBounds.Min.X, tankBounds.Min.Y-50, tankBounds.Max.X, tankBounds.Min.Y-200)
		case entities.Pillow, entities.Door:
			entity.Draw = false
		case entities.Phreader:
			entity.Sprite.AbleToBeUnfocusedAutomatically = true
		case entities.PiggyBank:
			entity.Sprite.AbleToBeUnfocusedAutomatically = true
			entity.AnimationMap = LoadPiggyBankAnimationMap(uSprite.X, uSprite.Y, float32(uSprite.Img.Bounds().Dy()))
		case entities.Magazine:
			entity.Draw = false
		}

		//lightingShader := shaders.LoadOnePointLightingNeutral()
		//sprite.Shader = lightingShader
		//LoadSpriteLightingParams(sprite.Sprite)

		uiEntities = append(uiEntities, entity)

	}

	return uiEntities, wbSprite, nil

}

func LoadSpritePositionData() (map[string]*drawables.SavePositionData, error) {
	var positions = make(map[string]*drawables.SavePositionData)
	spritePosition, err := assets.DataDir.ReadFile("data/spritePosition.json")
	if err != nil {
		return positions, err
	}
	json.Unmarshal(spritePosition, &positions)
	return positions, nil
}

func LoadPiggyBankAnimationMap(x, y, srcImageHeight float32) map[string]*sprite.Sprite {

	aniMap := make(map[string]*sprite.Sprite)
	img, err := util.LoadImageAssetAsEbitenImage("uiSprites/allowanceCollectedAni")
	if err != nil {
		log.Fatal(err, "cant load piggy bank animation thing")
	}
	animation := animations.NewAnimation(0, 7, 1, 5)
	spriteSheet := spritesheet.NewSpritesheet(8, 1, 149, 202)

	animatedSprite := sprite.Sprite{}
	animatedSprite.Img = img
	animatedSprite.Animation = animation
	animatedSprite.SpriteSheet = spriteSheet
	animatedSprite.ShaderParams = make(map[string]any)

	animatedSprite.X = x
	yOffSet := float32(animatedSprite.Img.Bounds().Dy()) - srcImageHeight
	animatedSprite.Y = y - yOffSet

	aniMap["allowance"] = &animatedSprite
	return aniMap
}

//imgLoader() []*ebitenImgs
//jsonLoader
