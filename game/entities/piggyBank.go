package entities

import (
	"github.com/acoco10/fishTankWebGame/game/entImportableLoaders"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

const (
	AmountAvailable = "amountAvailable"
	CashUpdate      = "CashUpdate"
)

func InitPiggyBankStateMachine() *StateMachine {
	States := make(map[int]*StateHandler)
	disabled := &StateHandler{Updater: NotClickable, TransitionOutFunc: nil, TransitionTo: 2}
	idle := &StateHandler{Updater: UISpriteIdleUpdater, TransitionOutFunc: PiggyBankTransitionFunc, TransitionTo: 3}
	pickedUp := &StateHandler{Updater: AnimationCompletedMonitor, TransitionTo: 1}
	States[1] = disabled
	States[2] = idle
	States[3] = pickedUp
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
		money := e.UiData.variables[AmountAvailable]
		e.DeInitEffects()
		if e.AnimationCycles >= 1 {
			e.Sprite.GetAnimation().Reset()
			e.Sprite.CurrentAnimation = ""
			gm := PiggyBankMoneyAddedGraphic(money, e.Sprite.X, e.Sprite.Y, e, gs)
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

func PiggyBankMoneyAddedGraphic(amtAdded float64, basex float32, basey float32, ent *Entity, gs GameState) *graphics.GraphicManager {
	gm := &graphics.GraphicManager{}

	gm.UpdateableFloat = gs.Player.Money
	gm.UpdateableFloatStop = gs.Player.Money + amtAdded

	gm.Timers = make(map[string]*util.Timer)
	gm.Timers[graphics.Trigger] = util.NewTimer(0.2)
	gm.Timers[graphics.Trigger].TimerUpdater = graphics.TriggerTimerUpdater

	gm.Timers[CashUpdate] = util.NewTimer(0.016 * 10)
	gm.Timers[CashUpdate].TimerUpdater = graphics.CashUpdaterTimerUpdater

	gm.Params = make(map[string]any)
	gm.Params[graphics.Trig2] = CashUpdate

	for i := 0; i < int(amtAdded/.25); i++ {
		spEff := entImportableLoaders.LoadStaticEffect("25c", basex, basey, "")
		params := make(map[string]any)
		params["opacity"] = float32(1.0)
		gs := graphics.SpriteGraphic{Sprite: *spEff, UpdateFunc: MoveSpriteToDestinationUp, Parameters: params}
		gs.SetDrawFunc(graphics.FadeIn)
		gm.GraphicsToBePublished = append(gm.GraphicsToBePublished, &gs)
	}

	cs := ebiten.ColorScale{}
	cs.SetR(0.9)
	cs.SetB(0.9)
	cs.SetG(0.9)
	cs.SetA(1.0)

	x := float64(ent.UiData.X) + float64(ent.Sprite.GetSpriteRect().Dx()/2)
	y := float64(ent.UiData.Y)
	gm.TextPosition = []float64{x, y}

	if gm.Timers != nil {
		gm.Timers[graphics.Trigger].TurnOn()
	}

	gm.FinishedFunc = func(in any) {
		inEnt, ok := in.(*Entity)
		if !ok {
			log.Fatal("piggy bank finished func was passed non entity as paramater")
		}
		inEnt.UiData.variables[AmountAvailable] = 0
	}

	return gm
}
