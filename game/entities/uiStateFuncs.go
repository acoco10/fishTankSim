package entities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/input"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/stringConstants"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/acoco10/fishTankWebGame/gameDevMath"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"image"
	color2 "image/color"
)

const (
	DoAtHovered = "hovered"

	//timers
	DoAtTime             = "doAtTime"
	PublishAtTime        = "publishAtTime"
	Transition           = "transition"
	Reset                = "reset"
	GenFood              = "GenFood"
	Freeze               = "freeze"
	PhGuessAddMoneyDelay = "PhGuessAddMoneyDelay"
	PickedUp             = "pickedUp"
	OnEnterFreeze        = "onEnterFreeze"
	ClickedCoolDown      = "clickCoolDown"
	//tags
	uiTag = "ui"

	//"Do at tags"
	CounterFinished = "counterFinished"
	Used1           = "used1"
	Used2           = "used2"
	Used3           = "used3"
)

func (e *Entity) SetUIState(state UiSpriteState) {
	e.UiData.state = state
}

func (sm *StateMachine) AppendState(newUpdater EntityUpdater, newtransitionFunc EntityTransitioner) {
	if sm.States[len(sm.States)] != nil {
		sm.States[len(sm.States)].TransitionTo = len(sm.States) + 1
		newState := &StateHandler{Updater: newUpdater, TransitionOutFunc: newtransitionFunc, TransitionTo: 1}
		sm.States[len(sm.States)+1] = newState
	}
}

func InitStateMachine(initState EntityUpdater, updateFunc EntityUpdater, transitionFunc1 EntityTransitioner, transitionFunc2 EntityTransitioner) *StateMachine {
	States := make(map[int]*StateHandler)
	idleFunc := &StateHandler{Updater: initState, TransitionOutFunc: transitionFunc1, TransitionTo: 2}

	//set to default idle updater if no custom state is provided`
	if initState == nil {
		idleFunc = &StateHandler{Updater: UISpriteIdleUpdater, TransitionOutFunc: transitionFunc1, TransitionTo: 2}
	}

	if updateFunc != nil {
		pickedUp := &StateHandler{Updater: updateFunc, TransitionOutFunc: transitionFunc2, TransitionTo: 1}
		States[2] = pickedUp
	}

	States[1] = idleFunc

	sm := &StateMachine{States: States, CurrentState: 1}
	sm.EveryUpdate = append(sm.EveryUpdate, MousedToBase, RemoveEffect)
	sm.EveryUpdateEarlyReturnConditions = append(sm.EveryUpdateEarlyReturnConditions, FreezeCheck, ZoomCheck)
	return sm
}

func FreezeCheck(ent *Entity, gs GameState) bool {
	return ent.HasFreeze()
}

func ZoomCheck(ent *Entity, gs GameState) bool {
	return gs.ZoomedFor != NotZoomed
}

func AddClickme(ent *Entity) {
	ent.Sprite.DOptsUpdaterParams["swirlSpeedX"] = -2
	ent.Sprite.DOptsUpdaterParams["swirlSpeedY"] = -1
	ent.Sprite.DOptsUpdaterParams["update"] = 0.1

	ent.Sprite.Shader = registry.ShaderMap[registry.Highlight]
	ent.Sprite.DOptsUpdaterTag = stringConstants.Swirl
}

func RemoveClickMe(ent *Entity) {
	ent.ClearKeepShader()
	ent.Sprite.Shader = nil
	ent.Sprite.DOptsUpdaterTag = ""
}

func AltImageFade(ent *Entity) {
	target := 0.0
	speed := 0.1
	minimum := 0.1

	sprite.SetOpacityUpdater(ent.Sprite, target, speed, minimum)

	filtered := ent.PublishedGraphicIDs[:0]
	for _, id := range ent.PublishedGraphicIDs {
		graphic, exists := graphics.GetGraphic(id)
		if !exists {
			continue
		}
		text, ok := graphic.(*graphics.FadeInText)
		if !ok {
			continue
		}
		//the minimum is set to 0.0 becuase text is more attention grabbing then a faded sprite
		text.SetFadeOut(float32(target), float32(speed), 0.0)
		filtered = append(ent.PublishedGraphicIDs, id)
	}
	ent.PublishedGraphicIDs = filtered
}

func PublishGraphicHovered(ent *Entity, gs GameState) {
	if !ent.Draw {
		return
	}

	switch Label(ent.UiData.Label) {
	case PiggyBank:
		if len(ent.PublishedGraphicIDs) == 0 && ent.UiData.variables["amountAvailable"] == 0 {
			var msg string
			var y float64
			if ent.UiData.Label == string(PiggyBank) {
				msg = fmt.Sprintf("$%0.2f", gs.Player.Money)
				y = float64(ent.UiData.Y)
			}

			cashGraphicID := graphics.NewFadeInTextGraphic(
				msg,
				float64(ent.UiData.X)+float64(ent.Sprite.GetSpriteRect().Dx()/2),
				y,
				120)
			ent.PublishedGraphicIDs = append(ent.PublishedGraphicIDs, cashGraphicID)
		}
	case FishFood:
		msg := fmt.Sprintf("%d/100", gs.Player.Food)
		y := float64(ent.UiData.Y) + float64(ent.Sprite.GetSpriteRect().Dy()) + 20
		foodGraphicID := graphics.NewFadeInTextGraphic(
			msg,
			float64(ent.UiData.X)+float64(ent.Sprite.GetSpriteRect().Dx()/2),
			y,
			120)
		ent.PublishedGraphicIDs = append(ent.PublishedGraphicIDs, foodGraphicID)
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

	if ent.HasUsed() {
		return
	}

	if ent.HasFreeze() {
		return
	}

	if ent.HasIdleClickEffect() {
		if ent.Draw {
			ent.ClearIdleClickEffect()
			ent.AddClickEffect("")
		}
	}

	if ent.Sprite.SpriteHovered() {
		if ent.DoAt[DoAtHovered] != nil {
			ent.DoAt[DoAtHovered](ent, gs)
		}
	}

	if gs.HoveredUiSprite == ent {
		if uiDat.state != Disabled {
			if us.Img == uiDat.MainImg {
				if !ent.HasNoOffset() {
					if us.X == uiDat.baseX {
						us.Y -= 5
						us.X += 5
					}
					if !ent.HasKeepShader() {
						us.Shader = registry.ShaderMap["Highlight"]
					}
				}
			}
			if registry.ClickCheck() {
				gs.FocusedEntity = ent
				if ent.UiData.Timers[PickedUp] != nil {
					ent.UiData.Timers[PickedUp].TurnOn()
				}
				ent.StateMachine.Transition(ent)
			}
		}
	} else {
		if ent.HasLowLight() && us.SpriteHovered() {
			us.Shader = registry.ShaderMap[registry.Lowlight]
			return
		}

		if !ent.HasKeepShader() {
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

	if registry.ClickCheck() {
		UnFocus(ent.Id)
		ent.UiData.baseX = ent.Sprite.X
		ent.UiData.baseY = ent.Sprite.Y
	}

}

func transitionToUpdateSkimmer(ent *Entity) {
	for _, p := range ParticleEntList {
		p.AddColoredOutlineShader(colornames.Red)
	}

	ent.AddClickEffect("")
}

func UpdateSkimmer(ent *Entity, gs GameState) {

	KeepDebrisContained(ent)

	if ent.HasFreeze() {
		return
	}

	if ent.UiData.HoverImg != nil {
		ent.Sprite.AddColoredOutlineShader(colornames.Yellow)
		ent.Sprite.Img = ent.UiData.HoverImg
		ent.Z = NotInTankZ
		ZSortEntities()
	}

	if registry.ClickCheck() {
		if ent.HasFishCaught() {
			ThrowAwayDebrisUpdater(ent, gs)
			return
		}
		distance := UiSpriteDistanceFromBase(ent)

		if distance < 100 {
			ent.StateMachine.States[ent.StateMachine.CurrentState].TransitionTo = 5
			ent.StateMachine.States[ent.StateMachine.CurrentState].TransitionOutFunc = nil
			ent.StateMachine.Transition(ent)
		}
	}

	if ent.CursorInActivationRect() {
		ent.Sprite.ChangeOutlineColor(colornames.Greenyellow)
		ent.Z = NotInTankZ
	}

	if registry.ClickCheck() && ent.CursorInActivationRect() {
		ent.Parameters.Ints[Power] = 0 // Update while pressing
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && ent.CursorInActivationRect() {
		if ent.HasFishCaught() {
			return
		}
		gs.CursorUpdater.Lock()
		ent.Sprite.XYUpdater = nil

		if ent.Parameters.Ints[Power] == 0 {
			deInit, speedUpdater := DrawCursorPowerEffect(EffectParams{})
			ent.AddDeInitHandler(deInit)
			ent.DeInitEffects()
			ent.effectUpdateHandler = speedUpdater
		}

		ent.Parameters.Ints[Power] = registry.ClickDuration()
		if ent.Parameters.Ints[Power]%20 == 0 {
			power := float64(ent.Parameters.Ints[Power] / 20)
			scalar := power * power
			ent.effectUpdateHandler(0.05 * scalar)
		} // Update while pressing
	}

	if registry.ReleaseCheck() && ent.CursorInActivationRect() {
		if ent.HasFishCaught() {
			ent.ClearFishCaught()
			return
		}

		ent.DeInitEffects()
		ent.StateMachine.Transition(ent)
	}

	if UnFocusCheck() {
		UiSpriteTurnOffEverything(ent)
		ent.StateMachine.Reset(ent)
	}

}

func AddMoveFuncToSprite(ent *Entity) {
	powScaled := min(float64(ent.Parameters.Ints[Power])/10.0, 20.0)
	powScaled = max(5.0, powScaled)
	ent.Sprite.XYUpdater = nil
	UpdateEntityZAndReSortEntitySlice(ent.Id, MidLayerZ)

	startX := float64(ent.Sprite.X)
	endX := float64(ent.Parameters.Rectangles[TankBounds].Min.X) - 200 + ((float64(powScaled) - 2) * 50)
	endX = min(float64(ent.Sprite.X)-75, endX) // clamp so it still goes slightly left
	startY := float64(ent.Sprite.Y)
	endY := float64(ent.Sprite.Y)
	bottomOfTank := float64(ent.Parameters.Rectangles[TankBounds].Max.Y) - startY - float64(ent.Sprite.GetSpriteRect().Dy())
	peakOffSet := powScaled / 10 * bottomOfTank

	ent.Sprite.MoveParams.MathFunction =
		gameDevMath.SetDownwardParabolaPath(
			startX,
			startY,
			endX,
			endY,
			peakOffSet,
		)

	distance := startX - endX

	xChangeClamped := -1 * distance / 50
	ent.Sprite.MoveParams.XChange = xChangeClamped
	ent.Sprite.MoveParams.FinalX = float32(endX)
	ent.Sprite.MoveParams.FinalY = float32(10)
	ent.Sprite.MoveParams.StartingX = ent.Sprite.X
	ent.Sprite.UpdateFunc = sprite.MathFunctionSpriteUpdater
}

func CatchDebrisDuringMove(ent *Entity, gs GameState) {

	skimmerNet := ent.Parameters.Rectangles[NetOpening]
	skimmerNet = util.MoveRectTo(skimmerNet, ent.Sprite.X+35, float32(ent.Sprite.GetSpriteRect().Max.Y-skimmerNet.Dy())-20)

	if ent.Sprite.MoveParams.DirX < 0 && ent.Sprite.MoveParams.DirY < 0 {
		for _, p := range ParticleEntList {

			if p.GetSpriteRect().Overlaps(skimmerNet) {
				ent.EventHub.Publish(DebrisCaptured{})
				p.bounds = &skimmerNet
			}
		}

		for _, creature := range CreatureList {
			if creature.Sprite.GetSpriteRect().Overlaps(skimmerNet) {
				creature.CreatureData.TankBoundaries = skimmerNet
				creature.Transition()
				ent.SetFishCaught()
			}
		}
	}

	if ent.Sprite.X <= float32(ent.Parameters.Rectangles[TankBounds].Min.X) {
		ent.Sprite.UpdateFunc = MoveSpriteToDestination
		ent.Sprite.DOptsUpdaterTag = ""
		ent.Sprite.DOptsUpdaterParams["destinationX"] = float64(ent.Parameters.Rectangles[TankBounds].Min.X + 1)
		ent.Sprite.DOptsUpdaterParams["destinationY"] = float64(ent.Sprite.MoveParams.FinalY)
		ent.Sprite.DOptsUpdaterParams["speed"] = 4.0
	}

	if ent.Sprite.UpdateFunc == nil {
		ent.StateMachine.Transition(ent)
	}
	KeepDebrisContained(ent)

}

func ResetPosition(ent *Entity) {
	x, y := util.GetScaledCursorPosition()
	ent.Sprite.UpdateFunc = MoveSpriteToDestination
	ent.Sprite.DOptsUpdaterTag = "cursor"
	ent.Sprite.DOptsUpdaterParams["destinationX"] = float64(x)
	ent.Sprite.DOptsUpdaterParams["destinationY"] = float64(y)
	ent.Sprite.DOptsUpdaterParams["speed"] = 7.0
	KeepDebrisContained(ent)
}

func ResetUpdaterAfterSpriteUpdate(ent *Entity, gs GameState) {
	KeepDebrisContained(ent)

	if ent.UiData.Sprite.XYUpdater == nil {
		if ent.Sprite.UpdateFunc == nil {
			gs.CursorUpdater.UnLock()
			AddUiSpriteXYUpdater(ent)
			ent.StateMachine.Transition(ent)
		}
	}
}

func UiSpriteDistanceFromBase(ent *Entity) float64 {
	return DistanceFunc(ent.Sprite.X, ent.UiData.baseX, ent.Sprite.Y, ent.UiData.baseY)
}

func ThrowAwayDebrisUpdater(ent *Entity, gs GameState) {
	if !ent.HasFishCaught() {
		ent.UiData.Timers[Reset].TurnOn()
	}

	if ent.Sprite.UpdateFunc == nil && ent.UiData.Sprite.XYUpdater != nil {
		ent.Sprite.Img = ent.UiData.AltImg
		for _, p := range ParticleEntList {
			p.Sprite.Shader = nil
			if p.bounds != nil {

				p.UnderWater = false
				p.bounds = nil
				p.baseVy = 9.8
				if !ent.HasFishCaught() {
					p.floorLevel = registry.Config.ScreenHeight
				}
			}
		}
		for _, creature := range CreatureList {
			if creature.CreatureData.State == Captured {
				creature.CreatureData.TankBoundaries = gs.Zbounds[12]
				creature.Transition()
			}
		}
		ent.updateEntities = make(map[uint32]struct{})
	}

	if ent.HasFishCaught() {
		ent.ClearFishCaught()
		ent.SetFreeze()
		ent.UiData.Timers[Freeze].TurnOn()
		return
	}

	KeepDebrisContained(ent)
}

func KeepDebrisContained(ent *Entity) {
	skimmerNet := ent.Parameters.Rectangles[Net]
	skimmerNet = util.MoveRectTo(skimmerNet, ent.Sprite.X+35, float32(ent.Sprite.GetSpriteRect().Max.Y-skimmerNet.Dy())-10)

	for _, p := range ParticleEntList {
		if p.bounds != nil {
			p.bounds = &skimmerNet
		}
	}
	for _, creature := range CreatureList {
		if creature.CreatureData.State == Captured {
			creature.CreatureData.TankBoundaries = skimmerNet
		}
	}
}

func UpdateUseOnTank(ent *Entity, gs GameState) {
	if registry.ClickCheck() && ent.CursorInActivationRect() {
		ent.StateMachine.Transition(ent)
	}
}

func UseOnTank(ent *Entity) {
	ent.UiData.state = Disabled
	UnFocus(ent.Id)
	ent.Sprite.XYUpdater = nil
	ent.SetUsed()
	ent.UiData.Timers["waitAndReset"] = util.NewTimer(2.5)
	ent.UiData.Timers["waitAndReset"].TurnOn()
	ent.Sprite.DOptsUpdaterTag = stringConstants.Swirl
	if ent.UiData.HoverImg != nil {
		ent.Sprite.Img = ent.UiData.HoverImg
	}

	UpdateEntityZAndReSortEntitySlice(ent.Id, MidLayerZ)

	if ent.HasUnDraw() {
		ent.Draw = false
	}
	if !ent.HasParticlesGenerated() {
		if ent.UiData.Label == string(Fertilizer) {
			spriteBounds := ent.Sprite.GetSpriteRect()
			fps := NewFertilizerParticleSystem(float64(ent.Sprite.X+float32(spriteBounds.Dx()/4)), float64(ent.Sprite.Y+float32(spriteBounds.Dy()/2)), ent.UiData.ActivationRect)
			fpent := &Entity{ParticleSystem: fps, Sprite: fps.Sprite}
			fpent.Z = MidLayerZ
			fpent.LifeTime = 8.0
			println("registering fertilizer particle entity")
			RegisterEntity(fpent)
			ent.SetParticlesGenerated()
			return
		}

		if ent.UiData.Label == string(PhBoost) {
			AddPHEffectParticles(ent.UiData.ActivationRect, 1)
			ent.SetParticlesGenerated()
		}

		if ent.UiData.Label == string(PhReduce) {
			AddPHEffectParticles(ent.UiData.ActivationRect, 2)
			ent.SetParticlesGenerated()
		}
	}
}

func (ent *Entity) CheckClickedInRightPlace() bool {
	if ent.HasNoActivationRect() {
		return ent.Sprite.SpriteHovered() && registry.ClickCheck()
	}

	return ent.Sprite.CheckMiddleOfSpriteInRect(ent.UiData.ActivationRect) && registry.ClickCheck()
}

func RemoveEffect(ent *Entity, gs GameState) {
	if ent.CheckClickedInRightPlace() {
		ent.DeInitEffects()
	}
}

func ClickedUpdater(ent *Entity, gs GameState) {

	if !ent.HasKeepShader() {
		ent.Sprite.DOptsUpdaterTag = ""
		ent.Sprite.Shader = nil
	}

	ent.DeInitEffects()

	if ent.HasCounter() {
		if ent.Parameters.Ints[IndexCounter] == ent.Parameters.Ints[IndexCounterMax]-1 {
			if ent.DoAt[Used1] != nil {
				ent.DoAt[Used1](ent, gs)
			}
		}
		if ent.Parameters.Ints[IndexCounter] == ent.Parameters.Ints[IndexCounterMax]-2 {
			if ent.DoAt[Used2] != nil {
				ent.DoAt[Used2](ent, gs)
			}
		}
		if ent.Parameters.Ints[IndexCounter] == ent.Parameters.Ints[IndexCounterMax]-3 {
			if ent.DoAt[Used3] != nil {
				ent.DoAt[Used3](ent, gs)
			}
		}

		ent.Parameters.Ints[IndexCounter]--
		if !ent.HasFreeze() {
			ent.SetAddIdleClickEffect()
		}

		if ent.Parameters.Ints[IndexCounter] <= 0 {
			ent.DoAt[CounterFinished](ent, gs)
			ent.Publish(events.UiSpriteCounterFinished{UISprite: ent.UiData.Label})
		}
	}

	if ent.HasLikeWindow() {
		ent.EventHub.Publish(events.WindowOpened{Window: ent.UiData.Label + uiTag})
		if UnFocusCheck() {
			turnOffEverythingUnFocusResetBundle(ent)
			ent.EventHub.Publish(events.WindowClosed{Window: ent.UiData.Label + uiTag})
		}
	}

	sp := ent.Sprite

	if ent.HasAltImage() {
		if ent.Sprite.Img != ent.UiData.HoverImg {
			img := ent.UiData.HoverImg
			sp.Img = img
		}
	}

	if ent.HasUpdater() {
		img := ent.UiData.HoverImg
		//if you don't do this, it won't be centered on the sprite
		x, y := util.GetScaledCursorPosition()
		sp.X = float32(x - img.Bounds().Dx()/2)
		sp.Y = float32(y - img.Bounds().Dy()/2)
		sp.XYUpdater = sprite.NewUpdater(sp)
	}

	if ent.HasAddClickEffect() {
		ent.AddClickEffect("")
	}

	if ent.HasOutline() {
		ent.Sprite.AddColoredOutlineShader(colornames.Yellow)
	}

	if ent.Z < 13 {
		UpdateEntityZAndReSortEntitySlice(ent.Id, 13)
	}

	if ent.HasAutoTransition1() {
		ent.StateMachine.Transition(ent)
	}

	if ent.HasClickTransition() {
		//this whole flag is for journal and notebook
		//basically a menu with journalmainimg as single entrypoint for both
		TextMenuClickedUpdater(ent)
	}

}

func (ent *Entity) AddClickEffect(text string) {
	x, y := util.MiddleOfRect(ent.UiData.ActivationRect)

	if ent.HasNoActivationRect() {
		x, y = ent.Sprite.MidPointF()
	}

	deinit1 := DrawControlEffect(x, y, EffectParams{}, ClickHere, text, nil)
	ent.AddDeInitHandler(deinit1)
	if ent.HasNoActivationRect() {
		return
	}

	color := color2.RGBA{R: 100, G: 232, B: 199, A: 1}
	deinit2 := DrawRectEntity(ent.UiData.ActivationRect, color, true)

	ent.AddDeInitHandler(deinit2)
}

func TextMenuClickedUpdater(ent *Entity) {
	ent.Sprite.UpdateFunc = MoveSpriteToDestination
	spacing := 60
	w := registry.Config.ScreenWidth
	h := registry.Config.ScreenHeight
	x := w / 2
	y := h / 2
	ent.Sprite.DOptsUpdaterParams["destinationX"] = float64(x - ent.Sprite.GetSpriteRect().Dx()/2 + spacing)
	ent.Sprite.DOptsUpdaterParams["destinationY"] = float64(y - ent.Sprite.GetSpriteRect().Dy()/2)
	ent.Sprite.DOptsUpdaterParams["speed"] = 8
	ent.SetFreeze()
	ent.UiData.Timers[Freeze].TurnOn()

	ent2, exists := GetEntity(ent.LinkedID)
	if exists {
		if ent2.HasKeepShader() {
			ent.ClearKeepShader()
			ent.Sprite.Shader = nil
		}
		if !ent2.Draw {
			ent.Sprite.DOptsUpdaterParams["destinationX"] -= float64(spacing)
		} else {
			img := ent2.UiData.HoverImg
			ent2.Sprite.Img = img
			ent2.Z = 15
			ent2.UiData.Timers[Freeze].TurnOn()
			ent2.Sprite.UpdateFunc = MoveSpriteToDestination
			ent2.Sprite.DOptsUpdaterParams["destinationX"] = float64(x - ent2.Sprite.GetSpriteRect().Dx()/2 - ent.Sprite.GetSpriteRect().Dx() - spacing/2)
			ent2.Sprite.DOptsUpdaterParams["destinationY"] = float64(y - ent2.Sprite.GetSpriteRect().Dy()/2)
			ent2.Sprite.DOptsUpdaterParams["speed"] = 6
			ent2.SetFreeze()
			ent2.StateMachine.Transition(ent2)
			ZSortEntities()
		}
	}

	ent.StateMachine.Transition(ent)
}

func ActivationRectUpdaterFishFood(ent *Entity, gs GameState) {
	if ent.Sprite.CheckMiddleOfSpriteInRect(ent.UiData.ActivationRect) {
		ent.Sprite.ChangeOutlineColor(colornames.Greenyellow)
		ent.Z = MidLayerZ
		ZSortEntities()
		if registry.ClickCheck() {
			x, y := util.GetScaledCursorPosition()
			pt := image.Point{x, y}

			leftRect := ent.UiData.ActivationRect
			leftRect.Max.X = leftRect.Max.X / 2
			ent.Sprite.Img = ent.UiData.AltImg

			ev2 := input.MouseButtonPressedUISpriteActivity{
				Point: &util.Point{X: float32(x), Y: float32(y), PType: util.Food},
			}
			if pt.In(leftRect) {
				ent.Sprite.Img = ent.UiData.HoverImg
				ev2.Point.X += float32(ent.Sprite.Img.Bounds().Dx())
				ev2.Point.Tag = "left"
			}
			ent.Parameters.Events[PointEvent] = ev2
			ent.Parameters.Ints[IndexCounter] = 5
			ent.UiData.Timers[GenFood].TurnOn()
		}

	} else {
		ent.Sprite.ChangeOutlineColor(colornames.Yellow)
	}

	if UnFocusCheck() {
		turnOffEverythingUnFocusResetBundle(ent)
	}
}

func MousedToBase(ent *Entity, gs GameState) {
	if ent.UiData.Timers[PickedUp] != nil {
		if ent.UiData.Timers[PickedUp].On == false && ent.StateMachine.CurrentState != 1 {
			if UiSpriteDistanceFromBase(ent) < 100 && registry.ClickCheck() {
				ent.Sprite.XYUpdater = nil
				turnOffEverythingUnFocusResetBundle(ent)
				ent.UiData.Timers[Freeze].TurnOn()
				ent.SetFreeze()
			}
		}
	}
}

func ActivationRectUpdaterPhReader(ent *Entity, gs GameState) {
	if ent.Sprite.CheckMiddleOfSpriteInRect(ent.UiData.ActivationRect) {
		ent.Sprite.ChangeOutlineColor(colornames.Greenyellow)
		if registry.ClickCheck() {
			DrawCircularProgressEffect(ent.Sprite.X, ent.Sprite.Y+float32(ent.Sprite.MidY()), EffectParams{Speed: 3, Cycles: 1})
			gs.CursorUpdater.Lock()
			ent.Sprite.DOptsUpdaterTag = stringConstants.Swirl
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
		turnOffEverythingUnFocusResetBundle(ent)
	}
}

func UsedInActivationRect(ent *Entity, gs GameState) {
	if ent.Sprite.DOptsUpdaterTag == "" {
		gs.CursorUpdater.UnLock()
	}

	if UnFocusCheck() {
		turnOffEverythingUnFocusResetBundle(ent)
	}
}

func AltImageHovered(ent *Entity, gs GameState) {
	if !ent.HasFreeze() {
		if !ent.HasKeepShader() {
			if ent.Sprite.SpriteHovered() {
				ent.Sprite.Shader = registry.ShaderMap["Highlight"]
			} else {
				ent.Sprite.Shader = nil
			}
		}

		if ent.Sprite.SpriteHovered() && registry.ClickCheck() {
			RemoveClickMe(ent)
			ent.StateMachine.Transition(ent)
			ent2, exists := GetEntity(ent.LinkedID)
			if exists {
				if ent.UiData.Label == string(Magazine) {
					ent2.Parameters.Events[EventAtTime] = events.UISpriteAction{UiSprite: string(GrandpasJournal), UiSpriteAction: "picked up"}
					ent2.UiData.Timers[PublishAtTime].TurnOn()
					ent2.UiData.Timers[Reset].TurnOn()
					ent2.Sprite.UpdateFunc = MoveSpriteToDestination
					ent2.Sprite.DOptsUpdaterParams["destinationX"] = float64(registry.Config.ScreenWidth)
					ent2.Sprite.DOptsUpdaterParams["destinationY"] = float64(registry.Config.ScreenHeight)
					ent2.Sprite.DOptsUpdaterParams["speed"] = 5
				} else {
					RemoveClickMe(ent2)
					ent2.ClearFreeze()
					UiSpriteTurnOffEverything(ent2)
					UnFocus(ent2.Id)
					ent2.StateMachine.Reset(ent2)
				}
			}

		}

		if UnFocusCheck() {
			turnOffEverythingUnFocusResetBundle(ent)
			ent2, exists := GetEntity(ent.LinkedID)
			if exists {
				ent2.ClearFreeze()
			}
		}
	}
}

func turnOffEverythingUnFocusResetBundle(ent *Entity) {
	ent.StateMachine.Reset(ent)
	UiSpriteTurnOffEverything(ent)
	UnFocus(ent.Id)
	ent2, exists := GetEntity(ent.LinkedID)
	if exists {
		UiSpriteTurnOffEverything(ent2)
		UnFocus(ent2.Id)
		ent2.StateMachine.Reset(ent2)
	}
}

func FadeOutOnNotHovered(ent *Entity, gs GameState) {
	if ent.Sprite.SpriteHovered() {
		ent.UiData.Timers["UnFocus"].Reset()
	}
	if !ent.Sprite.SpriteHovered() || UnFocusCheck() {
		ent.UiData.Timers["UnFocus"].TurnOn()
		if ent.Sprite.UpdateFunc == nil && ent.UiData.Timers["UnFocus"].TimerState == util.Done {
			AltImageFade(ent)
		}
		ent.StateMachine.Transition(ent)
	}

	if UnFocusCheck() {
		UiSpriteTurnOffEverything(ent)
		ent.StateMachine.Reset(ent)
	}

}

func TurnOffEveryThingOnSpriteAnimationComplete(ent *Entity, gs GameState) {
	if ent.Sprite.UpdateFunc == nil {
		ent.StateMachine.Reset(ent)
		UiSpriteTurnOffEverything(ent)
		UnFocus(ent.Id)
		ent.StateMachine.Reset(ent)
	}
}

func PublishPickedUpEvent(ent *Entity) {

	if ent.HasLikeWindow() {
		ent.EventHub.Publish(events.WindowClosed{Window: ent.UiData.Label + uiTag})
	}

	action := "picked up"

	if len(ent.Parameters.StringLists[Events]) > 0 {

		action += ent.Parameters.StringLists[Events][ent.Parameters.StringListsIndex[Events]]
		if ent.Parameters.StringListsIndex[Events] < len(ent.Parameters.StringLists[Events])-1 {
			ent.Parameters.StringListsIndex[Events]++
		} else {
			ent.Parameters.StringListsIndex[Events] = 0
		}

	}

	ev := events.UISpriteAction{
		UiSprite:       ent.UiData.Label,
		UiSpriteAction: action,
	}

	if ent.HasDeInitAfterUSed() {
		RemoveEntity(ent.Id)
	}

	ent.EventHub.Publish(ev)
	if ent.HasRevert() {
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
