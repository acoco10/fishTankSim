package entities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/entImportableLoaders"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/system"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"image"
	"math"
	"strconv"
	"strings"
)

type UiSpriteState uint8

const (
	Selected UiSpriteState = iota
	HoveredOver
	ClickedWhileBeingSelected
	Idle
	Clickable
	ExtraSpriteAnimationCompleted
	Activatable //meeting condition to be activated ie in "activation zone" on screen
	Animation
	JustFocused
	FinishedButStayOpen
	Disabled
)

type UiSpriteData struct {
	*sprite.Sprite
	baseX, baseY           float32
	MainImg                *ebiten.Image
	HoverImg               *ebiten.Image
	AltImg                 *ebiten.Image
	AltOffsetX, AltOffsetY float32
	*sprite.XYUpdater
	*tasks.EventHub
	state                     UiSpriteState
	stateWas                  UiSpriteState
	gameMode                  registry.GameMode
	clicked                   bool
	Draggable                 bool
	Label                     string
	highlight                 bool
	screenHeight, screenWidth int
	debugGraphicPublishedID   int
	extraSprite               *sprite.Sprite
	Environment               *system.Environment
	ActivationRect            image.Rectangle
	timers                    map[string]*util.Timer
	variables                 map[string]float64
	Flags                     map[string]bool
	stringVariables           map[string]string
}

func (e *Entity) UpdateUiSpriteEntity(gs *GameState) {
	e.UpdateUiSprite(gs)
}

func (e *Entity) SetUIState(state UiSpriteState) {
	e.UiData.state = state
}

func InitStateMachine(updateFunc func(entity *Entity, gs GameState), transitionFunc1 func(entity *Entity), transitionFunc2 func(entity *Entity)) *StateMachine {
	States := make(map[int]*StateHandler)
	idle := &StateHandler{Updater: UISpriteIdleUpdater, TransitionFunc: transitionFunc1, TransitionTo: 2}
	pickedUp := &StateHandler{Updater: updateFunc, TransitionFunc: transitionFunc2, TransitionTo: 1}
	States[1] = idle
	States[2] = pickedUp
	sm := &StateMachine{States: States, CurrentState: 1}
	return sm
}

func UISpriteIdleUpdater(ent *Entity, gs GameState) {
	if ent.UiData.Flags["used"] {
		return
	}
	if gs.FocusedEntity == ent {
		ent.StateMachine.Transition(ent)
	}
	us := ent.Sprite
	uiDat := ent.UiData
	if gs.HoveredUiSprite == ent {
		if uiDat.state != Disabled {
			if us.Img == uiDat.MainImg {
				if us.X == uiDat.baseX {
					us.Y -= 5
					us.X += 5
					us.Shader = registry.ShaderMap["Highlight"]
				}
			}
		}
	} else {
		us.Shader = nil
		us.X = uiDat.baseX
		us.Y = uiDat.baseY
	}
}

func PositionUpdate(ent *Entity, gs GameState) {
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		ent.StateMachine.Transition(ent)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		UnFocus(ent.Id)
	}

}

func UpdateSkimmer(ent *Entity, gs GameState) {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && ent.CursorInActivationRect() {
		UpdateEntityZAndReSortEntitySlice(ent.Id, 6)
		gs.CursorUpdater.ChangeSpeed(0.2)
		skimmerBounds := gs.Zbounds[0]
		skimmerBounds.Max.X += ent.Sprite.GetSpriteRect().Dx()
		skimmerBounds.Min.X -= ent.Sprite.GetSpriteRect().Dy()
		gs.CursorUpdater.SetBounds(gs.Zbounds[0])
	}
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		UpdateEntityZAndReSortEntitySlice(ent.Id, 13)
		gs.CursorUpdater.ResetSpeed()
		gs.CursorUpdater.ResetBounds()
		if inpututil.IsKeyJustPressed(ebiten.KeyE) {
			UiSpriteTurnOffEverything(ent)
			ent.StateMachine.Transition(ent)
		}
	}

}

func UpdateUseOnTank(ent *Entity, gs GameState) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && ent.CursorInActivationRect() {
		ent.StateMachine.Transition(ent)
	}
}

func UseOnTank(ent *Entity) {

	ent.UiData.state = Disabled
	UnFocus(ent.Id)
	ent.Sprite.XYUpdater = nil
	ent.UiData.Flags["used"] = true
	ent.UiData.timers["waitAndReset"] = util.NewTimer(2.5)
	ent.UiData.timers["waitAndReset"].TurnOn()
	ent.Sprite.DOptsUpdaterTag = "swirl"
	if ent.UiData.HoverImg != nil {
		ent.Sprite.Img = ent.UiData.HoverImg
	}
	if ent.Z != 6 {
		UpdateEntityZAndReSortEntitySlice(ent.Id, 6)
	}
	if !ent.UiData.Flags["particlesGenerated"] {
		spriteBounds := ent.Sprite.GetSpriteRect()
		fps := NewFertilizerParticleSystem(float64(ent.Sprite.X+float32(spriteBounds.Dx()/4)), float64(ent.Sprite.Y+float32(spriteBounds.Dy()/2)), ent.UiData.ActivationRect)
		fpent := &Entity{ParticleSystem: fps}
		fpent.Z = 6
		fpent.LifeTime = 8.0
		println("registering fertilizer particle entity")
		RegisterEntity(fpent)
		ent.UiData.Flags["particlesGenerated"] = true
	}

}

func (ent *Entity) CursorInActivationRect() bool {
	x, y := util.GetScaledCursorPosition()
	pt := image.Point{x, y}
	return pt.In(ent.UiData.ActivationRect)
}

func AddUiSpriteXYUpdater(ent *Entity) {
	if ent.UiData.Flags["used"] {
		return
	}
	ent.Sprite.XYUpdater = sprite.NewUpdater(ent.Sprite)
	if ent.UiData.Label == string(Skimmer) {
		ent.Sprite.XYUpdater.SetCustomOffset(0.0, 25)
	}
}

/*func (us *UiSpriteData) UpdatePosition() {
	if us.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		us.XYUpdater = sprite.NewUpdater(us.Sprite)
		us.Shader = registry.ShaderMap["Outline"]
	}

	if us.XYUpdater != nil {
		us.XYUpdater.Update()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && us.XYUpdater != nil {
		us.XYUpdater = nil
	}

	UpdateUiSpriteTimers(us)
}*/

func UpdateUiSpriteTimers(ent *Entity, gs GameState) {
	us := ent.UiData
	for name, timer := range us.timers {
		switch name {

		case "clickMeBuffer":
			state := timer.Update()
			if state == util.Done {
				us.highlight = false
				timer.TurnOff()
			}

		case "graphicDeInit":
			state := timer.Update()
			if state == util.Done {
				graphics.DeInitGraphics(us.Sprite.PublishedGraphicId)
				timer.TurnOff()
			}
		case "waitAndReset":
			state := timer.Update()
			if state == util.Done {
				UiSpriteTurnOffEverything(ent)
				UnFocus(ent.Id)
				timer.TurnOff()
				gs.CursorUpdater.ResetBounds()
				gs.CursorUpdater.ResetSpeed()
				if ent.UiData.Flags["oneOff"] {
					RemoveEntity(ent.Id)
				}
			}

		}

	}
}

func (e *Entity) DebugActivationRect() {
	us := e.UiData

	if us.Label == string(Phreader) && us.gameMode == registry.Position && us.debugGraphicPublishedID == 0 {
		id := graphics.NewRectGraphic(us.ActivationRect, colornames.Cornflowerblue)
		us.debugGraphicPublishedID = id
	}

	if us.Label == string(FishFood) && us.gameMode == registry.Position && us.debugGraphicPublishedID == 0 {
		id := graphics.NewRectGraphic(us.ActivationRect, colornames.Yellowgreen)
		us.debugGraphicPublishedID = id
	}
}

func (e *Entity) specificUiSpriteBehaviourUpdater(gs *GameState) {
	uiDat := e.UiData

	switch uiDat.Label {
	case string(Phreader):
		PHReaderUpdate(e, *gs)
	case string(FishFood):
		FishFoodUpdate(e)
	case string(Thermometer):
		if uiDat.Sprite.Img == uiDat.HoverImg && inpututil.IsKeyJustPressed(ebiten.KeyE) {
			graphics.DeInitGraphics(uiDat.Sprite.PublishedGraphicId)
			uiDat.Sprite.Img = uiDat.MainImg
			UpdateEntityZAndReSortEntitySlice(e.Id, 2)
		}
	}
}

func (e *Entity) hoveredUpdater(gs *GameState) {
	us := e.Sprite
	uiDat := e.UiData

	if !us.Focused && uiDat.state != Disabled {
		if us.Img == uiDat.MainImg {
			if us.X == uiDat.baseX {
				us.Y -= 5
				us.X += 5
				UpdateEntityZAndReSortEntitySlice(e.Id, 2)
			}
			us.Shader = registry.ShaderMap["Highlight"]
		} else {
			us.Shader = nil
		}
	}
	switch uiDat.Label {
	case string(Phreader), string(Magazine), string(GrandpasJournal):
		uiDat.PublishPickedUpEventIfClicked()
	case string(Thermometer):
		AltImageWhenClickedUpdaterStatic(e)
	case string(Pillow), string(Door):
		if e.Draw && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			uiDat.Publish(events.UISpriteAction{UiSprite: uiDat.Label, UiSpriteAction: "clicked"})
			turnOffClickMeEffect(uiDat)
		}
	}
}

func PiggyBankMoneyAddedGraphic(amtAdded float64, basex float32, basey float32) *graphics.GraphicManager {
	gm := &graphics.GraphicManager{}

	for i := 0; i < int(amtAdded/.25); i++ {
		timer := util.NewTimer(0.2)
		gm.Timers = make(map[string]*util.Timer)
		gm.Timers["Trigger"] = timer
		spEff := entImportableLoaders.LoadStaticEffect("25c", basex, basey, "")
		params := make(map[string]any)
		params["opacity"] = float32(1.0)
		gs := graphics.SpriteGraphic{Sprite: *spEff, UpdateFunc: MoveSpriteToDestinationUp, Parameters: params}
		gs.SetDrawFunc(graphics.FadeIn)
		gm.GraphicsToBePublished = append(gm.GraphicsToBePublished, &gs)
	}
	if gm.Timers != nil {
		gm.Timers["Trigger"].TurnOn()
	}
	return gm
}

func UiSpriteTurnOffEverything(ent *Entity) {
	// all the shit that needs to happen when we set the ui sprite/ data back to idle
	uiDat := ent.UiData
	sp := ent.Sprite
	sp.CurrentAnimation = ""
	if ent.Z != 0 {
		UpdateEntityZAndReSortEntitySlice(ent.Id, 0)
	}
	if len(uiDat.PublishedGraphicId) != 0 {
		graphics.DeInitGraphics(uiDat.PublishedGraphicId)
	}
	sp.Shader = nil
	sp.LinkedSprite = nil
	uiDat.state = Idle
	sp.XYUpdater = nil
	uiDat.returnToBase()
	sp.Img = uiDat.MainImg
	uiDat.DOptsUpdaterParams = make(map[string]float64)
}

func ClickForTime(ent *Entity, gs GameState, doAtTime func(ent *Entity)) {
	us := ent.Sprite
	uiDat := ent.UiData
	x, y := util.GetScaledCursorPosition()
	pt := image.Point{x, y}

	if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 120 && pt.In(uiDat.ActivationRect) {
		uiDat.Flags["clickForTime"] = false
		ent.Sprite.UpdateFunc = nil
		uiDat.state = Animation
		us.XYUpdater = nil
		doAtTime(ent)
	}
}

func MoveSpriteToDestination(sp *sprite.Sprite) {

	destinationX := 420.0
	destinationY := float64(registry.Config.ScreenHeight / 4)
	speed := 8.0

	// Calculate rotation needed to reach π (flipped)

	// Calculate the distance to destination
	dx := destinationX - float64(sp.X)
	dy := destinationY - float64(sp.Y)

	// Calculate the total distance
	distance := math.Sqrt(dx*dx + dy*dy)

	// If we're close enough, stop moving
	if distance < speed {
		sp.X = float32(destinationX)
		sp.Y = float32(destinationY)
		sp.UpdateFunc = nil

		return
	}
	sp.X += float32(dx / distance * speed)
	sp.Y += float32(dy / distance * speed)

}

func MoveSpriteToDestinationAndSpin(ui *sprite.Sprite) {

	destinationX := 400.0
	destinationY := 230.0
	speed := 4.0

	// Calculate rotation needed to reach π (flipped)

	// Calculate the distance to destination
	dx := destinationX - float64(ui.X)
	dy := destinationY - float64(ui.Y)

	// Calculate the total distance
	distance := math.Sqrt(dx*dx + dy*dy)

	// If we're close enough, stop moving
	if distance < speed {
		ui.X = float32(destinationX)
		ui.Y = float32(destinationY)
		ui.DOptsUpdaterParams["degree"] = math.Pi
		ui.UpdateFunc = nil
		return
	}

	travelTime := distance / speed

	targetRotation := math.Pi
	rotationNeeded := targetRotation - ui.DOptsUpdaterParams["degree"]

	// Handle rotation wrapping (if current rotation > π)
	if rotationNeeded < 0 {
		rotationNeeded += 2 * math.Pi
	}

	// Calculate rotation speed to arrive flipped
	rotationSpeed := rotationNeeded / travelTime
	ui.DOptsUpdaterParams["degree"] += rotationSpeed

	ui.X += float32(dx / distance * speed)
	ui.Y += float32(dy / distance * speed)

}

func (e *Entity) UpdateUiSprite(gs *GameState) {
	UpdateUiSpriteTimers(e, *gs)
	if e.StateMachine != nil {
		return
	}

	uiDat := e.UiData
	us := e.Sprite

	if uiDat.Label == string(WhiteBoard) {
		return
	}

	if e == gs.HoveredUiSprite {
		e.hoveredUpdater(gs)
	}

	if !us.Focused {
		if uiDat.state != Idle && uiDat.state != Disabled && uiDat.state != Clickable {
			uiDat.state = Idle
			UiSpriteTurnOffEverything(e)
		}
		if e != gs.HoveredUiSprite && us.X != uiDat.baseX {
			us.X = uiDat.baseX
			us.Y = uiDat.baseY
			UpdateEntityZAndReSortEntitySlice(e.Id, 0)
			us.Shader = nil
		}
		return
	}

	us.Scale = 1.0
	e.DebugActivationRect()

	if uiDat.state == Idle {
		uiDat.state = JustFocused
	}

	e.specificUiSpriteBehaviourUpdater(gs)
}

func (us *UiSpriteData) returnToBase() {
	us.highlight = false
	us.X = us.baseX
	us.Y = us.baseY
	if !registry.Config.Zoom {
		//if we are zooming in we dont want to play any sound effects
		ev := events.UISpriteAction{}
		ev.UiSprite = us.Label
		ev.UiSpriteAction = "put back"
		us.EventHub.Publish(ev)
	}

	us.XYUpdater = nil
}

func AltImageWhenClickedUpdaterStatic(ent *Entity) {
	us := ent.UiData
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && us.Img != us.HoverImg {
		println("thermometer message triggered")
		if us.Label == "thermometer" {
			AddTempGuage(us)
		}
		us.Scale = 1.0
		us.Sprite.Img = us.HoverImg
		UpdateEntityZAndReSortEntitySlice(ent.Id, 13)
	}

}

func AltImageWhenClickedUpdater(ent *Entity, gs GameState) {
	if ent.Sprite.Img != ent.UiData.HoverImg {
		img := ent.UiData.HoverImg
		x, y := util.GetScaledCursorPosition()
		sp := ent.Sprite
		sp.X = float32(x - img.Bounds().Dx()/2)
		sp.Y = float32(y - img.Bounds().Dy()/2)
		sp.Img = img
		sp.XYUpdater = sprite.NewUpdater(sp)
		if ent.Z < 2 {
			UpdateEntityZAndReSortEntitySlice(ent.Id, 2)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		ent.Sprite.Img = ent.UiData.MainImg
		ent.StateMachine.Transition(ent)
	}

}

func AddTempGuage(us *UiSpriteData) {
	width := float32(4)
	x := float32(us.HoverImg.Bounds().Dx() / 2)
	y := float32(us.HoverImg.Bounds().Dy() - 3)
	height := float32(us.Environment.Temperature-62) * 5
	vector.StrokeLine(us.HoverImg, x, y, x, y-height, width, colornames.Red, false)
	us.X = us.X - float32(us.HoverImg.Bounds().Dx()/4)
	us.Sprite.SavePublishedGraphicID(graphics.NewFadeInTextGraphicSmall(
		"Temperature:"+strconv.Itoa(us.Environment.Temperature),
		float64(us.X)/registry.Config.ResolutionScalingF+float64(us.HoverImg.Bounds().Dx()/2)/registry.Config.ResolutionScalingF,
		float64(us.Y)/registry.Config.ResolutionScalingF-float64(us.HoverImg.Bounds().Dy())/4,
	))
}

func AddTextGraphic(sp *sprite.Sprite, text string) int {
	id := graphics.NewFadeInTextGraphic(text, float64(sp.X)-float64(sp.Img.Bounds().Dx()), float64(sp.Y)-float64(sp.Img.Bounds().Dy()))
	return id
}

func NewUiSprite(environment *system.Environment, imgs []*ebiten.Image, hub *tasks.EventHub, x, y float32, label string) *UiSpriteData {

	var paramaMappa = make(map[string]any)

	uis := UiSpriteData{Sprite: &sprite.Sprite{X: x, Y: y, Scale: 1.0, ShaderParams: paramaMappa}}
	uis.ShaderParams = paramaMappa
	uis.baseX = x
	uis.baseY = y
	uis.Label = label
	uis.EventHub = hub
	uis.Environment = environment
	uis.ShaderParams = paramaMappa

	uis.timers = map[string]*util.Timer{}
	uis.timers["clickMeBuffer"] = util.NewTimer(1)
	uis.timers["graphicDeInit"] = util.NewTimer(3)
	uis.timers["waitAndReset"] = util.NewTimer(1)

	uis.Img = &ebiten.Image{}
	uis.Img = imgs[0]
	uis.MainImg = imgs[0]

	//set alt img + offset for alt
	if len(imgs) > 1 {
		uis.HoverImg = imgs[1]
		x1 := imgs[0].Bounds().Dx()
		y1 := imgs[0].Bounds().Dy()

		x2 := imgs[1].Bounds().Dx()
		y2 := imgs[1].Bounds().Dy()

		uis.AltOffsetX = float32(x1 - x2)
		uis.AltOffsetY = float32(y1 - y2)
	}

	if len(imgs) > 2 {
		uis.AltImg = imgs[2]
	}

	uis.state = Idle
	if label == string(PiggyBank) {
		uis.state = Disabled
	}
	uis.gameMode = registry.Normal

	uis.stringVariables = make(map[string]string)
	uis.Flags = make(map[string]bool)
	uis.variables = make(map[string]float64)

	return &uis
}

func UiSpriteSubs(hub *tasks.EventHub, uis *Entity) {

	hub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		switch ev.ButtonText {
		case "Mode":
			uis.UiData.SwitchGameMode()

		}
	})

	switch uis.UiData.Label {
	case string(PiggyBank):
		hub.Subscribe(events.MoneyAvailable{}, func(e tasks.Event) {
			ev := e.(events.MoneyAvailable)
			uis.Sprite.DOptsUpdaterTag = "swirl"
			uis.Sprite.Shader = registry.ShaderMap["Highlight"]
			uis.UiData.state = Clickable
			uis.UiData.variables["amountAvailable"] = ev.Amount
		})
	case string(Pillow):

		hub.Subscribe(tasks.AllTasksCompleted{}, func(e tasks.Event) {
			println("got all tasks done at pillow, drawing sprite")
			initClickMeEffect(uis.UiData)
			uis.Draw = true
		})

		hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
			uis.Draw = false
		})

	case string(Door):
		hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
			ev := e.(events.NewDay)
			if ev.DayType == "Chores" || ev.DayType == "Camp" {
				initClickMeEffect(uis.UiData)
				uis.Draw = true
			}
		})

		hub.Subscribe(events.LeavingFishScene{}, func(e tasks.Event) {
			uis.Draw = false
		})

	case string(Phreader):
		hub.Subscribe(events.PHGuess{}, func(e tasks.Event) {
			phev := e.(events.PHGuess)
			var text string
			var wbText string
			if math.Abs(phev.Guess-uis.UiData.Environment.ModifiedPHLevel) < .1 {
				text = "Right On! Stash you bonus quarter in your piggy bank"
				wbText = "PH: " + strconv.FormatFloat(phev.Guess, 'f', 2, 32)
				mev := events.MoneyAvailable{Amount: .25}
				hub.Publish(mev)
			} else if phev.Guess < uis.UiData.Environment.ModifiedPHLevel {
				text = "Too Low"
				wbText = "PH: >" + strconv.FormatFloat(phev.Guess, 'f', 1, 32)
			} else if phev.Guess > uis.UiData.Environment.ModifiedPHLevel {
				text = "Too High"
				wbText = "PH: <" + strconv.FormatFloat(phev.Guess, 'f', 1, 32)
			}
			if uis.Sprite.LinkedSprite != nil {
				uis.Sprite.SavePublishedGraphicID(AddTextGraphic(uis.Sprite.LinkedSprite, text))
			}
			fmt.Printf("Before replace: '%s'\n", wbText)
			wbText = strings.ReplaceAll(wbText, ".", " . ")
			fmt.Printf("After replace: '%s'\n", wbText)

			ev := events.WriteToWhiteBoard{PreferredPosition: "bottomLeft", Msg: wbText}
			hub.Publish(ev)
		})

		hub.Subscribe(events.DayOver{}, func(e tasks.Event) {
			uis.UiData.Flags["triggered"] = false
		})
	case string(Magazine):
		hub.Subscribe(tasks.TaskCreated{}, func(e tasks.Event) {
			ev := e.(tasks.TaskCreated)
			if ev.Task.Text == "3. Buy a new fish." {
				initClickMeEffect(uis.UiData)
			}
		})
		hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
			ev := e.(events.NewDay)
			if ev.Day > 1 {
				uis.Draw = true
			}
		})
	case string(GrandpasJournal):
		hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
			uis.Sprite.DOptsUpdaterTag = "swirl"
		})
		hub.Subscribe(events.UISpriteAction{}, func(e tasks.Event) {
			ev := e.(events.UISpriteAction)

			if ev.UiSprite == string(GrandpasJournal) {
				if uis.Sprite.DOptsUpdaterTag == "swirl" {
					uis.Sprite.DOptsUpdaterTag = ""
				}

			}
		})
	}
}

func (us *UiSpriteData) SwitchGameMode() {
	switch us.gameMode {
	case registry.Normal:
		us.gameMode = registry.Position
	case registry.Position:
		us.gameMode = registry.Normal
		if us.debugGraphicPublishedID != 0 {
			graphics.DeInitGraphicId(us.debugGraphicPublishedID)
			us.debugGraphicPublishedID = 0
		}
	}
}

func (us *UiSpriteData) SavePosition() drawables.SavePositionData {

	sp := drawables.SavePositionData{}
	sp.X = us.X
	sp.Y = us.Y
	sp.Name = us.Label
	return sp
}

func (us *UiSpriteData) PublishPickedUpEventIfClicked() {
	if us.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ev := events.UISpriteAction{
			UiSprite:       us.Label,
			UiSpriteAction: "picked up",
		}
		us.EventHub.Publish(ev)
	}

}

func initClickMeEffect(us *UiSpriteData) {
	cs := ebiten.ColorScale{}
	cs.SetR(0.1)
	cs.SetB(0.2)
	cs.SetG(1.0)
	cs.SetA(1.0)
	msg := "Click Me"
	us.Sprite.SavePublishedGraphicID(graphics.NewOutlineGraphicText(&msg, 24, float64(us.X), float64(us.Y), true, cs, float64(us.Img.Bounds().Dx()), true))

	ols := shaders.LoadOutlineShader()

	us.Sprite.Shader = ols
	us.Sprite.ShaderParams["Opacity"] = float32(0.0)
	us.Sprite.ShaderParams["OutlineColor"] = [4]float32{0.2, 0.7, 0.2, 1.0}
	us.Sprite.UpdateShaderParams = shaders.UpdatePulseWithText
}

func turnOffClickMeEffect(us *UiSpriteData) {

	us.Sprite.Shader = nil
	if us.timers != nil {
		us.timers["clickMeBuffer"].TurnOn()
	}
	graphics.DeInitGraphics(us.Sprite.PublishedGraphicId)
}

func FlipHighlight(us *UiSpriteData) {
	switch us.highlight {
	case true:
		us.highlight = false
	case false:
		us.highlight = true
	}
}
