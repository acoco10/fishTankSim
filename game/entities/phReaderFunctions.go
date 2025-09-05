package entities

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
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
				phps := NewPhReaderParticleSystem(float64(e.Sprite.X)+float64(e.Sprite.GetSpriteRect().Dx()/2), float64(e.Sprite.Y)+float64(e.Sprite.GetSpriteRect().Dy()), gs.Zbounds[e.Z-1])
				RegisterEntity(&Entity{ParticleSystem: phps, LifeTime: 20, Z: e.Z})
			}
		}
		ClickForTime(e, gs, phReaderDoAtTime)
	}

	if uiDat.state == Selected && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if e.Z < 10 {
			UpdateEntityZAndReSortEntitySlice(e.Id, 10)
		}
	}

	if uiDat.state == JustFocused {
		e.EventHub.Publish(events.UISpriteAction{e.UiData.Label, "clicked"})
		UpdateEntityZAndReSortEntitySlice(e.Id, 13)
		boxSprite := &sprite.Sprite{Img: uiDat.MainImg, X: uiDat.baseX, Y: uiDat.baseY}
		boxEnt := &Entity{Sprite: boxSprite}
		boxEnt.Z = 0
		RegisterEntity(boxEnt)
		e.LinkedID = boxEnt.Id
		uiDat.state = Selected
		AltImageWhenClickedUpdater(e, gs)
	}

	if us.UpdateFunc == nil && uiDat.state == Animation {
		uiDat.state = ExtraSpriteAnimationCompleted
	}

	if uiDat.state == ExtraSpriteAnimationCompleted {
		uiDat.state = FinishedButStayOpen
		ev := events.UISpriteAction{UiSprite: "phreader", UiSpriteAction: "ph reading"}
		uiDat.EventHub.Publish(ev)
	}
}

func phReaderDoAtTime(ent *Entity) {
	ent.Sprite.DOptsUpdaterTag = ""

	RemoveEntity(ent.LinkedID)
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
	us.LinkedSprite = sp

	//us.UpdateShaderParams = shaders.UpdateCounter
}
