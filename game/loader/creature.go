package loader

import (
	"fmt"
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spriteSheet"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/geometry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"math/rand"
)

func LoadFishImg(fType entities.FishList, level int) (*ebiten.Image, error) {
	var fishImgName string
	switch fType {
	case entities.Fish:
		fishImgName = fmt.Sprintf("fish%dSpriteSheet", level)
	case entities.MollyFish:
		fishImgName = fmt.Sprintf("mollyFish%dSpriteSheet", level)
	case entities.Guppy:
		fishImgName = fmt.Sprintf("guppy%dSpriteSheet", level)
	}
	img, err := LoadImageAssetAsEbitenImage("fishSpriteSheets/" + fishImgName)
	if err != nil {
		return &ebiten.Image{}, err
	}
	return img, nil
}

func LoadFishNormal(fType entities.FishList, level int) (*ebiten.Image, error) {
	var fishImgName string
	switch fType {
	case entities.Fish:
		fishImgName = fmt.Sprintf("fish%dNormalSpriteSheet", level)
	case entities.MollyFish:
		fishImgName = fmt.Sprintf("mollyFish%dNormalSpriteSheet", level)
	case entities.Guppy:
		fishImgName = fmt.Sprintf("guppy%dSpriteSheet", level)

	}
	img, err := LoadImageAssetAsEbitenImage("fishSpriteSheets/" + fishImgName)

	if err != nil {
		log.Printf("normal map not found for: %s, trying with lazy export name", fishImgName)
		switch fType {
		case entities.Fish:
			fishImgName = fmt.Sprintf("fish%dSpriteSheet_n", level)
		case entities.MollyFish:
			fishImgName = fmt.Sprintf("mollyFish%dSpriteSheet_n", level)
		case entities.Guppy:
			fishImgName = fmt.Sprintf("guppy%dSpriteSheet_n", level)
		}
		img, err = LoadImageAssetAsEbitenImage("fishSpriteSheets/" + fishImgName)
		if err != nil {
			log.Printf("normal map not found for: %s", fishImgName)
			return nil, err
		}
	}

	return img, nil
}

func LoadFishSprite(creatureType entities.FishList, creatureLvl int) (*sprite.AnimatedSprite, error) {
	var c sprite.AnimatedSprite
	c.Sprite = &sprite.Sprite{}
	c.Scale = 1

	img, err := LoadFishImg(creatureType, creatureLvl)
	if err != nil {
		return nil, err
	}

	normalImg, err := LoadFishNormal(creatureType, creatureLvl)
	if err != nil {
		return nil, err
	}

	switch creatureLvl {
	case 1:
		c.Img = img
		c.Animation = animations.NewAnimation(0, 3, 1, 15)
		c.SpriteSheet = spritesheet.NewSpritesheet(4, 1, 21, 9)
		if normalImg != nil {
			c.NormalMap = normalImg
			c.ShaderParams = make(map[string]any)
			c.ShaderParams["Cursor"] = [2]float64{440, 170}
			shader := shaders.LoadNormalMapShader()
			c.Shader = shader
		}
	case 2:
		c.Img = img
		c.Animation = animations.NewAnimation(0, 3, 1, 15)
		c.SpriteSheet = spritesheet.NewSpritesheet(4, 1, 26, 13)
	case 3:
		c.Img = img
		c.Animation = animations.NewAnimation(0, 3, 1, 15)
		c.SpriteSheet = spritesheet.NewSpritesheet(4, 1, 65, 37)
		if normalImg != nil {
			c.NormalMap = normalImg
			c.ShaderParams = make(map[string]any)
			c.ShaderParams["Cursor"] = [2]float64{440, 170}
			shader := shaders.LoadNormalMapShader()
			c.Shader = shader
		}
	default:
		c.Img = img
		c.Animation = animations.NewAnimation(0, 3, 1, 15)
		c.SpriteSheet = spritesheet.NewSpritesheet(4, 1, 24, 11)
	}

	return &c, nil
}

func LoadFishSpriteAltAnimations(fType entities.FishList) (*sprite.AnimatedSprite, error) {
	c := sprite.AnimatedSprite{}
	c.Sprite = &sprite.Sprite{}

	switch fType {
	case entities.MollyFish:
		println("Loading Molly Fish Animation")
		img, err := LoadImageAssetAsEbitenImage("fishSpriteSheets/mollyFishSpinAnimation")
		if err != nil {
			return &c, err
		}

		c.Img = img
	}

	c.Animation = animations.NewAnimation(0, 3, 1, 15)
	c.SpriteSheet = spritesheet.NewSpritesheet(7, 1, 21, 9)

	return &c, nil
}

func NewFish(hub *tasks.EventHub, tankSize geometry.Rect, saveData entities.SavedFish) *entities.Creature {

	timers := make(map[entities.FishState]*entities.Timer)
	randDuration := rand.Float64() * 50
	timers[entities.Swimming] = entities.NewTimer(randDuration)
	timers[entities.Eating] = entities.NewTimer(0.5)
	timers[entities.Resting] = entities.NewTimer(10)

	fs, err := entities.GenFishStats(entities.FishList(saveData.FishType), saveData.Name)
	if err != nil {
		log.Fatal(err)
	}
	if fs == nil {
		println("Fish stats returning empty pointer")
	}

	fs.Size = saveData.Size

	c := entities.Creature{
		[]*geometry.Point{},
		hub,
		tankSize,
		timers,
		entities.Swimming,
		false,
		false,
		fs,
		sprite.NewAnimatedSprite(),
		false,
	}

	c.AnimatedSprite, err = LoadFishSprite(c.FishType, c.Size)
	if err != nil {
		log.Fatal(err)
	}

	//LoadRotatingHighlightOutlineAnimated(c.AnimatedSprite)

	firstPoint := c.RandomTarget()
	c.AddTargetPointToQueue(firstPoint)

	c.X = rand.Float32()*200 + tankSize.X1
	c.Y = rand.Float32()*100 + tankSize.Y1

	entities.CreatureEventSubscriptions(&c)

	return &c
}

func LoadLevlUpSprite(c *entities.Creature) {
	newFish, err := LoadFishSprite(c.FishType, c.Size)

	log.Printf("loading fish type: %s to level:%d", string(c.FishType), c.Size)

	if err != nil {
		log.Fatal("tried to lvl up a fish that dont exish")
	}

	c.AnimatedSprite = newFish
}
