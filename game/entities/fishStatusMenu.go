package entities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/entImportableLoaders"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/stringConstants"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
)

const (
	StatusSpriteX = 40
	StatusSpriteY = 155
)

const (
	goldFishTitle  = "Goldfish"
	guppyFishTitle = "Guppy"
	angelFishTitle = "Angel Fish"
	kirbensisTile  = "Kirbensis"
	mollyFishTitle = "Common Molly"
)

var FishTitles = map[FishList]string{
	GoldFish:       goldFishTitle,
	MollyFish:      mollyFishTitle,
	AngelFish:      angelFishTitle,
	Kirbensis:      kirbensisTile,
	mollyFishTitle: mollyFishTitle,
}

func LoadFishFactorsMenu(data CreatureData, fishId uint32) {
	img, err := util.LoadImageAssetAsEbitenImage("menuAssets/status")
	if err != nil {
		log.Fatal(err)
	}

	x := float32(260)
	y := float32(120)

	menuSp := sprite.Sprite{Img: img, X: 0, Y: 0, Scale: 2}

	buffEnt := &Entity{Sprite: &sprite.Sprite{Img: menuSp.MakeSpriteSizedBuffer(2), X: x, Y: y}, Z: 1}
	buffEnt.Sprite.IsBuffer = true
	buffEnt.Sprite.UnFocusable = true

	menuEnt := &Entity{Sprite: &menuSp, Z: 1}
	menuEnt.Sprite.BufferDst = buffEnt.Sprite.Img

	sg := entImportableLoaders.LoadEffect("stomach")

	sg.BufferDst = buffEnt.Sprite.Img
	currentFrameShouldBe := (data.Hunger) * 3
	if data.HealthState == Stressed {
		sg = entImportableLoaders.LoadEffect("stressedStomach")
		currentFrameShouldBe = (data.Hunger - 3) * 3
	}
	sg.Scale = 2.0
	sg.X = StatusSpriteX
	sg.Y = StatusSpriteY

	sg.SetFrame(currentFrameShouldBe)

	if currentFrameShouldBe > sg.GetAnimation().LastF {
		sg.SetFrame(sg.GetAnimation().LastF)
	}

	sg.AnimationPaused = true
	stomachEnt := &Entity{Sprite: sg}
	stomachEnt.Sprite.GetAnimation().SpeedInTPS = 20
	stomachEnt.Sprite.Scale = 2
	stomachEnt.Parameters.Ints[stopFrame] = currentFrameShouldBe

	growthSprite := entImportableLoaders.LoadEffect("growth")
	growthSprite.AnimationPaused = true
	growthSprite.X = StatusSpriteX + float32(growthSprite.GetSpriteRect().Dx()) + 40
	growthSprite.Y = StatusSpriteY + 4
	growthSprite.Scale = 2
	growthEnt := &Entity{Sprite: growthSprite}
	growthEnt.Sprite.BufferDst = buffEnt.Sprite.Img

	currentFrameForProgress := int(data.progress)
	if currentFrameForProgress > growthSprite.GetAnimation().LastF {
		currentFrameForProgress = growthSprite.GetAnimation().LastF
	}

	growthEnt.Parameters.Ints[stopFrame] = currentFrameForProgress
	buffEnt.Parameters.EntIds[linkedCreature] = fishId
	buffEnt.SetOverZoom()
	RegisterEntity(stomachEnt)
	RegisterEntity(buffEnt)
	RegisterEntity(menuEnt)
	RegisterEntity(growthEnt)

	buffEnt.Parameters.EntIds[LinkedGraphic1] = growthEnt.Id

	buffEnt.DrawTimeFunc = AddTextToSpriteUpdater

	buffEnt.Sprite.DOptsUpdaterTag = stringConstants.Sway

	buffEnt.UpdateFunc = PlayAnimationUsingEntAsBufferUntilStopFrame
	buffEnt.LinkedID = stomachEnt.Id
	buffEnt.Parameters.Ints[AltSpriteWidth] = stomachEnt.Sprite.GetSpriteRect().Dx()
	buffEnt.Parameters.Ints[AltSpriteHeight] = stomachEnt.Sprite.GetSpriteRect().Dy()

	buffEnt.EventHub.Subscribe(events.UnFocusEvent{}, func(e tasks.Event) {
		RemoveEntity(buffEnt.Id)
		RemoveEntity(menuEnt.Id)
		RemoveEntity(stomachEnt.Id)
		RemoveEntity(growthEnt.Id)
	})

	menuEnt.EventHub.Subscribe(CreatureReachedPoint{}, func(e tasks.Event) {
		ev := e.(CreatureReachedPoint)

		if ev.PointTypeReached == util.Food && ev.CreatureID == fishId {
			if stomachEnt.Parameters.Ints[stopFrame] < stomachEnt.Sprite.GetAnimation().LastF {
				stomachEnt.Sprite.AnimationPaused = false
				stomachEnt.Parameters.Ints[stopFrame] += 3
			}
		}
	})
}

func PlayAnimationUsingEntAsBufferUntilStopFrame(ent *Entity, gs GameState) {
	ent2, exists := ent.GetLinkedEnt()
	if !exists {
		log.Printf("buffered animation stop frame script got linked sprite ent that doesnt exist at update time")
		return
	}
	stopF := ent2.Parameters.Ints[stopFrame]
	if ent2.Sprite.GetAnimation().Frame() >= stopF {
		ent2.Sprite.AnimationPaused = true
	}

}

func AddTextToSpriteUpdater(ent *Entity, state GameState) {
	fish, exists := GetEntity(ent.Parameters.EntIds[linkedCreature])
	if !exists {
		log.Printf("fish status fish id not working for get entity")
		return
	}

	data := fish.CreatureData.FishStats
	//176 242 115
	cs := util.ConvertRGBAtoEbitenCS(color.RGBA{176, 242, 115, 255})
	cs.ScaleAlpha(0.8)
	fontName := "RockSalt_16"

	face := registry.FontMap[fontName]
	if data.name == "" {
		data.name = "Fish With No Name"
	}

	w, h := util.MeasureText(data.name, 18, "RockSalt_18")

	offSet := h*2 + 4

	baseX := float64(ent.Sprite.Img.Bounds().Dx() / 16)
	baseY := h - 4
	centerX := float64(ent.Sprite.GetSpriteRect().Dx()) / 2

	tOpts := &text.DrawOptions{}
	tOpts.ColorScale = cs
	tOpts.GeoM.Translate(centerX-w/2, baseY)
	text.Draw(ent.Sprite.Img, data.name, registry.FontMap["RockSalt_18"], tOpts)

	tOpts.GeoM.Reset()
	tOpts.ColorScale = cs
	tOpts.GeoM.Translate(baseX, baseY+offSet)
	text.Draw(ent.Sprite.Img, FishTitles[data.FishType], face, tOpts)

	tOpts.GeoM.Reset()
	tOpts.ColorScale = cs
	tOpts.GeoM.Translate(baseX, baseY+offSet+h)
	text.Draw(ent.Sprite.Img, fmt.Sprintf("Age:%d", data.age), face, tOpts)

	tOpts.GeoM.Reset()
	tOpts.ColorScale = cs
	tOpts.GeoM.Translate(baseX, baseY+offSet+h*2)
	text.Draw(ent.Sprite.Img, fmt.Sprintf("Size:%d", data.Size), face, tOpts)

	tOpts.GeoM.Reset()
	tOpts.ColorScale = cs
	tOpts.GeoM.Translate(baseX, baseY+offSet+h*8)
	text.Draw(ent.Sprite.Img, "Hunger", registry.FontMap["RockSalt_12"], tOpts)

	tOpts.GeoM.Reset()
	tOpts.ColorScale = cs
	tOpts.GeoM.Translate(baseX+120, baseY+offSet+h*8)
	text.Draw(ent.Sprite.Img, "Growth", registry.FontMap["RockSalt_12"], tOpts)

	prog, exists := GetEntity(ent.Parameters.EntIds[LinkedGraphic1])
	if !exists {
		log.Fatal("linked progress graphic to menu entity for fish status doesnt exist")
	}

	prog.Sprite.SetFrame(int(data.progress))

}
