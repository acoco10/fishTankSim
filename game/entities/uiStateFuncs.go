package entities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/stringConstants"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
)

const (
	DoAtHovered = "hovered"
	//timers
	DoAtTime      = "doAtTime"
	PublishAtTime = "publishAtTime"
)

func PublishGraphicHovered(ent *Entity, gs GameState) {
	if len(ent.PublishedGraphicIDs) == 0 && ent.UiData.variables["amountAvailable"] == 0 {
		var msg string
		var y float64
		if ent.UiData.Label == string(PiggyBank) {
			msg = fmt.Sprintf("$%0.2f", gs.Player.Money)
			y = float64(ent.UiData.Y)
		}
		if ent.UiData.Label == string(FishFood) {
			msg = fmt.Sprintf("%d/100", gs.Player.Food)
			y = float64(ent.UiData.Y) + float64(ent.Sprite.GetSpriteRect().Dy()) + 20
		}
		cashGraphicID := graphics.NewFadeInTextGraphic(
			msg,
			float64(ent.UiData.X)+float64(ent.Sprite.GetSpriteRect().Dx()/2),
			y,
			120)
		ent.PublishedGraphicIDs = append(ent.PublishedGraphicIDs, cashGraphicID)
	}
}

func NotClickable(ent *Entity, gs GameState) {
	us := ent.Sprite
	if us.SpriteHovered() {
		us.Shader = registry.ShaderMap[registry.Lowlight]
		if ent.DoAt[DoAtHovered] != nil {
			ent.DoAt[DoAtHovered](ent, gs)
		}
	} else {
		us.Shader = nil
	}

}

func UISpriteIdleUpdater(ent *Entity, gs GameState) {
	if gs.MouseFlags.WindowOpen {
		return
	}
	us := ent.Sprite
	uiDat := ent.UiData

	if ent.UiData.Flags["used"] {
		return
	}

	if ent.Sprite.SpriteHovered() {
		if ent.DoAt[DoAtHovered] != nil {
			ent.DoAt[DoAtHovered](ent, gs)
		}
	}

	if gs.FocusedEntity == ent {
		ent.StateMachine.Transition(ent)
	}

	if gs.HoveredUiSprite == ent {
		if uiDat.state != Disabled {
			if us.Img == uiDat.MainImg {
				if us.X == uiDat.baseX {
					us.Y -= 5
					us.X += 5
					if !uiDat.Flags["keepShader"] {
						us.Shader = registry.ShaderMap["Highlight"]
					}
				}
			}

		}
	} else {
		if !uiDat.Flags["keepShader"] {
			if uiDat.Flags["lowLight"] && us.SpriteHovered() {
				return
			}
			us.Shader = nil
			us.X = uiDat.baseX
			us.Y = uiDat.baseY
		}
	}
}

func PositionUpdate(ent *Entity, gs GameState) {
	if UnFocusCheck() {
		ent.StateMachine.Transition(ent)
	}

	if ClickCheck() {
		UnFocus(ent.Id)
		ent.UiData.baseX = ent.Sprite.X
		ent.UiData.baseY = ent.Sprite.Y
	}

}

func UpdateSkimmer(ent *Entity, gs GameState) {
	if ClickCheck() && ent.CursorInActivationRect() {
		UpdateEntityZAndReSortEntitySlice(ent.Id, MidLayerZ)
		gs.CursorUpdater.ChangeSpeed(0.2)
		skimmerBounds := gs.Zbounds[0]
		skimmerBounds.Max.X += ent.Sprite.GetSpriteRect().Dx()
		skimmerBounds.Min.X -= ent.Sprite.GetSpriteRect().Dy()
		gs.CursorUpdater.SetBounds(gs.Zbounds[0])
	}
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		UpdateEntityZAndReSortEntitySlice(ent.Id, NotInTankZ)
		gs.CursorUpdater.ResetSpeed()
		gs.CursorUpdater.ResetBounds()
		if UnFocusCheck() {
			UiSpriteTurnOffEverything(ent)
			ent.StateMachine.Transition(ent)
		}
	}

}

func UpdateUseOnTank(ent *Entity, gs GameState) {
	if ClickCheck() && ent.CursorInActivationRect() {
		ent.StateMachine.Transition(ent)
	}
}

func UseOnTank(ent *Entity) {
	ent.UiData.state = Disabled
	UnFocus(ent.Id)
	ent.Sprite.XYUpdater = nil
	ent.UiData.Flags["used"] = true
	ent.UiData.Timers["waitAndReset"] = util.NewTimer(2.5)
	ent.UiData.Timers["waitAndReset"].TurnOn()
	ent.Sprite.DOptsUpdaterTag = "swirl"
	if ent.UiData.HoverImg != nil {
		ent.Sprite.Img = ent.UiData.HoverImg
	}
	if ent.UiData.Flags["resort"] {
		UpdateEntityZAndReSortEntitySlice(ent.Id, MidLayerZ)
		ent.UiData.Flags["resort"] = false
	}
	if ent.UiData.Flags["unDraw"] {
		ent.Draw = false
	}
	if !ent.UiData.Flags["particlesGenerated"] {
		if ent.UiData.Label == string(Fertilizer) {
			spriteBounds := ent.Sprite.GetSpriteRect()
			fps := NewFertilizerParticleSystem(float64(ent.Sprite.X+float32(spriteBounds.Dx()/4)), float64(ent.Sprite.Y+float32(spriteBounds.Dy()/2)), ent.UiData.ActivationRect)
			fpent := &Entity{ParticleSystem: fps, Sprite: fps.Sprite}
			fpent.Z = MidLayerZ
			fpent.LifeTime = 8.0
			println("registering fertilizer particle entity")
			RegisterEntity(fpent)
			ent.UiData.Flags["particlesGenerated"] = true
			return
		}

		if ent.UiData.Label == string(PhBoost) {
			AddPHEffectParticles(ent.UiData.ActivationRect, 1)
			ent.UiData.Flags["particlesGenerated"] = true

		}

		if ent.UiData.Label == string(PhReduce) {
			AddPHEffectParticles(ent.UiData.ActivationRect, 2)
			ent.UiData.Flags["particlesGenerated"] = true
		}
	}
}

func PublishPickedUpEventIfClicked(ent *Entity, gs GameState) {
	if ent.Sprite.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		PublishPickedUpEvent(ent)
	}

}

func AltImageWhenClickedUpdater(ent *Entity, gs GameState) {
	if ent.Sprite.Img != ent.UiData.HoverImg {
		img := ent.UiData.HoverImg
		sp := ent.Sprite
		sp.Img = img
		if ent.UiData.Flags["updater"] {
			//if you don't do this, it won't be centered on the sprite
			x, y := util.GetScaledCursorPosition()
			sp.X = float32(x - img.Bounds().Dx()/2)
			sp.Y = float32(y - img.Bounds().Dy()/2)
			sp.XYUpdater = sprite.NewUpdater(sp)
		}
		if ent.UiData.Flags["outline"] {
			ent.Sprite.AddColoredOutlineShader(colornames.Yellow)
		}
		if ent.effectHandler != nil {
			//clear exclamation or other graphics
			ent.effectHandler()
		}

		if ent.Z < 13 {
			UpdateEntityZAndReSortEntitySlice(ent.Id, 13)
		}

		if ent.UiData.Flags["autoTransition1"] {
			ent.StateMachine.Transition(ent)

		}
		if ent.UiData.Flags["clickTransition"] {
			//this whole flag is for journal and notebook, need a more generalized way of dealing with linked entities that need specific behaviour, linked state updater?
			ent.Sprite.UpdateFunc = MoveSpriteToDestination
			spacing := 60
			w := registry.Config.ScreenWidth
			h := registry.Config.ScreenHeight
			x := w / 2
			y := h / 2
			ent.Sprite.DOptsUpdaterParams["destinationX"] = float64(x - ent.Sprite.GetSpriteRect().Dx()/2 + spacing)
			ent.Sprite.DOptsUpdaterParams["destinationY"] = float64(y - ent.Sprite.GetSpriteRect().Dy()/2)
			ent.Sprite.DOptsUpdaterParams["speed"] = 8
			ent.UiData.Timers["freeze"].TurnOn()

			ent2, exists := GetEntity(ent.LinkedID)
			if exists {
				ent2.Draw = true
				ent2.Z = 15
				ent2.UiData.Timers["freeze"].TurnOn()
				ent2.Sprite.UpdateFunc = MoveSpriteToDestination
				ent2.Sprite.DOptsUpdaterParams["destinationX"] = float64(x - ent2.Sprite.GetSpriteRect().Dx()/2 - ent.Sprite.GetSpriteRect().Dx() - spacing/2)
				ent2.Sprite.DOptsUpdaterParams["destinationY"] = float64(y - ent2.Sprite.GetSpriteRect().Dy()/2)
				ent2.Sprite.DOptsUpdaterParams["speed"] = 6
				ent2.StateMachine.Transition(ent2)
				ZSortEntities()
			}
			ent.StateMachine.Transition(ent)
		}
	}
}

func ActivationRectUpdater(ent *Entity, gs GameState) {
	if ent.Sprite.CheckMiddleOfSpriteInRect(ent.UiData.ActivationRect) {
		ent.Sprite.ChangeOutlineColor(colornames.Greenyellow)
		if ClickCheck() {
			if ent.UiData.Flags[stringConstants.Swirl] {
				ent.Sprite.DOptsUpdaterTag = stringConstants.Swirl
			}
			ent.Sprite.XYUpdater = nil
			ent.Z = MidLayerZ
			ZSortEntities()
			ent.UiData.Timers[DoAtTime].TurnOn()
			ent.StateMachine.Transition(ent)
		}
	} else {
		ent.Sprite.ChangeOutlineColor(colornames.Yellow)
	}

	if UnFocusCheck() {
		UiSpriteTurnOffEverything(ent)
		ent.StateMachine.Reset(ent)
	}
}

func UsedInActivationRect(ent *Entity, gs GameState) {
	if UnFocusCheck() {
		UiSpriteTurnOffEverything(ent)
		ent.StateMachine.Reset(ent)
	}
}

func AltImageHovered(ent *Entity, gs GameState) {
	if ent.UiData.Flags["freezeDone"] {
		if ent.Sprite.SpriteHovered() {
			ent.Sprite.Shader = registry.ShaderMap["Highlight"]
		} else {
			if !ent.UiData.Flags["keepShader"] {
				ent.Sprite.Shader = nil
			}
		}
		if ent.Sprite.SpriteHovered() && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			ent.UiData.Flags["freezeDone"] = false
			ent.StateMachine.Transition(ent)
			ent2, exists := GetEntity(ent.LinkedID)
			if exists {
				ent2.UiData.Flags["freezeDone"] = false
				UiSpriteTurnOffEverything(ent2)
				UnFocus(ent2.Id)
				ent2.StateMachine.Reset(ent2)
			}
		}
		if UnFocusCheck() {
			ent.UiData.Flags["freezeDone"] = false
			ent.StateMachine.Reset(ent)
			UiSpriteTurnOffEverything(ent)
			UnFocus(ent.Id)
			ent2, exists := GetEntity(ent.LinkedID)
			if exists {
				ent2.UiData.Flags["freezeDone"] = false
				UiSpriteTurnOffEverything(ent2)
				UnFocus(ent2.Id)
				ent2.StateMachine.Reset(ent2)
			}
		}
	}
}

func TurnOffEveryThingOnUnFocus(ent *Entity, gs GameState) {
	if UnFocusCheck() || ClickCheck() && !ent.Sprite.SpriteHovered() {
		ent.StateMachine.Reset(ent)
		UiSpriteTurnOffEverything(ent)
		UnFocus(ent.Id)
	}
}

func PublishPickedUpEvent(ent *Entity) {
	ent.Sprite.DOptsUpdaterTag = ""
	ev := events.UISpriteAction{
		UiSprite:       ent.UiData.Label,
		UiSpriteAction: "picked up",
	}
	ent.EventHub.Publish(ev)
	if ent.UiData.Flags["revert"] {
		UnFocus(ent.Id)
		UiSpriteTurnOffEverything(ent)
	}
}

func DisabledState(ent *Entity, gs GameState) {
	//state will only be transitioned out of by external factors?
	//subs or other entities
	//we could give it a flag to leave (thinking emoji)
	return
}
