package entities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/entImportableLoaders"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/system"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
	"math/rand"
	"strings"
)

var animationPaths []string

func LoadFishImg(fType FishList, level int, tag string) (*ebiten.Image, error) {
	var fishImgName string
	switch fType {
	case GoldFish:
		fishImgName = fmt.Sprintf("goldFish%d%sSpriteSheet", level, tag)
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

func LoadFishAnimations(creatureType FishList, creatureLvl int) (map[string]*sprite.Animation, error) {

	AnimationMap := make(map[string]*sprite.Animation)

	tags := []string{"", "Eating", "Forward", "Backwards", "Depth"}
	for _, tag := range tags {
		nimg, err := LoadFishNormal(creatureType, creatureLvl, "")
		if err != nil {
			log.Println("normal map not found for fish:", creatureType, "ERROR:", err)
		}
		img := loadFishAnimationImg(creatureType, creatureLvl, tag)
		animation := loadFishAnimationData(creatureType, creatureLvl, tag)
		animation.Img = img
		animation.NormalImg = nimg

		if tag == "" {
			AnimationMap["swimming"] = animation
			continue
		}

		name := strings.ToLower(tag)
		AnimationMap[name] = animation
	}

	return AnimationMap, nil
}

func PopulateAnimationDirectory(tags []string) {

}

func loadFishAnimationData(creatureType FishList, creatureLvl int, tag string) *sprite.Animation {
	path := fmt.Sprintf("data/animationData/%s%d%sAnimation.json", string(creatureType), creatureLvl, tag)
	fmt.Printf("Animation path = %s\n", path)

	fallBackPath := fmt.Sprintf("data/animationData/%s%d%sAnimation.json", string(GoldFish), 1, tag)
	fmt.Printf("Animation path = %s\n", path)

	ani, err := entImportableLoaders.LoadAnimation(path)
	if err != nil {
		log.Printf("wrong animation path or misnamed file")
		ani, err = entImportableLoaders.LoadAnimation(fallBackPath)
		if err != nil {
			log.Fatal("Even the backup doesnt work", err)
		}
	}
	return ani
}

func loadFishAnimationImg(creatureType FishList, creatureLvl int, tag string) *ebiten.Image {
	img, err := LoadFishImg(creatureType, creatureLvl, tag)
	if err != nil {
		log.Println("Trying to load fish with level-1")
		img, err = LoadFishImg(GoldFish, 1, tag)
		if err != nil {
			log.Fatal("this is the you fucked up path", err)
		}
	}
	return img
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
		println("GoldFish stats returning empty pointer")
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

func LoadFishSprite(creature *Entity) {
	if creature.CreatureData.Size < 4 {
		//no level 4+ sprites yet
		Sprite, err := LoadFishAnimations(creature.CreatureData.FishType, creature.CreatureData.Size)
		if err != nil {
			//just load goldfish img as default level up sprite
			Sprite, _ = LoadFishAnimations("GoldFish", creature.CreatureData.Size)
		}
		creature.Sprite = &sprite.Sprite{CurrentAnimation: "swimming", Scale: 1, AbleToBeUnfocusedAutomatically: true, ShaderParams: make(map[string]any)}
		creature.Sprite.ShaderParams["Cursor"] = []float64{0, 0, 100}
		creature.Sprite.ShaderParams["TankDepthZ"] = creature.Z
		creature.Sprite.AnimationMap = Sprite
	}
}
