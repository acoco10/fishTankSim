package gameEntities

import (
	"encoding/json"
	"fishTankWebGame/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

type UISpriteLabel string

const (
	Records  UISpriteLabel = "records"
	FishFood UISpriteLabel = "fishFood"
	FishBook UISpriteLabel = "book"
)

var uiElements = []UISpriteLabel{Records, FishBook, FishFood}

func loadUiSpritesImgs(label UISpriteLabel) ([]*ebiten.Image, error) {
	var imgs []*ebiten.Image
	tags := []string{"Main", "Outline", "Alt"}

	for _, tag := range tags {
		assetName := string(label) + tag
		img, err := LoadImageAssetAsEbitenImage("uiSprites/" + assetName)
		if err != nil {
			log.Printf("%s not found for loading uiSprite %s, proceeding with loading other files. error msg: %s", assetName, string(label), err)
		} else {
			imgs = append(imgs, img)
		}
	}
	return imgs, nil
}

func LoadUISprites(hub EventHub, screenWidth, screenHeight int) ([]DrawableSprite, error) {
	var sprites []DrawableSprite

	spritePositions, err := loadSpritePositionData()
	if err != nil {
		return nil, err
	}

	for _, elem := range uiElements {
		x := spritePositions[string(elem)].X
		y := spritePositions[string(elem)].Y
		imgs, err := loadUiSpritesImgs(elem)
		if err != nil {
			return sprites, err
		}
		sprite := NewUiSprite(imgs, &hub, x, y, string(elem), screenWidth, screenHeight)
		if elem == FishFood {
			s2 := FishFoodSprite{sprite}
			sprites = append(sprites, &s2)
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
