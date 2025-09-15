package entities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"image"
	"log"
	"math"
)

func MakeFishMenu(forFish uint32) {

	fish, exists := GetEntity(forFish)

	if !exists {
		log.Fatal("make fish menu called with entity that is nil or called with wrong id")
	}

	menuBackGround, err := util.LoadImageAssetAsEbitenImage("uiSprites/fishFactors")
	if err != nil {
		log.Fatal(err)
	}

	bground := MakeSpriteEntity(
		menuBackGround,
		float32((registry.Config.ResolutionWidth/2)-menuBackGround.Bounds().Dx()),
		float32(registry.Config.ResolutionHeight/4*3)-float32(menuBackGround.Bounds().Dy()))

	cs := ebiten.ColorScale{}
	clr := colornames.Gold
	cs.Scale(float32(clr.R), float32(clr.G), float32(clr.B), float32(clr.A))
	name := fish.CreatureData.name

	nameId := graphics.NewNkTextGraphic(&name, 16, float64(bground.Sprite.X+(264))/2, float64(bground.Sprite.Y+28)/2, false, cs, 1, false, 0)

	currrentLevel := fmt.Sprintf("Level: %d", fish.CreatureData.Size)
	currlvlId := graphics.NewNkTextGraphic(&currrentLevel, 12, float64(bground.Sprite.X+(70))/2, float64(bground.Sprite.Y+50)/2, false, cs, 1, false, 0)

	ageString := "Age: 1 day"
	ageId := graphics.NewNkTextGraphic(&ageString, 12, float64(bground.Sprite.X+(130))/2, float64(bground.Sprite.Y+50)/2, false, cs, 1, false, 0)

	msg := "Progress:"
	factorsID := graphics.NewNkTextGraphic(&msg, 12, float64(bground.Sprite.X+(62*2))/2, (float64(bground.Sprite.Y)+float64((bground.Sprite.Img.Bounds().Dy()-52)*2))/2, false, cs, 1, false, 0)

	bground.Sprite.PublishedGraphicId = append(bground.Sprite.PublishedGraphicId, nameId, factorsID, currlvlId, ageId)

	fish.LinkedID = bground.Id

	bground.NoZoom = true

	icons, err := util.LoadImageAssetAsEbitenImage("uiSprites/fishFactorIcons")
	if err != nil {
		log.Fatal(err)
	}

	iconLabels := []string{"thumbsUp", "thumbsNeutral", "thumbsDown", "otherFish", "structures", "temperature", "ph"}
	imageMap, indMap := util.ChopUpIcons(icons, iconLabels, 32)

	iconLabels = iconLabels[3:]
	bground.Z = 13
	//buffer := float32(60.0 * registry.Config.ZoomFactor)
	spacing := 84

	lastSprite := &Entity{}
	for i, label := range iconLabels {
		if label == "otherFish" || label == "structures" {
			continue
		}
		iconSprite := &Entity{}
		iconSprite = MakeSpriteEntity(imageMap[label], bground.Sprite.X+float32(bground.Sprite.GetSpriteRect().Dx()-120+i*spacing), bground.Sprite.Y+float32(80))
		iconSprite.Z = 14
		iconSprite.NoZoom = true
		switch label {
		case "temperature":
			threshHolds := []float64{5.0, 10.0}
			val := math.Abs(float64(fish.CreatureData.IdealTemperature - fish.CreatureData.Environment.Temperature))
			condImg := CheckIconValue(val, threshHolds, indMap)
			iconSprite.Sprite.LinkedSprite = &sprite.Sprite{
				Img:   &condImg,
				X:     iconSprite.Sprite.X + 45,
				Y:     iconSprite.Sprite.Y,
				Scale: registry.Config.ZoomFactor}
		case "ph":
			threshHolds := []float64{0.25, 0.5}
			val := math.Abs(float64(fish.CreatureData.IdealPH - fish.CreatureData.Environment.ModifiedPHLevel))
			condImg := CheckIconValue(val, threshHolds, indMap)
			iconSprite.Sprite.LinkedSprite = &sprite.Sprite{
				Img: &condImg,
				X:   iconSprite.Sprite.X + 35, Y: iconSprite.Sprite.Y,
				Scale: registry.Config.ZoomFactor}
		case "otherFish":
		case "structures":
		}

		if bground.LinkedID == 0 {
			bground.LinkedID = iconSprite.Id
		} else {
			lastSprite.LinkedID = iconSprite.Id
		}

		lastSprite = iconSprite
	}

	x0 := int(bground.Sprite.X + 92)
	y0 := int(bground.Sprite.Y) + bground.Sprite.GetSpriteRect().Dy()*2 - 95
	bar := image.Rect(x0, y0, x0+50, y0+12)
	graphics.NewRectGraphic(bar, colornames.Red)
	//todo abstract the stupid progress bar
	width := int((fish.CreatureData.progress / fish.CreatureData.nextLevel) * 50)
	println("rect inned width:", width)

	innerBar := image.Rect(x0+2, y0+2, x0+2+width, y0+10)
	graphics.NewFilledRectGraphicWithPointerWidth(innerBar, colornames.Green, &fish.CreatureData.progress, fish.CreatureData.nextLevel)

	stomachId := InitStomachGraphic(bground, fish.Id)

	lastSprite.LinkedID = stomachId

	ZSortEntities()

}

func CheckIconValue(val float64, thresholds []float64, iconMap map[string]*ebiten.Image) ebiten.Image {
	if val < thresholds[0] {
		return *iconMap["thumbsUp"]
	}
	if val < thresholds[1] {
		return *iconMap["thumbsNeutral"]
	}
	if val >= thresholds[1] {
		return *iconMap["thumbsDown"]
	}
	return *iconMap["thumbsNeutral"]
}

func InitStomachGraphic(menuBackground *Entity, fishId uint32) uint32 {
	stomachSprite := LoadStomachGraphic()
	stomachSprite.X = menuBackground.Sprite.X + 86
	stomachSprite.Y = menuBackground.Sprite.Y + 82
	stomachSprite.Shader = registry.ShaderMap["Stomach"]
	stomachSprite.ShaderParams = make(map[string]any)
	stomachSprite.ShaderParams["Fullness"] = 0.0
	stomachSprite.ShaderParams["FishId"] = fishId
	stomachSprite.UpdateShaderParams = UpdateFullness
	stomachEnt := &Entity{Sprite: stomachSprite}
	stomachEnt.Z = 13
	RegisterEntity(stomachEnt)
	stomachEnt.NoZoom = true
	return stomachEnt.Id
}

func UpdateFullness(params map[string]any) map[string]any {

	entId := params["FishId"].(uint32)
	targetedFish, exists := GetEntity(entId)
	if !exists {
		log.Println("fullness graphic trying to update with nil fish entity")
		return params
	}

	params["Fullness"] = float64(targetedFish.CreatureData.Hunger) / float64(targetedFish.CreatureData.MaxHunger)

	return params
}
