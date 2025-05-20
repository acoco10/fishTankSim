package sprite

import (
	"encoding/json"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

type UISpriteLabel string

const (
	Records    UISpriteLabel = "records"
	FishFood   UISpriteLabel = "fishFood"
	FishBook   UISpriteLabel = "book"
	WhiteBoard UISpriteLabel = "whiteBoard"
	Plant      UISpriteLabel = "plant"
)

func loadUiSpritesImgs(label UISpriteLabel) ([]*ebiten.Image, error) {
	var imgs []*ebiten.Image
	tags := []string{"Main", "Outline", "Alt"}

	for _, tag := range tags {
		assetName := string(label) + tag
		img, err := loader.LoadImageAssetAsEbitenImage("uiSprites/" + assetName)
		if err != nil {
			log.Printf("%s not found for loading UiSprite %s, proceeding with loading other files. error msg: %s", assetName, string(label), err)
		} else {
			imgs = append(imgs, img)
		}
	}
	return imgs, nil
}

func LoadUISprites(spritesToLoad []UISpriteLabel, hub *events.EventHub, screenWidth, screenHeight int) ([]DrawableSprite, error) {
	var sprites []DrawableSprite

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
		sprite := NewUiSprite(imgs, hub, x, y, string(elem), screenWidth, screenHeight)
		if elem == FishFood {
			ffSprite := FishFoodSprite{sprite}
			ffSprite.Subscribe()
			sprites = append(sprites, &ffSprite)
			continue
		}
		if elem == WhiteBoard {
			wbSprite := WhiteBoardSprite{UiSprite: sprite}
			wbSprite.Subscribe(hub)
			sprites = append(sprites, &wbSprite)
			continue
		}
		sprites = append(sprites, sprite)

	}

	return sprites, nil

}

func loadSpritePositionData() (map[string]*SavePositionData, error) {
	var positions = make(map[string]*SavePositionData)
	spritePosition, err := assets.DataDir.ReadFile("data/spritePosition.json")
	if err != nil {
		return positions, err
	}
	json.Unmarshal(spritePosition, &positions)
	return positions, nil
}

//imgLoader() []*ebitenImgs
//jsonLoader
