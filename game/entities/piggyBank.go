package entities

import (
	"github.com/acoco10/fishTankWebGame/game/events"
)

func InitPiggyBankStateMachine() *StateMachine {
	States := make(map[int]*StateHandler)
	idle := &StateHandler{Updater: UISpriteIdleUpdater, TransitionFunc: PiggyBankTransitionFunc, TransitionTo: 2}
	pickedUp := &StateHandler{Updater: AnimationCompletedMonitor, TransitionTo: 1}
	States[1] = idle
	States[2] = pickedUp
	sm := &StateMachine{States: States, CurrentState: 1}
	return sm
}

func PiggyBankTransitionFunc(ent *Entity) {
	ent.Sprite.DOptsUpdaterTag = ""
	uiDat := ent.UiData
	uiDat.Publish(events.UISpriteAction{UiSprite: uiDat.Label, UiSpriteAction: "clicked"})
	ent.Sprite.SetAnimation("allowance")
	ent.UiData.state = Animation
}

func AnimationCompletedMonitor(e *Entity, gs GameState) {
	if e.Sprite.CurrentAnimation != "" {
		money := e.UiData.variables["amountAvailable"]
		if e.AnimationCycles >= 1 {
			e.Sprite.GetAnimation().Reset()
			e.Sprite.CurrentAnimation = ""
			gm := PiggyBankMoneyAddedGraphic(money, e.Sprite.X, e.Sprite.Y)
			e.GraphicManager = gm
			ev := events.MoneyAdded{Amount: money}
			e.EventHub.Publish(ev)
			e.AnimationCycles = 0
			UiSpriteTurnOffEverything(e)
			UnFocus(e.Id)
			e.UiData.state = Disabled
			e.StateMachine.Transition(e)

		}
	}
}
