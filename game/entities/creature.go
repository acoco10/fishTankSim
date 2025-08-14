package entities

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/fishTankWebGame/game/entImportableLoaders"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/system"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"math/rand"
)

func LoadFishImg(fType FishList, level int) (*ebiten.Image, error) {
	var fishImgName string
	switch fType {
	case Fish:
		fishImgName = fmt.Sprintf("fish%dSpriteSheet", level)
	case MollyFish:
		fishImgName = fmt.Sprintf("mollyFish%dSpriteSheet", level)
	case Guppy:
		fishImgName = fmt.Sprintf("guppy%dSpriteSheet", level)
	case Kirbensis:
		fishImgName = fmt.Sprintf("kirbensis%dSpriteSheet", level)
	}
	img, err := util.LoadImageAssetAsEbitenImage("fishSpriteSheets/" + fishImgName)
	if err != nil {
		return &ebiten.Image{}, err
	}
	return img, nil
}

func LoadFishNormal(fType FishList, level int) (*ebiten.Image, error) {
	normalImgPath := fmt.Sprintf("fishSpriteSheets/assets/images/fishSpriteSheets/%s%dSpriteSheet_n", string(fType), level)
	img, err := util.LoadImageAssetAsEbitenImage(normalImgPath)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func LoadFishSprite(creatureType FishList, creatureLvl int) (*sprite.Sprite, error) {

	var c *sprite.Sprite

	c = &sprite.Sprite{}
	c.Scale = 1

	img, err := LoadFishImg(creatureType, creatureLvl)
	if err != nil {
		return nil, err
	}

	normalImg, err := LoadFishNormal(creatureType, creatureLvl)
	if err != nil {
		return nil, err
	}
	c.Img = img
	path := fmt.Sprintf("data/animationData/%s%dAnimation.json", string(creatureType), creatureLvl)

	fmt.Printf("Animation path = %s\n", path)

	ani, sps, err := entImportableLoaders.LoadAnimation(path)
	if err != nil {
		log.Printf("wrong animation path or misnamed file")
		return nil, err
	}

	c.SpriteSheet = sps
	c.Animation = ani

	if normalImg != nil {
		c.NormalMap = normalImg
		c.ShaderParams = make(map[string]any)
		c.ShaderParams["Cursor"] = [2]float64{400, 50}
		c.Shader = registry.ShaderMap["NormalMap"]
	}

	c.AbleToBeUnfocusedAutomatically = true
	return c, nil
}

func LoadFishSpriteAltAnimations(fType FishList) (*sprite.AnimatedSprite, error) {
	c := sprite.AnimatedSprite{}
	c.Sprite = &sprite.Sprite{}

	switch fType {
	case MollyFish:
		println("Loading Molly Fish Animation")
		img, err := util.LoadImageAssetAsEbitenImage("fishSpriteSheets/mollyFishSpinAnimation")
		if err != nil {
			return &c, err
		}

		c.Img = img
	}

	c.Animation = animations.NewAnimation(0, 3, 1, 15)

	return &c, nil
}

func NewFishData(environment *system.Environment, hub *tasks.EventHub, tankSize image.Rectangle, saveData SavedFish) *CreatureData {

	timers := make(map[FishState]*util.Timer)
	randDuration := rand.Float64() * 50
	timers[Swimming] = util.NewTimer(randDuration)
	timers[Eating] = util.NewTimer(0.5)
	timers[Resting] = util.NewTimer(10)

	fs, err := GenFishStats(FishList(saveData.FishType), saveData.Name)
	if err != nil {
		log.Fatal("fish stats", err)
	}
	if fs == nil {
		println("Fish stats returning empty pointer")
	}

	fs.Size = saveData.Size

	c := CreatureData{

		TargetPoint:        nil,
		ParticlePointQueue: map[uint32]*util.Point{},
		EventHub:           hub,
		TankBoundaries:     tankSize,
		Timers:             timers,
		State:              Swimming,
		TickClicked:        false,
		Environment:        environment,
		FishStats:          fs,
		Flip:               false,
	}
	//LoadRotatingHighlightOutlineAnimated(c.AnimatedSprite)
	c.ParticlePointQueue = make(map[uint32]*util.Point)
	return &c
}

func LoadLevlUpSprite(creature *Entity) {
	if creature.CreatureData.Size < 4 {
		levelUpSprite, err := LoadFishSprite(creature.CreatureData.FishType, creature.CreatureData.Size)
		if err != nil {
			//just load goldfish img as default level up sprite
			levelUpSprite, _ = LoadFishSprite("Fish", creature.CreatureData.Size)
		}
		creature.Sprite = levelUpSprite
	}
}
