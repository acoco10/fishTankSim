package entities

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

func PHReaderUpdate(e *Entity, gs GameState) {
	us := e.Sprite
	uiDat := e.UiData

	b := us.Img.Bounds()

	pt := image.Point{X: int(us.X) + b.Dx()/2, Y: int(us.Y) + b.Dy()/2}

	if !pt.In(uiDat.ActivationRect) {
		us.ShaderParams["OutlineColor"] = [4]float64{1, 1, 0, 1}
	}

	if uiDat.state == Selected && pt.In(uiDat.ActivationRect) {
		us.ShaderParams["OutlineColor"] = [4]float64{0, 1, 0, 1}
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			if e.Z > 10 {
				UpdateEntityZAndReSortEntitySlice(e.Id, 10)
				e.Sprite.DOptsUpdaterTag = "swirl"
				phps := NewPhReaderParticleSystem(float64(e.Sprite.X)+float64(e.Sprite.GetSpriteRect().Dx()/2), float64(gs.Zbounds[e.Z-1].Min.Y), gs.Zbounds[e.Z-1])
				phps2 := NewPhReaderParticleSystem(float64(e.Sprite.X)+float64(e.Sprite.GetSpriteRect().Dx()/2), float64(gs.Zbounds[e.Z-1].Min.Y), gs.Zbounds[e.Z-1])
				RegisterEntity(&Entity{ParticleSystem: phps, LifeTime: 20, Z: e.Z, Sprite: phps.Sprite})
				RegisterEntity(&Entity{ParticleSystem: phps2, LifeTime: 20, Z: e.Z, Sprite: phps2.Sprite})
				uiDat.Timers["doAtTime"].TurnOn()
			}
		}
	}

	if uiDat.state == Selected && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if e.Z < 10 {
			UpdateEntityZAndReSortEntitySlice(e.Id, 10)
			e.Sprite.DOptsUpdaterTag = ""

		}

	}

	if uiDat.state == JustFocused {
		e.EventHub.Publish(events.UISpriteAction{e.UiData.Label, "clicked"})
		UpdateEntityZAndReSortEntitySlice(e.Id, 13)
		boxSprite := &sprite.Sprite{Img: uiDat.MainImg, X: uiDat.baseX, Y: uiDat.baseY}
		boxEnt := &Entity{Sprite: boxSprite}
		boxEnt.Z = 0
		RegisterEntity(boxEnt)
		e.TempLinkedID = boxEnt.Id
		uiDat.state = Selected
		AltImageWhenClickedUpdater(e, gs)
	}

	if us.UpdateFunc == nil && uiDat.state == Animation {
		uiDat.state = ExtraSpriteAnimationCompleted
	}

	if uiDat.state == ExtraSpriteAnimationCompleted {
		uiDat.state = FinishedButStayOpen

	}
}

func PhReaderDoAtTime(ent *Entity, gs GameState) {
	ent.Sprite.DOptsUpdaterTag = ""

	RemoveEntity(ent.TempLinkedID)
	us := ent.Sprite
	uiDat := ent.UiData

	UpdateEntityZAndReSortEntitySlice(ent.Id, 13)

	us.Shader = registry.ShaderMap["PH"]
	us.XYUpdater = nil
	us.ShaderParams["PHValue"] = uiDat.Environment.NaturalPHLevel
	us.ShaderParams["Point"] = []float64{3, 10}
	us.ShaderParams["Radius"] = 3.0
	us.UpdateFunc = MoveSpriteToDestinationAndSpin

	sp := &sprite.Sprite{Img: uiDat.AltImg, X: uiDat.baseX, Y: uiDat.baseY}
	sp.UpdateFunc = MoveSpriteToDestination
	sp.DOptsUpdaterParams = make(map[string]float64)
	sp.DOptsUpdaterParams["speed"] = 8.0
	sp.DOptsUpdaterParams["destinationX"] = 420
	sp.DOptsUpdaterParams["destinationY"] = float64(registry.Config.ScreenHeight / 10)
	us.LinkedSprite = sp

	ent.UiData.Timers[PublishAtTime] = util.NewTimer(0.5)
	ent.UiData.Timers[PublishAtTime].TimerUpdater = PublishAtTimerUpdater
	ent.UiData.Timers[PublishAtTime].TurnOn()
	ent.Parameters[PublishAtTime] = events.UISpriteAction{UiSprite: "phreader", UiSpriteAction: "ph reading"}

	//us.UpdateShaderParams = shaders.UpdateCounter
}
