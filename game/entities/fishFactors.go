package entities

import (
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
	clr := colornames.Gray
	cs.Scale(float32(clr.R), float32(clr.G), float32(clr.B), float32(clr.A))
	name := "Lily"
	nameId := graphics.NewNkTextGraphic(&name, 22, float64(bground.Sprite.X+(135*2))/2, float64(bground.Sprite.Y+34)/2, false, cs, 1, false)

	bground.Sprite.PublishedGraphicId = append(bground.Sprite.PublishedGraphicId, nameId)

	fish.LinkedID = bground.Id

	bground.Sprite.Z = 4

	icons, err := util.LoadImageAssetAsEbitenImage("uiSprites/fishFactorIcons")
	if err != nil {
		log.Fatal(err)
	}

	iconLabels := []string{"thumbsUp", "thumbsNeutral", "thumbsDown", "otherFish", "structures", "temperature", "ph"}
	imageMap, indMap := ChopUpIcons(icons, iconLabels, 32)

	iconLabels = iconLabels[3:]

	buffer := float32(60.0 * registry.Config.ZoomFactor)
	spacing := float32(64.0 * registry.Config.ZoomFactor)

	lastSprite := &Entity{}
	for i, label := range iconLabels {
		iconSprite := &Entity{}
		if i < 2 {
			iconSprite = MakeSpriteEntity(imageMap[label], bground.Sprite.X+buffer+(float32(i+1)*spacing), bground.Sprite.Y+80)
		} else {
			iconSprite = MakeSpriteEntity(imageMap[label], bground.Sprite.X+buffer+float32(i-1)*spacing, bground.Sprite.Y+spacing+40)
		}
		iconSprite.Sprite.Z = 5
		switch label {
		case "temperature":
			threshHolds := []float64{5.0, 10.0}
			val := math.Abs(float64(fish.CreatureData.idealTemperature - fish.CreatureData.Environment.Temperature))
			condImg := CheckIconValue(val, threshHolds, indMap)
			iconSprite.Sprite.LinkedSprite = &sprite.Sprite{
				Img: &condImg,
				X:   iconSprite.Sprite.X + 64, Y: iconSprite.Sprite.Y, Z: iconSprite.Sprite.Z,
				Scale: registry.Config.ZoomFactor}
		case "ph":
			threshHolds := []float64{2.5, 5.0}
			val := math.Abs(float64(fish.CreatureData.idealPH - fish.CreatureData.Environment.NaturalPHLevel))
			condImg := CheckIconValue(val, threshHolds, indMap)
			iconSprite.Sprite.LinkedSprite = &sprite.Sprite{
				Img: &condImg,
				X:   iconSprite.Sprite.X + 64, Y: iconSprite.Sprite.Y, Z: iconSprite.Sprite.Z, Scale: registry.Config.ZoomFactor}
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

func ChopUpIcons(inputImage *ebiten.Image, labels []string, size int) (map[string]*ebiten.Image, map[string]*ebiten.Image) {
	imageMap := make(map[string]*ebiten.Image)
	indMap := make(map[string]*ebiten.Image)

	for i, icon := range labels {
		//horizontal slice of square images
		rect := image.Rect(i*size, 0, (i+1)*size, size)

		if i < 3 {
			indMap[icon] = ebiten.NewImageFromImage(inputImage.SubImage(rect))
		} else {
			imageMap[icon] = ebiten.NewImageFromImage(inputImage.SubImage(rect))
		}
	}

	return imageMap, indMap
}

func InitStomachGraphic(menuBackground *Entity, fishId uint32) uint32 {
	stomachSprite := LoadStomachGraphic()
	stomachSprite.X = menuBackground.Sprite.X + 40
	stomachSprite.Y = menuBackground.Sprite.Y + 80
	stomachSprite.Shader = registry.ShaderMap["Stomach"]
	stomachSprite.ShaderParams = make(map[string]any)
	stomachSprite.ShaderParams["Fullness"] = 0.0
	stomachSprite.ShaderParams["FishId"] = fishId
	stomachSprite.UpdateShaderParams = UpdateFullness
	stomachEnt := &Entity{Sprite: stomachSprite}
	stomachEnt.Sprite.Z = 5
	RegisterEntity(stomachEnt)
	return stomachEnt.Id
}

func UpdateFullness(params map[string]any) map[string]any {

	entId := params["FishId"].(uint32)
	targetedFish, exists := GetEntity(entId)
	if !exists {
		log.Println("fullness graphic trying to update with nil fish entity")
		return params
	}

	params["Fullness"] = float64(targetedFish.CreatureData.Hunger) / float64(targetedFish.CreatureData.maxHunger)
	println("fullness =", params["Fullness"].(float64))

	return params
}
