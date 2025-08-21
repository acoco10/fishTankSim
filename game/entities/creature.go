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

func LoadFishImg(fType FishList, level int, tag string) (*ebiten.Image, error) {
	var fishImgName string
	switch fType {
	case Fish:
		fishImgName = fmt.Sprintf("fish%d%sSpriteSheet", level, tag)
	case MollyFish:
		fishImgName = fmt.Sprintf("mollyFish%d%sSpriteSheet", level, tag)
	case Guppy:
		fishImgName = fmt.Sprintf("guppy%d%sSpriteSheet", level, tag)
	case Kirbensis:
		fishImgName = fmt.Sprintf("kirbensis%d%sSpriteSheet", level, tag)
	}
	img, err := util.LoadImageAssetAsEbitenImage("fishSpriteSheets/" + fishImgName)
	if err != nil {
		return &ebiten.Image{}, err
	}
	return img, nil
}

func LoadFishNormal(fType FishList, level int, tag string) (*ebiten.Image, error) {
	normalImgPath := fmt.Sprintf("fishSpriteSheets/assets/images/fishSpriteSheets/%s%dSpriteSheet_n", string(fType), level)
	img, err := util.LoadImageAssetAsEbitenImage(normalImgPath)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func LoadFishSprite(creatureType FishList, creatureLvl int) (map[string]*sprite.Sprite, error) {

	var c *sprite.Sprite

	c = &sprite.Sprite{}
	eatingC := &sprite.Sprite{}
	c.Scale = 1

	img, err := LoadFishImg(creatureType, creatureLvl, "")
	if err != nil {
		return nil, err
	}

	normalImg, err := LoadFishNormal(creatureType, creatureLvl, "")
	if err != nil {
		return nil, err
	}
	c.Img = img

	tag := "Eating"
	eatingImg, err := LoadFishImg(creatureType, creatureLvl, tag)
	if err != nil {
		eatingImg, err = LoadFishImg(MollyFish, 1, tag)
		if err != nil {
			log.Fatal("this is the you fucked up path", err)
		}
	}

	eatingNormalImg, _ := LoadFishNormal(creatureType, creatureLvl, tag)

	c.Img = img
	eatingC.Img = eatingImg

	path := fmt.Sprintf("data/animationData/%s%dAnimation.json", string(creatureType), creatureLvl)
	fmt.Printf("Animation path = %s\n", path)

	eatingPath := fmt.Sprintf("data/animationData/%s%d%sAnimation.json", string(creatureType), creatureLvl, tag)
	fmt.Printf("Animation path = %s\n", path)

	eani, esps, err := entImportableLoaders.LoadAnimation(eatingPath)
	if err != nil {
		log.Printf("wrong animation path or misnamed file")
		eani, esps, err = entImportableLoaders.LoadAnimation("data/animationData/mollyFish1EatingAnimation.json")
		if err != nil {
			log.Fatal("Even the backup doesnt work", err)
		}
	}

	ani, sps, err := entImportableLoaders.LoadAnimation(path)
	if err != nil {
		log.Printf("wrong animation path or misnamed file")
		return nil, err
	}
	c.Animation = ani
	c.SpriteSheet = sps

	eatingC.SpriteSheet = esps
	eatingC.Animation = eani
	eatingC.NormalMap = eatingNormalImg
	eatingC.ShaderParams = make(map[string]any)

	if normalImg != nil {
		c.NormalMap = normalImg
		c.ShaderParams = make(map[string]any)
		c.ShaderParams["Cursor"] = [2]float64{400, 50}
		c.ShaderParams["TankDepthZ"] = rand.Float64() * 0.5
		c.Shader = registry.ShaderMap["NormalMap"]
	}

	AnimationMap := make(map[string]*sprite.Sprite)
	AnimationMap["swimming"] = c
	AnimationMap["eating"] = eatingC

	c.AbleToBeUnfocusedAutomatically = true

	return AnimationMap, nil
}

func LoadFishSpriteAltAnimations(fType FishList) (*sprite.AnimatedSprite, error) {
	c := sprite.AnimatedSprite{}
	c.Sprite = &sprite.Sprite{}

	c.Animation = animations.NewAnimation(0, 3, 1, 15)

	return &c, nil
}

func NewFishData(environment *system.Environment, hub *tasks.EventHub, tankSize image.Rectangle, saveData SavedFish) *CreatureData {

	timers := make(map[FishState]*util.Timer)
	randDuration := rand.Float64() * 50
	timers[Swimming] = util.NewTimer(randDuration)
	timers[Eating] = util.NewTimer(1.5)
	timers[Resting] = util.NewTimer(10)

	fs, err := GenFishStats(FishList(saveData.FishType), saveData.Name)
	if err != nil {
		log.Fatal("fish stats", err)
	}
	if fs == nil {
		println("Fish stats returning empty pointer")
	}

	fs.Size = saveData.Size

	tankSize.Max.X -= 5
	tankSize.Min.X += 5

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

func LoadLevelUpSprite(creature *Entity) {
	if creature.CreatureData.Size < 4 {
		levelUpSprite, err := LoadFishSprite(creature.CreatureData.FishType, creature.CreatureData.Size)
		if err != nil {
			//just load goldfish img as default level up sprite
			levelUpSprite, _ = LoadFishSprite("Fish", creature.CreatureData.Size)
		}
		ms := levelUpSprite["swimming"]
		creature.Sprite = ms
		creature.AnimationMap = levelUpSprite
	}
}
