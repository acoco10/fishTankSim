package interactableUIObjects

import (
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/input"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/system"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"image"
	"math"
	"strconv"
)

type uiSpriteState uint8

const (
	Selected uiSpriteState = iota
	HoveredOver
	ClickedWhileBeingSelected
	Idle
	Clickable
	ExtraSpriteAnimationCompleted
	Animation
)

type gameMode uint8

const (
	Position gameMode = iota
	Normal
)

type UiSprite struct {
	*sprite.Sprite
	baseX, baseY           float32
	MainImg                *ebiten.Image
	HoverImg               *ebiten.Image
	AltImg                 *ebiten.Image
	AltOffsetX, AltOffsetY float32
	*sprite.XYUpdater
	*tasks.EventHub
	state    uiSpriteState
	stateWas uiSpriteState
	timers   map[string]*entities.Timer
	gameMode
	clicked                   bool
	Draggable                 bool
	Label                     string
	highlight                 bool
	screenHeight, screenWidth int
	graphicPublishedID        int
	extraSprite               *sprite.Sprite
	Environment               *system.Environment
	activationRect            image.Rectangle
}

func (us *UiSprite) Draw(screen *ebiten.Image) {
	if us.Label == string(Phreader) {
		//util.StrokeRectFromImageRect(us.activationRect, screen)
	}
	us.Sprite.Draw(screen)
	if us.extraSprite != nil {
		us.extraSprite.Draw(screen)
	}
}

func (us *UiSprite) Highlighted() bool {
	return us.highlight
}

func (us *UiSprite) Update() {

	switch us.gameMode {
	case Normal:
		us.UpdateNormal()
	case Position:
		us.UpdatePosition()
	}

}

func (us *UiSprite) UpdatePosition() {
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
}

func UpdateUiSpriteTimers(us *UiSprite) {
	for name, timer := range us.timers {
		switch name {

		case "clickMeBuffer":
			state := timer.Update()
			if state == entities.Done {
				us.highlight = false
				timer.TurnOff()

			}

		case "graphicDeInit":
			state := timer.Update()
			if state == entities.Done {
				graphics.DeInitGraphicId(us.graphicPublishedID)
				timer.TurnOff()
			}

		}
	}
}

func (us *UiSprite) specificBehaviourUpdater() {
	switch us.Label {
	case string(Thermometer):
		AltImageWhenClickedUpdaterStatic(us)
	case string(Magazine):
		us.PublishPickedUpEventIfClicked()
	case string(Phreader):
		if us.state == ExtraSpriteAnimationCompleted && inpututil.IsKeyJustPressed(ebiten.KeyE) {
			us.extraSprite = nil
			us.state = Idle
			us.returnToBase()
			us.Img = us.MainImg
			us.DoptsUpdaterParams = make(map[string]float64)
			ev := events.UISpriteAction{UiSprite: "phreader", UiSpriteAction: "ph reading"}
			us.EventHub.Publish(ev)
		}

		pt := image.Point{X: int(us.X), Y: int(us.Y)}

		if us.Sprite.SpriteHovered() {
			ClickForTime(us, phReaderDoAtTime)
		}

		if pt.In(us.activationRect) && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && us.state != Animation && us.state != ExtraSpriteAnimationCompleted {
			println("setting highlight to false")
			us.highlight = false
		}

		if us.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			AltImageWhenClickedUpdater(us)
		}
	}
}

func ClickForTime(us *UiSprite, doAtTime func(us *UiSprite)) {

	x, y := ebiten.CursorPosition()
	pt := image.Point{x, y}

	if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 120 && pt.In(us.activationRect) {
		us.state = Animation
		us.highlight = true
		println("setting highlight to true")
		us.XYUpdater = nil
		println("adding move sprite to destination updated to ph tab")
		doAtTime(us)
	}
}

func phReaderDoAtTime(us *UiSprite) {
	us.Sprite.Shader = registry.ShaderMap["PH"]
	us.Sprite.ShaderParams["PHValue"] = us.Environment.NaturalPHLevel
	us.Sprite.ShaderParams["Point"] = []float64{3, 10}
	us.Sprite.ShaderParams["Radius"] = 3.0
	us.Sprite.UpdateFunc = MoveSpriteToDestinationAndSpin

	sp := &sprite.Sprite{Img: us.AltImg, X: us.baseX, Y: us.baseY}
	sp.UpdateFunc = MoveSpriteToDestination
	us.extraSprite = sp
	us.state = ExtraSpriteAnimationCompleted
	//us.UpdateShaderParams = shaders.UpdateCounter
}

func MoveSpriteToDestination(ui *sprite.Sprite) {

	destinationX := 420.0
	destinationY := 310.0
	speed := 6.0

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
		ui.UpdateFunc = nil

		return
	}
	ui.X += float32(dx / distance * speed)
	ui.Y += float32(dy / distance * speed)

}

func MoveSpriteToDestinationAndSpin(ui *sprite.Sprite) {

	destinationX := 400.0
	destinationY := 400.0
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
		ui.DoptsUpdaterParams["degree"] = math.Pi
		ui.UpdateFunc = nil

		return
	}

	travelTime := distance / speed

	targetRotation := math.Pi
	rotationNeeded := targetRotation - ui.DoptsUpdaterParams["degree"]

	// Handle rotation wrapping (if current rotation > π)
	if rotationNeeded < 0 {
		rotationNeeded += 2 * math.Pi
	}

	// Calculate rotation speed to arrive flipped
	rotationSpeed := rotationNeeded / travelTime
	ui.DoptsUpdaterParams["degree"] += rotationSpeed

	ui.X += float32(dx / distance * speed)
	ui.Y += float32(dy / distance * speed)

}

func (us *UiSprite) UpdateNormal() {
	us.clicked = false
	UpdateUiSpriteTimers(us)
	us.specificBehaviourUpdater()
	us.Sprite.Update()
	if us.extraSprite != nil {
		us.extraSprite.Update()
	}

	if us.state == HoveredOver && us.stateWas != HoveredOver {
		if us.Img == us.MainImg {
			us.Shader = registry.ShaderMap["Outline"]
			us.ShaderParams = make(map[string]any)
			us.ShaderParams["OutlineColor"] = [4]float64{1, 1, 0, 1}
		}
	}

	if us.Shader != nil && (us.state != HoveredOver && us.state != Selected) {
		//
		//us.Shader = nil
	}

	if us.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && us.state == Selected && us.Draggable {

		if us.XYUpdater == nil {
			ev := events.UISpriteAction{}
			ev.UiSprite = us.Label
			ev.UiSpriteAction = "picked up"
			us.EventHub.Publish(ev)
		}

		us.XYUpdater = sprite.NewUpdater(us.Sprite)
	}

	if us.XYUpdater != nil {
		us.XYUpdater.Update()
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && us.state == ClickedWhileBeingSelected && !us.clicked {
		us.clicked = true
	}

	baseDis := math.Hypot(float64(us.X-us.baseX), float64(us.Y-us.baseY))

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if us.state == Selected && baseDis < 100 && us.stateWas == Selected {
			us.returnToBase()
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyE) {
		if baseDis != 0 {
			us.returnToBase()
		}
	}

	us.stateWas = us.state
}

func (us *UiSprite) returnToBase() {
	us.state = HoveredOver
	us.highlight = false
	us.X = us.baseX
	us.Y = us.baseY
	ev := events.UISpriteAction{}
	ev.UiSprite = us.Label
	ev.UiSpriteAction = "put back"
	us.EventHub.Publish(ev)
	us.XYUpdater = nil
}

func AltImageWhenClickedUpdaterStatic(us *UiSprite) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && us.SpriteHovered() && us.Img != us.HoverImg {
		println("thermometer message triggered")
		if us.Label == "thermometer" {
			AddTempGuage(us)
		}
		us.Img = us.HoverImg
		us.highlight = true
	}

	if us.Img == us.HoverImg && inpututil.IsKeyJustPressed(ebiten.KeyE) {
		graphics.DeInitGraphicId(us.graphicPublishedID)
		us.Img = us.MainImg
		us.highlight = false
		us.X = us.X + float32(us.HoverImg.Bounds().Dx())
	}

}

func AltImageWhenClickedUpdater(us *UiSprite) {
	x, y := ebiten.CursorPosition()

	us.X = float32(x - us.HoverImg.Bounds().Dx()/2)
	us.Y = float32(y - us.HoverImg.Bounds().Dy()/2)
	us.Img = us.HoverImg
	us.XYUpdater = sprite.NewUpdater(us.Sprite)
}

func AddTempGuage(us *UiSprite) {
	width := float32(4)
	x := float32(us.HoverImg.Bounds().Dx()/2) - 1
	y := float32(us.HoverImg.Bounds().Dy() - 3)
	height := float32(us.Environment.Temperature-62) * 5
	vector.StrokeLine(us.HoverImg, x, y, x, y-height, width, colornames.Red, false)
	us.X = us.X - float32(us.HoverImg.Bounds().Dx()/4)
	us.graphicPublishedID = graphics.NewFadeInTextGraphic("Temperature:"+strconv.Itoa(us.Environment.Temperature), float64(us.X+80), float64(us.Y))
}

func AddTextGraphic(us *UiSprite, text string, x float64, y float64) {
	us.graphicPublishedID = graphics.NewFadeInTextGraphic(text, x, y)
}

func NewUiSprite(environment *system.Environment, imgs []*ebiten.Image, hub *tasks.EventHub, x, y float32, label string, screenWidth, screenHeight int) *UiSprite {

	var paramaMappa = make(map[string]any)

	uis := UiSprite{Sprite: &sprite.Sprite{X: x, Y: y}}
	uis.ShaderParams = paramaMappa
	uis.baseX = x
	uis.baseY = y
	uis.Label = label
	uis.EventHub = hub
	uis.Environment = environment

	uis.timers = map[string]*entities.Timer{}
	uis.timers["clickMeBuffer"] = entities.NewTimer(1)
	uis.timers["graphicDeInit"] = entities.NewTimer(3)

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

	Subs(hub, &uis)

	uis.state = Idle
	uis.gameMode = Normal

	return &uis
}

func Pubs(hub *tasks.EventHub, uis UiSprite) {
	ev := events.UISpriteLayedOut{
		Label: uis.Label,
		X:     uis.X,
		Y:     uis.Y,
	}
	hub.Publish(ev)
}

func Subs(hub *tasks.EventHub, uis *UiSprite) {
	hub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		switch ev.ButtonText {
		case "Mode":
			uis.SwitchGameMode()
		}
	})
	hub.Subscribe(events.UISpriteAction{}, func(e tasks.Event) {
		ev := e.(events.UISpriteAction)
		if ev.UiSprite == uis.Label {
			switch ev.UiSpriteAction {
			case "highlight":
				FlipHighlight(uis)
			}
		}
	})
	uis.EventHub.Subscribe(events.FishTankLayout{}, func(e tasks.Event) {
		ev := e.(events.FishTankLayout)
		x := ev.Rectangle.Min.X
		y := ev.Rectangle.Min.Y
		switch uis.Label {
		case string(Phreader):
			uis.activationRect = image.Rect(x, y+20, ev.Rectangle.Dx()+210, 100)
		}
	})
	hub.Subscribe(input.CursorPressed{}, func(e tasks.Event) {
		if uis.SpriteHovered() {
			uis.DoClick()
		}
	})
}

func (us *UiSprite) DoClick() {

}

func (us *UiSprite) SwitchGameMode() {
	switch us.gameMode {
	case Normal:
		us.gameMode = Position
	case Position:
		us.gameMode = Normal
	}
}

func (us *UiSprite) SavePosition() drawables.SavePositionData {

	sp := drawables.SavePositionData{}
	sp.X = us.X
	sp.Y = us.Y
	sp.Name = us.Label
	return sp
}

func (us *UiSprite) PublishPickedUpEventIfClicked() {
	if us.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ev := events.UISpriteAction{
			UiSprite:       us.Label,
			UiSpriteAction: "picked up",
		}
		us.EventHub.Publish(ev)
	}

}

func initClickMeEffect(us *UiSprite) {
	cs := ebiten.ColorScale{}
	cs.SetR(0.1)
	cs.SetB(0.2)
	cs.SetG(1.0)
	cs.SetA(1.0)
	msg := "Click Me"
	us.graphicPublishedID = graphics.NewGraphicText(&msg, 24, float64(us.X), float64(us.Y), true, cs, float64(us.Img.Bounds().Dx()), true)

	us.highlight = true
	ols := shaders.LoadOutlineShader()

	us.Shader = ols
	us.ShaderParams["Opacity"] = float32(0.0)
	us.ShaderParams["OutlineColor"] = [4]float32{0.2, 0.7, 0.2, 1.0}
	us.UpdateShaderParams = shaders.UpdatePulseWithText
}

func turnOffClickMeEffect(us *UiSprite) {
	us.Shader = nil
	if us.timers != nil {
		us.timers["clickMeBuffer"].TurnOn()
	}
	graphics.DeInitGraphicId(us.graphicPublishedID)
}

func FlipHighlight(us *UiSprite) {
	switch us.highlight {
	case true:
		us.highlight = false
	case false:
		us.highlight = true
	}
}
