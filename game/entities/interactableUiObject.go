package entities

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/stringConstants"
	"github.com/acoco10/fishTankWebGame/game/system"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"math"
	"math/rand"
	"strconv"
)

func UnFocusCheck() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyE) || inpututil.IsKeyJustPressed(ebiten.KeyEscape)
}

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
	Day1EventIds           []int
	MainImg                *ebiten.Image
	HoverImg               *ebiten.Image
	AltImg                 *ebiten.Image
	NewItemCondition       func()
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
	stringVariables           map[string]string
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
	if ent.HasUsed() {
		return
	}

	ent.Sprite.XYUpdater = sprite.NewUpdater(ent.Sprite)
	if ent.UiData.Label == string(Skimmer) {
		ent.Sprite.Shader = nil
		ent.Sprite.XYUpdater.SetCustomOffset(0.0, 25)
		transitionToUpdateSkimmer(ent)
	}
}

func UpdateUiSpriteTimers(ent *Entity, gs GameState) {
	us := ent.UiData

	if ent.UiData.Timers != nil {
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
					if ent.HasOneOff() {
						fmt.Println("Item used:", ent.UiData.Label)
						ent.EventHub.Publish(events.ItemUsed{Name: ent.UiData.Label})
						RemoveEntity(ent.Id)
					}
				}
			case Transition:
				state := timer.Update()
				if state == util.Done {
					ent.StateMachine.Transition(ent)
					timer.TurnOff()
				}
			case Reset:
				state := timer.Update()
				if state == util.Done {
					ent.StateMachine.Reset(ent)
					turnOffEverythingUnFocusResetBundle(ent)
					UnFocus(ent.Id)
					timer.TurnOff()
					ent.ClearFreeze()
					if ent.LinkedID != 0 {
						ent2, _ := GetEntity(ent.LinkedID)
						ent2.StateMachine.Reset(ent2)
						turnOffEverythingUnFocusResetBundle(ent2)
						ent2.ClearFreeze()
					}
				}
			case Freeze:
				state := timer.Update()
				if state == util.Done {
					ent.ClearFreeze()
					timer.TurnOff()
				}
			case PhGuessAddMoneyDelay:
				state := timer.Update()
				if state == util.Done {
					mev := events.MoneyAvailable{Amount: .25}
					ent.EventHub.Publish(mev)
					timer.TurnOff()
				}
			case PublishAtTime:
				state := timer.Update()
				if state == util.Done {
					ent.EventHub.Publish(ent.Parameters.Events[EventAtTime])
					timer.TurnOff()
					if ent.UiData.Label == string(GrandpasJournal) {
						ent.Draw = false
						ent.EventHub.Subscribe(events.WindowClosed{}, func(e tasks.Event) {
							ent.Draw = true
						})
					}
				}
			case PickedUp:
				state := timer.Update()
				if state == util.Done {
					timer.TurnOff()
				}
			case "unFocusBuffer":
				state := timer.Update()
				if state == util.Done {
					timer.TurnOff()
				}
			case DoAtTime:
				state := timer.Update()
				if state == util.Done {
					ent.DoAt[DoAtTime](ent, gs)
					timer.TurnOff()
				}
			case GenFood:
				state := timer.Update()
				if state == util.Done {
					if ent.Parameters.Ints[IndexCounter] > 0 {
						ent.Parameters.Ints[IndexCounter]--
						ent.EventHub.Publish(ent.Parameters.Events[PointEvent])
					} else {
						ent.Sprite.Img = ent.UiData.MainImg
						ent.EventHub.Publish(RequestZoom{ZoomedForFeeding})
						timer.TurnOff()
					}
				}
			case OnEnterFreeze:
				state := timer.Update()
				if state == util.Done {
					timer.TurnOff()
				}
			}
		}
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

	if len(ent.PublishedGraphicIDs) != 0 {
		graphics.DeInitGraphics(ent.PublishedGraphicIDs)
	}

	ent.ClearKeepShader()

	sp.LinkedSprite = nil
	uiDat.state = Idle
	sp.XYUpdater = nil
	uiDat.Scale = 0.0
	uiDat.returnToBase()
	if uiDat.DontDraw {
		ent.UiData.state = Disabled
		//ent.Draw = false
	}

	sp.Img = uiDat.MainImg
	sp.DOptsUpdaterParams = make(map[string]float64)
}

func ClickForTime(ent *Entity, gs GameState, doAtTime func(ent *Entity)) {
	us := ent.Sprite
	uiDat := ent.UiData
	x, y := util.GetScaledCursorPosition()
	pt := image.Point{x, y}

	if registry.ClickDuration() > 120 && pt.In(uiDat.ActivationRect) {
		ent.ClearClickForTime()
		ent.Sprite.UpdateFunc = nil
		uiDat.state = Animation
		us.XYUpdater = nil
		doAtTime(ent)
	}
}

func MoveSpriteToDestination(sp *sprite.Sprite) {
	if sp.DOptsUpdaterTag == "cursor" {
		x, y := util.GetScaledCursorPosition()
		sp.DOptsUpdaterParams["destinationX"] = float64(x)
		sp.DOptsUpdaterParams["destinationY"] = float64(y)
	}

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

func (ent *Entity) UpdateUiSprite(gs *GameState) {
	if registry.Config.Zoom {
		return
	}

	UpdateUiSpriteTimers(ent, *gs)
	if gs.MouseFlags.WindowOpen {
		return
	}
	if !gs.MouseFlags.WindowOpen && registry.Config.Zoom {
		if ent.Sprite.X != ent.UiData.baseX {
			ent.Sprite.Shader = nil
			ent.Sprite.X = ent.UiData.baseX
			ent.Sprite.Y = ent.UiData.baseY
			ent.Z = 0
			ZSortEntities()
		}
	}
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
	tempGraphic := graphics.NewFadeInTextGraphic(
		"Temperature:"+strconv.Itoa(us.Environment.Temperature),
		float64(int(us.X)+us.HoverImg.Bounds().Dx()/2)+20, float64(us.Y),
		0,
	)
	ent.PublishedGraphicIDs = append(ent.PublishedGraphicIDs, tempGraphic)
	ent.EventHub.Publish(WriteToWhiteBoard{Msg: fmt.Sprintf("Temp: %d", us.Environment.Temperature), PreferredPosition: LowerRight, NoErase: true})
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

	//set alt particleImg + offset for alt
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
	uis.variables = make(map[string]float64)

	return &uis
}

func PublishNewItemText(newitem string) {
	graphics.NewFadeInTextGraphicCenteredFadeOnClick(fmt.Sprintf("New Item: %s!", newitem), 0)
}

func UiSpriteSubs(hub *tasks.EventHub, uis *Entity) {

	hub.Subscribe(events.Zoom{}, func(event tasks.Event) {
		uis.PreZoomDraw = uis.Draw
		uis.Draw = false
	})

	hub.Subscribe(events.UnZoom{}, func(event tasks.Event) {
		uis.Draw = uis.PreZoomDraw
	})

	hub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		switch ev.ButtonText {
		case "Mode":
			uis.UiData.SwitchGameMode()
		}
	})

	id := hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		ev := e.(events.NewDay)
		if ev.Day == 1 {
			if uis.HasDontDrawFirstDay() {
				uis.Draw = false
			}
		}
	})

	uis.EndOfDayNUnSubscribeEvents[1] = append(uis.EndOfDayNUnSubscribeEvents[1], tasks.CreatedEvent{Id: id, Ev: &events.NewDay{}}) //unsubscribe from all of these after day one

	switch Label(uis.UiData.Label) {

	case FishFood:
		tcID := hub.Subscribe(tasks.TaskCreated{}, func(e tasks.Event) {
			ev := e.(tasks.TaskCreated)
			if ev.Task.EventType.Type() == "AllFishFed" {
				PublishNewItemText("Fish Food")
				uis.Sprite.DOptsUpdaterTag = stringConstants.Swirl
				uis.Sprite.Shader = registry.ShaderMap[registry.Highlight]
				uis.Draw = true
			}
		})

		uis.EndOfDayNUnSubscribeEvents[1] = append(uis.EndOfDayNUnSubscribeEvents[1], tasks.CreatedEvent{Id: tcID, Ev: tasks.TaskCreated{}}) //unsubscribe from all of these after day one

		uis.UiData.Day1EventIds = append(uis.UiData.Day1EventIds, tcID) //unsubscribe from all of these after day one

	case PiggyBank:
		//day 1 case
		tcID := hub.Subscribe(events.MoneyAvailable{}, func(e tasks.Event) {
			if !uis.Draw {
				PublishNewItemText("Piggy Bank")
				uis.Draw = true
			}

		})

		uis.EndOfDayNUnSubscribeEvents[1] = append(uis.EndOfDayNUnSubscribeEvents[1], tasks.CreatedEvent{Id: tcID, Ev: events.MoneyAvailable{}}) //unsubscribe from all of these after day one
		//normal case
		hub.Subscribe(events.MoneyAvailable{}, func(e tasks.Event) {
			uis.SetAddIdleClickEffect()
			uis.SetNoActivationRect()
			ev := e.(events.MoneyAvailable)
			uis.Sprite.DOptsUpdaterTag = stringConstants.Swirl
			uis.Sprite.Shader = registry.ShaderMap["Highlight"]
			uis.Sprite.DOptsUpdaterParams["swirlSpeedX"] = 0.5
			uis.Sprite.DOptsUpdaterParams["swirlSpeedY"] = 5
			uis.UiData.state = Clickable
			uis.StateMachine.Transition(uis)
			uis.UiData.variables["amountAvailable"] = ev.Amount
		})

	case Thermometer:
		tcID := hub.Subscribe(tasks.TaskCreated{}, func(e tasks.Event) {
			ev := e.(tasks.TaskCreated)
			if ev.Task.EventType.Type() == "WriteToWhiteBoard" {
				PublishNewItemText("Thermometer")
				uis.Sprite.DOptsUpdaterTag = stringConstants.Swirl
				uis.Sprite.DOptsUpdaterParams["swirlSpeedX"] = 1
				uis.Sprite.Shader = registry.ShaderMap[registry.Highlight]
				uis.Draw = true
			}
		})

		uis.EndOfDayNUnSubscribeEvents[1] = append(uis.EndOfDayNUnSubscribeEvents[1], tasks.CreatedEvent{Id: tcID, Ev: tasks.TaskCreated{}}) //unsubscribe from all of these after day one

		hub.Subscribe(events.Zoom{}, func(e tasks.Event) {
			uis.Sprite.DOptsUpdaterParams["opacity"] = 0.15
		})
		hub.Subscribe(events.UnZoom{}, func(e tasks.Event) {
			uis.Sprite.DOptsUpdaterParams["opacity"] = 0
		})
	case LightSwitch:
		hub.Subscribe(events.BedTime{}, func(e tasks.Event) {
			ev := e.(events.BedTime)
			if ev.Day == 1 {
				// on day one dont prompt user to click light switch until they have browsed the stor for the first time
				hub.Subscribe(events.WindowClosed{}, func(e tasks.Event) {
					ev2 := e.(events.WindowClosed)
					if ev2.Window == string(Magazine) {
						uis.StateMachine.Transition(uis)
						uis.UiData.state = Clickable
						uis.Sprite.DOptsUpdaterTag = stringConstants.Swirl
						uis.Sprite.DOptsUpdaterParams["swirlSpeedX"] = 0.05
						uis.Sprite.DOptsUpdaterParams["swirlSpeedY"] = 2
						uis.Sprite.Shader = registry.ShaderMap[registry.Highlight]
						uis.SetKeepShader()
					}
				})
			} else {
				uis.UiData.state = Clickable
			}
		})
		hub.Subscribe(events.DayOver{}, func(e tasks.Event) {
			uis.ClearKeepShader()
			uis.Sprite.DOptsUpdaterTag = ""
			uis.Sprite.Shader = nil
			uis.UiData.state = Disabled
		})
	case Phreader:

		hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
			ev := e.(events.NewDay)
			if ev.Day != 1 {
				uis.UiData.state = Idle
				uis.ClearLowLight()
			}
		})
		tcID := hub.Subscribe(tasks.TaskCreated{}, func(e tasks.Event) {
			ev := e.(tasks.TaskCreated)
			if ev.Task.EventType.Type() == "PHGuess" {
				PublishNewItemText("pH Strips")
				uis.Sprite.DOptsUpdaterTag = stringConstants.Swirl
				uis.Sprite.Shader = registry.ShaderMap[registry.Highlight]
				uis.Draw = true
			}
		})

		uis.EndOfDayNUnSubscribeEvents[1] = append(uis.EndOfDayNUnSubscribeEvents[1], tasks.CreatedEvent{Id: tcID, Ev: tasks.TaskCreated{}}) //unsubscribe from all of these after day one
		hub.Subscribe(events.PHGuess{}, func(e tasks.Event) {
			UiSpriteTurnOffEverything(uis)
			uis.StateMachine.Reset(uis)
			UnFocus(uis.Id)
			phev := e.(events.PHGuess)
			text, wbText := makePhGuessMsg(uis.UiData.Environment.ModifiedPHLevel, phev.Guess)
			if text == "Right On!" {
				uis.UiData.Timers[PhGuessAddMoneyDelay].TurnOn()

			}

			if uis.Sprite.LinkedSprite != nil {
				uis.Sprite.SavePublishedGraphicID(AddTextGraphic(uis.Sprite.LinkedSprite, text))
			}

			graphics.NewFadeInTextGraphicCentered(text, 120)
			ev := WriteToWhiteBoard{PreferredPosition: LowerLeft, Msg: wbText, NoErase: true}
			hub.Publish(ev)
			uis.UiData.state = Disabled
			uis.SetLowLight()
		})
	case PlantPack:
		hub.Subscribe(events.UISpriteAction{}, func(e tasks.Event) {
			ev := e.(events.UISpriteAction)
			if ev.UiSprite == string(OldContainer) {
				if ev.UiSpriteAction == "picked upplant" {
					PublishNewItemText("Plant Pack")
					uis.Draw = true
				}
			}
		})
	case RockProp:
		hub.Subscribe(events.UISpriteAction{}, func(e tasks.Event) {
			ev := e.(events.UISpriteAction)
			if ev.UiSprite == string(OldContainer) {
				if ev.UiSpriteAction == "picked uprock" {
					PublishNewItemText("Boring rock")
					uis.Draw = true
				}
			}

		})
	case Magazine:
		tcID := hub.Subscribe(events.BedTime{}, func(e tasks.Event) {
			uis.Draw = true
			uis.UiData.state = Disabled
			uis.Sprite.DOptsUpdaterTag = stringConstants.Swirl
			uis.UiData.Shader = registry.ShaderMap["Highlight"]
			uis.SetKeepShader()
		})

		uis.EndOfDayNUnSubscribeEvents[1] = append(uis.EndOfDayNUnSubscribeEvents[1], tasks.CreatedEvent{Id: tcID, Ev: &tasks.TaskCreated{}}) //unsubscribe from all of these after day one
	case Skimmer:
		tcID := hub.Subscribe(tasks.TaskCreated{}, func(e tasks.Event) {
			ev := e.(tasks.TaskCreated)
			if ev.Task.EventType.Type() == "DebrisCaptured" {
				PublishNewItemText("Skimmer")
				uis.Sprite.DOptsUpdaterTag = stringConstants.Swirl
				uis.Sprite.Shader = registry.ShaderMap[registry.Highlight]
				uis.Draw = true
			}
		})
		uis.EndOfDayNUnSubscribeEvents[2] = append(uis.EndOfDayNUnSubscribeEvents[2], tasks.CreatedEvent{Id: tcID, Ev: tasks.TaskCreated{}})

	case GrandpasJournal:
		hub.Subscribe(events.UISpriteAction{}, func(e tasks.Event) {
			ev := e.(events.UISpriteAction)
			if ev.UiSprite == string(OldContainer) {
				if ev.UiSpriteAction == "picked upjournal" {
					PublishNewItemText("Journal")
					uis.SetAddIdleClickEffect()
					uis.Draw = true
				}
			}

		})

		hub.Subscribe(events.NewSpecies{}, func(e tasks.Event) {
			uis.Sprite.DOptsUpdaterTag = stringConstants.Swirl
			uis.UiData.Shader = registry.ShaderMap["Highlight"]
			uis.SetKeepShader()
		})
		tcID := hub.Subscribe(events.BedTime{}, func(e tasks.Event) {
			uis.UiData.MainImg = uis.Parameters.Images[Alternate]
			uis.Sprite.Img = uis.UiData.MainImg
			PublishNewItemText("Magazine")
			uis.Sprite.DOptsUpdaterTag = stringConstants.Swirl
			uis.Sprite.Shader = registry.ShaderMap[registry.Highlight]
			uis.SetKeepShader()
			uis.Draw = true
		})

		uis.EndOfDayNUnSubscribeEvents[1] = append(uis.EndOfDayNUnSubscribeEvents[1], tasks.CreatedEvent{Id: tcID, Ev: tasks.TaskCreated{}}) //unsubscribe from all of these after day one
	}
}

func makePhGuessMsg(actual float64, guess float64) (string, string) {
	var text string
	var wbText string

	dif := math.Abs(actual - guess)

	if dif < .1 {
		text = "Right On!"
		wbText = "PH: " + strconv.FormatFloat(guess, 'f', 2, 32)
	} else if dif < .3 {
		text += "A little "
	} else if dif < .7 {
		text += "Way "
	}

	if text != "Right On!" {
		if guess < actual {
			text += "too Low!"
			wbText = "PH: >" + strconv.FormatFloat(guess, 'f', 1, 32)
		} else {
			text += "too High!"
			wbText = "PH: <" + strconv.FormatFloat(guess, 'f', 1, 32)
		}
	}
	return text, wbText
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
