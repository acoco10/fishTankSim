package loader

import (
	"encoding/json"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/stringConstants"
	"github.com/acoco10/fishTankWebGame/game/system"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

func LoadUiSpritesImgs(label entities.Label) ([]*ebiten.Image, error) {
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
	loaded := make(map[string]uint32)
	for i, elem := range spritesToLoad {

		var spritePosition drawables.SavePositionData
		if spritePositions[string(elem)] == nil {
			//if we have a new sprite, just load it at whiteboard
			spritePosition = *spritePositions[string(entities.WhiteBoard)]
		} else {
			spritePosition = *spritePositions[string(elem)]
		}

		x := spritePosition.X
		y := spritePosition.Y

		imgs, err2 := LoadUiSpritesImgs(elem)
		if err2 != nil {
			return uiEntities, wbSprite, err
		}

		uSprite := entities.NewUiSprite(environment, imgs, hub, x, y, string(elem))
		uSprite.Sprite.DOptsUpdaterParams = make(map[string]float64)
		uSprite.ActivationRect = image.Rect(tankBounds.Min.X, tankBounds.Min.Y+50, tankBounds.Max.X, tankBounds.Min.Y-200)

		entity := &entities.Entity{UiData: uSprite, Sprite: uSprite.Sprite}

		entity.EventHub = hub
		entity.LayerIndex = i
		entities.RegisterEntity(entity)
		entities.UiSpriteSubs(hub, entity)

		entity.DoAt = make(map[string]func(entity *entities.Entity, gs entities.GameState))
		loaded[string(elem)] = entity.Id
		switch elem {

		case entities.WhiteBoard:
			uSprite.UnFocusable = true
			wbSprite = &entities.WhiteBoardSprite{UiSprite: entity}
			wbSprite.Subscribe(hub)
			wbSprite.Init(hub)
		case entities.LightSwitch:
			entity.Z = 13
			uSprite.UnFocusable = true
			uSprite.Flags["noOffset"] = true
		case entities.FishFood:
			uSprite.ActivationRect = image.Rect(tankBounds.Min.X, tankBounds.Min.Y-80, tankBounds.Max.X, tankBounds.Min.Y-200)
		case entities.Skimmer:
			sm := entities.InitStateMachine(nil, entities.AltImageWhenClickedUpdater, entities.AddUiSpriteXYUpdater, nil)
			//phys := physics.NewNetBody(float64(uSprite.Sprite.X), float64(int(uSprite.Sprite.Y)+uSprite.Sprite.GetSpriteRect().Dy()), 40, -math.Pi/2, math.Pi/2)
			entity.StateMachine = sm
		case entities.Pillow:
			entity.Draw = false
			uSprite.UnFocusable = true
		case entities.Phreader:
			sm := entities.InitStateMachine(nil, entities.AltImageWhenClickedUpdater, nil, nil)
			sm.AppendState(entities.ActivationRectUpdater, nil)
			sm.AppendState(entities.UsedInActivationRect, nil)
			entity.StateMachine = sm
			entity.DoAt = make(map[string]func(entity *entities.Entity, gs entities.GameState))
			entity.DoAt[entities.DoAtTime] = entities.PhReaderDoAtTime
			entity.UiData.Timers[entities.DoAtTime] = util.NewTimer(2.0)

			entity.UiData.Flags[stringConstants.Swirl] = true
			entity.UiData.Flags["updater"] = true
			entity.UiData.Flags["outline"] = true
			entity.UiData.Flags["autoTransition1"] = true
			entity.Sprite.AbleToBeUnfocusedAutomatically = true
		case entities.PiggyBank:
			entity.Sprite.AbleToBeUnfocusedAutomatically = true
			b := uSprite.Img.Bounds()
			entity.Sprite.AnimationMap = LoadPiggyBankAnimationMap(b.Dy(), b.Dx())
			sm := entities.InitPiggyBankStateMachine()
			entity.StateMachine = sm
			entity.UiData.Flags["lowLight"] = true
			entity.DoAt[entities.DoAtHovered] = entities.PublishGraphicHovered
		case entities.Thermometer:
			entity.Z = 11
			entity.UiData.BaseZ = 11
			sm := entities.InitStateMachine(nil, entities.AltImageWhenClickedUpdater, entities.AddTempGuage, nil)
			sm.AppendState(entities.TurnOffEveryThingOnUnFocus, nil)
			entity.StateMachine = sm
			entity.UiData.Flags["autoTransition1"] = true
			entity.UiData.Flags["updater"] = false
		case entities.Magazine:
			sm := entities.InitStateMachine(entities.DisabledState, entities.AltImageHovered,
				nil, entities.PublishPickedUpEvent)
			entity.Draw = false
			entity.UiData.Timers["freeze"] = util.NewTimer(0.5)
			entity.UiData.Flags["revert"] = true
			entity.Sprite.UnFocusable = true
			entity.UiData.DontDraw = true
			entity.StateMachine = sm

		case entities.GrandpasJournal:
			sm := entities.InitStateMachine(nil, entities.AltImageWhenClickedUpdater,
				nil, nil)
			sm.AppendState(entities.AltImageHovered, entities.PublishPickedUpEvent)
			entity.UiData.Timers["freeze"] = util.NewTimer(0.5)
			entity.StateMachine = sm
			entity.UiData.Flags["clickTransition"] = true
			entity.UiData.Flags["revert"] = true
			entity.Sprite.NoShaderOnFocus = true
		}

		uiEntities = append(uiEntities, entity)

	}

	//set linked ids now that all ids are loaded
	ent, exists := entities.GetEntity(loaded[string(entities.GrandpasJournal)])
	if exists {
		ent.LinkedID = loaded[string(entities.Magazine)]
	}

	ent2, exists := entities.GetEntity(loaded[string(entities.Magazine)])
	if exists {
		ent2.LinkedID = loaded[string(entities.GrandpasJournal)]
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

func LoadPiggyBankAnimationMap(srcImgHeight int, srcImgWidth int) map[string]*sprite.Animation {

	aniMap := make(map[string]*sprite.Animation)
	img, err := util.LoadImageAssetAsEbitenImage("uiSprites/allowanceCollectedAni")
	if err != nil {
		log.Fatal(err, "cant load piggy bank animation thing")
	}
	spriteSheet := sprite.NewSpriteSheet(8, 1, 149, 202)
	animation := sprite.NewAnimation(spriteSheet, 0, 7, 1, 5)
	animation.Img = img
	animation.OffSetY = float32(img.Bounds().Dy() - srcImgHeight)
	animation.OffSetX = float32(img.Bounds().Dx()/animation.LastF + 1 - srcImgWidth)
	aniMap["allowance"] = animation
	return aniMap
}

//imgLoader() []*ebitenImgs
//jsonLoader
