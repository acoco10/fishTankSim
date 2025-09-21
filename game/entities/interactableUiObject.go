package entities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/drawables"
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
	"image/color"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

func UnFocusCheck() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyE) || inpututil.IsKeyJustPressed(ebiten.KeyEscape)
}

func ClickCheck() bool { return inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) }

func ClickDuration() int { return inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) }

const MidLayerZ = 6
const NotInTankZ = 13

type UiSpriteState uint8

const (
	Selected UiSpriteState = iota
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
	BaseZ                  int
	DontDraw               bool
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
	Timers                    map[string]*util.Timer
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

func (sm *StateMachine) AppendState(newUpdater EntityUpdater, newtransitionFunc EntityTransitioner) {
	if sm.States[len(sm.States)] != nil {
		sm.States[len(sm.States)].TransitionTo = len(sm.States) + 1
		newState := &StateHandler{Updater: newUpdater, TransitionFunc: newtransitionFunc, TransitionTo: 1}
		sm.States[len(sm.States)+1] = newState
	}
}

func InitStateMachine(initState EntityUpdater, updateFunc EntityUpdater, transitionFunc1 EntityTransitioner, transitionFunc2 EntityTransitioner) *StateMachine {
	States := make(map[int]*StateHandler)
	idleFunc := &StateHandler{Updater: initState, TransitionFunc: transitionFunc1, TransitionTo: 2}

	//set to default idle updater if no custom state is provided`
	if initState == nil {
		idleFunc = &StateHandler{Updater: UISpriteIdleUpdater, TransitionFunc: transitionFunc1, TransitionTo: 2}
	}
	if updateFunc != nil {
		pickedUp := &StateHandler{Updater: updateFunc, TransitionFunc: transitionFunc2, TransitionTo: 1}
		States[2] = pickedUp
	}
	States[1] = idleFunc
	sm := &StateMachine{States: States, CurrentState: 1}
	return sm
}

func AddPHEffectParticles(rect image.Rectangle, textureTag uint32) {
	for i := 0; i < 3; i++ {
		x := float64(50 + rect.Min.X + (i * 100))
		y := float64(rect.Max.Y) - 45 - 3*rand.NormFloat64()
		fps := NewGenericParticleSystem(x, y, rect, textureTag)

		fps.PConfig = &ParticleConfig{
			XVariance:         100,
			YVariance:         20,
			XVelocityVariance: 15,
			YVelocityVariance: 50,
			BaseYVelocity:     -150,
			MaxLife:           5,
			Scale:             1.0,
			AlphaDecay:        0.1,
			RotationSpeed:     0.01,
		}

		fpent := &Entity{ParticleSystem: fps, Sprite: fps.Sprite}
		fpent.Z = MidLayerZ
		fpent.LifeTime = 4.0
		RegisterEntity(fpent)
	}
}

func (ent *Entity) CursorInActivationRect() bool {
	x, y := util.GetScaledCursorPosition()
	pt := image.Point{x, y}
	return pt.In(ent.UiData.ActivationRect)
}

func AddUiSpriteXYUpdater(ent *Entity) {
	if ent.UiData.Flags["done"] {
		return
	}
	ent.Sprite.XYUpdater = sprite.NewUpdater(ent.Sprite)
	if ent.UiData.Label == string(Skimmer) {
		ent.Sprite.XYUpdater.SetCustomOffset(0.0, 25)
	}
}

func UpdateUiSpriteTimers(ent *Entity, gs GameState) {
	us := ent.UiData
	for name, timer := range us.Timers {
		if timer.TimerUpdater != nil {
			timer.TimerUpdater(timer, ent)
		}
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
					fmt.Println("Item used:", ent.UiData.Label)
					ent.EventHub.Publish(events.ItemUsed{Name: ent.UiData.Label})
					RemoveEntity(ent.Id)
				}
			}
		case "transition1":
			state := timer.Update()
			if state == util.Done {
				ent.StateMachine.Transition(ent)
				timer.TurnOff()
			}
		case "freeze":
			state := timer.Update()
			if state == util.Done {
				ent.UiData.Flags["freezeDone"] = true
				timer.TurnOff()
			}
		case "unFocusBuffer":
			state := timer.Update()
			if state == util.Done {
				//ent.UiData.state = Idle
				timer.TurnOff()
			}
		case DoAtTime:
			state := timer.Update()
			if state == util.Done {
				ent.DoAt[DoAtTime](ent, gs)
				timer.TurnOff()
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
	if gs.MouseFlags.WindowOpen {
		return
	}
	uiDat := e.UiData

	switch uiDat.Label {
	case string(Phreader):
		PHReaderUpdate(e, *gs)
	case string(FishFood):
		FishFoodUpdate(e)
	case string(Thermometer):
		if uiDat.Sprite.Img == uiDat.HoverImg && UnFocusCheck() {
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
			if us.X == uiDat.baseX && !uiDat.Flags["noOffset"] {
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
	case string(Phreader), string(Magazine):
		PublishPickedUpEventIfClicked(e, *gs)
	case string(LightSwitch):
		if e.Draw && ClickCheck() {
			uiDat.Publish(events.UISpriteAction{UiSprite: uiDat.Label, UiSpriteAction: "clicked"})
			e.effectHandler()
			UiSpriteTurnOffEverything(e)

			//turnOffClickMeEffect(uiDat)
		}
	case string(FishFood):
		PublishGraphicHovered(e, *gs)
	}
}

func UiSpriteTurnOffEverything(ent *Entity) {
	// all the shit that needs to happen when we set the ui sprite/ data back to idle
	uiDat := ent.UiData
	sp := ent.Sprite
	sp.CurrentAnimation = ""
	if ent.Z != ent.UiData.BaseZ {
		UpdateEntityZAndReSortEntitySlice(ent.Id, uiDat.BaseZ)
	}
	if len(uiDat.PublishedGraphicId) != 0 {
		graphics.DeInitGraphics(uiDat.PublishedGraphicId)
	}

	if !ent.UiData.Flags["keepShader"] {
		sp.Shader = nil
	}
	sp.LinkedSprite = nil
	uiDat.state = Idle
	sp.XYUpdater = nil
	uiDat.Scale = 0.0
	uiDat.returnToBase()
	sp.Img = uiDat.MainImg
	if uiDat.DontDraw {
		ent.Draw = false
	}

	uiDat.DOptsUpdaterParams = make(map[string]float64)
}

func ClickForTime(ent *Entity, gs GameState, doAtTime func(ent *Entity)) {
	us := ent.Sprite
	uiDat := ent.UiData
	x, y := util.GetScaledCursorPosition()
	pt := image.Point{x, y}

	if ClickDuration() > 120 && pt.In(uiDat.ActivationRect) {
		uiDat.Flags["clickForTime"] = false
		ent.Sprite.UpdateFunc = nil
		uiDat.state = Animation
		us.XYUpdater = nil
		doAtTime(ent)
	}
}

func MoveSpriteToDestination(sp *sprite.Sprite) {
	destinationX := sp.DOptsUpdaterParams["destinationX"]
	destinationY := sp.DOptsUpdaterParams["destinationY"]
	speed := sp.DOptsUpdaterParams["speed"]

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

	destinationX := 250.0
	destinationY := 50.0
	maxScale := 4.0
	speed := 8.0

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

	scaleIncreaseNeeded := maxScale - ui.Scale
	totalScaleChange := maxScale - 1
	ui.Scale += scaleIncreaseNeeded / totalScaleChange
	ui.X += float32(dx / distance * speed)
	ui.Y += float32(dy / distance * speed)

}

func (e *Entity) UpdateUiSprite(gs *GameState) {
	if registry.Config.Zoom {
		return
	}

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
	} else {
		if uiDat.Label == string(GrandpasJournal) {
			uiDat.Img = uiDat.MainImg
		}
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

	//us.Scale = 1.0 why is this here?
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

func CenterSprite(ent *Entity) {
	w := registry.Config.ScreenWidth
	h := registry.Config.ScreenHeight

	x := w/2 - ent.Sprite.GetSpriteRect().Dx()/2
	y := h/2 - ent.Sprite.GetSpriteRect().Dy()/2

	ent.Sprite.X = float32(x)
	ent.Sprite.Y = float32(y)
}

func MoveToCenter(ent *Entity, speed float64) {
	w := registry.Config.ScreenWidth
	h := registry.Config.ScreenHeight

	x := w/2 - ent.Sprite.GetSpriteRect().Dx()/2
	y := h/2 - ent.Sprite.GetSpriteRect().Dy()/2

	ent.Sprite.UpdateFunc = MoveSpriteToDestination
	ent.Sprite.DOptsUpdaterParams["destinationY"] = float64(y)
	ent.Sprite.DOptsUpdaterParams["destinationX"] = float64(x)
	ent.Sprite.DOptsUpdaterParams["speed"] = speed
}

func AddTempGuage(ent *Entity) {
	us := ent.UiData
	width := float32(4)
	x := float32(us.HoverImg.Bounds().Dx()/2) + 1
	y := float32(us.HoverImg.Bounds().Dy() - 2)
	height := float32(us.Environment.Temperature-62) * 5
	colr := color.RGBA{255, 20, 10, 150}
	vector.StrokeLine(us.HoverImg, x, y, x, y-height, width, colr, false)
	us.X = us.X - float32(us.HoverImg.Bounds().Dx()/4)
	us.Sprite.SavePublishedGraphicID(graphics.NewFadeInTextGraphic(
		"Temperature:"+strconv.Itoa(us.Environment.Temperature),
		float64(int(us.X)+us.HoverImg.Bounds().Dx()/2)+20, float64(us.Y),
		0,
	))
	ent.EventHub.Publish(WriteToWhiteBoard{Msg: fmt.Sprintf("Temp: %d", us.Environment.Temperature), PreferredPosition: "bottomRight", NoErase: true})
}

func AddTextGraphic(sp *sprite.Sprite, text string) int {
	id := graphics.NewFadeInTextGraphic(text, float64(sp.X)-float64(sp.Img.Bounds().Dx()), float64(sp.Y)-float64(sp.Img.Bounds().Dy()), 0)
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

	uis.Timers = map[string]*util.Timer{}
	uis.Timers["clickMeBuffer"] = util.NewTimer(1)
	uis.Timers["graphicDeInit"] = util.NewTimer(3)
	uis.Timers["waitAndReset"] = util.NewTimer(1)
	uis.Timers["unFocusBuffer"] = util.NewTimer(0.5)

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
	if label == string(PiggyBank) || label == string(LightSwitch) {
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
			uis.StateMachine.Transition(uis)
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

	case string(Thermometer):
		hub.Subscribe(events.Zoom{}, func(e tasks.Event) {
			uis.Sprite.DOptsUpdaterParams["opacity"] = 0.15
		})
		hub.Subscribe(events.UnZoom{}, func(e tasks.Event) {
			uis.Sprite.DOptsUpdaterParams["opacity"] = 0
		})
	case string(LightSwitch):
		hub.Subscribe(events.BedTime{}, func(e tasks.Event) {

			uis.UiData.state = Clickable
			uis.UiData.Shader = registry.ShaderMap["PulseHighlight"]
			uis.Sprite.ShaderParams["Counter"] = 0
			uis.Sprite.ShaderParams["MaxCounter"] = uis.Sprite.GetSpriteRect().Dy() * 10
			uis.UiData.Flags["KeepShader"] = true
			uis.Sprite.UpdateShaderParams = shaders.UpdateCounter
			uis.effectHandler = LoadFollowEffectAsEnt("exclamation", uis.Id, hub, nil)
		})
		hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
			uis.Sprite.Shader = nil
		})

	case string(Phreader):
		hub.Subscribe(events.PHGuess{}, func(e tasks.Event) {
			UiSpriteTurnOffEverything(uis)
			uis.StateMachine.Reset(uis)
			UnFocus(uis.Id)
			phev := e.(events.PHGuess)
			var text string
			var wbText string
			if math.Abs(phev.Guess-uis.UiData.Environment.ModifiedPHLevel) < .1 {
				text = "Right On!"
				wbText = "PH: " + strconv.FormatFloat(phev.Guess, 'f', 2, 32)
				mev := events.MoneyAvailable{Amount: .25}
				hub.Publish(mev)
			} else if phev.Guess < uis.UiData.Environment.ModifiedPHLevel {
				text = "Too Low!"
				wbText = "PH: >" + strconv.FormatFloat(phev.Guess, 'f', 1, 32)
			} else if phev.Guess > uis.UiData.Environment.ModifiedPHLevel {
				text = "Too High!"
				wbText = "PH: <" + strconv.FormatFloat(phev.Guess, 'f', 1, 32)
			}
			if uis.Sprite.LinkedSprite != nil {
				uis.Sprite.SavePublishedGraphicID(AddTextGraphic(uis.Sprite.LinkedSprite, text))
			}

			graphics.NewFadeInTextGraphicCentered(text, 120)
			ev := WriteToWhiteBoard{PreferredPosition: "bottomLeft", Msg: wbText, NoErase: true}
			hub.Publish(ev)
		})

		hub.Subscribe(events.DayOver{}, func(e tasks.Event) {
			uis.UiData.Flags["triggered"] = false
		})

	case string(GrandpasJournal):
		hub.Subscribe(events.NewSpecies{}, func(e tasks.Event) {
			params := make(map[string]any)
			params["position"] = "center"
			uis.effectHandler = LoadFollowEffectAsEnt("exclamation", uis.Id, hub, params)
			uis.Sprite.DOptsUpdaterTag = "swirl"
			uis.UiData.Shader = registry.ShaderMap["Highlight"]
			uis.UiData.Flags["keepShader"] = true
		})
		hub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
			ev := e.(events.ButtonClickedEvent)
			if strings.HasPrefix(ev.ButtonText, "FishUnlocked:") {
				uis.UiData.Flags["keepShader"] = false
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

func initClickMeEffect(us *UiSpriteData) {
	cs := ebiten.ColorScale{}
	cs.SetR(0.1)
	cs.SetB(0.2)
	cs.SetG(1.0)
	cs.SetA(1.0)
	msg := "Click Me"
	us.Sprite.SavePublishedGraphicID(graphics.NewPulseGraphic(msg, float64(us.X), float64(us.Y), 120))

	ols := registry.ShaderMap["Outline"]

	us.Sprite.Shader = ols
	us.Sprite.ShaderParams["Opacity"] = float32(0.0)
	us.Sprite.ShaderParams["OutlineColor"] = [4]float32{0.2, 0.7, 0.2, 1.0}
	us.Sprite.UpdateShaderParams = shaders.UpdatePulseWithText
}

func turnOffClickMeEffect(us *UiSpriteData) {

	us.Sprite.Shader = nil
	if us.Timers != nil {
		us.Timers["clickMeBuffer"].TurnOn()
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
