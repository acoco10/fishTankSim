package interactableUIObjects

import (
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
)

type uiSpriteState uint8

const (
	Selected uiSpriteState = iota
	HoveredOver
	ClickedWhileBeingSelected
	Idle
	Clickable
)

type gameMode uint8

const (
	Position gameMode = iota
	Normal
)

type UiSprite struct {
	*sprite.Sprite
	baseX, baseY           float32
	HoverImg               *ebiten.Image
	AltImg                 *ebiten.Image
	AltOffsetX, AltOffsetY float32
	*sprite.XYUpdater
	*tasks.EventHub
	state    uiSpriteState
	stateWas uiSpriteState
	gameMode
	clicked                   bool
	Label                     string
	screenHeight, screenWidth int
}

func (us *UiSprite) Draw(screen *ebiten.Image) {
	us.Sprite.Draw(screen)
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
	us.updateState()
	if us.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		us.XYUpdater = sprite.NewUpdater(us.Sprite)
	}

	if us.XYUpdater != nil {
		us.XYUpdater.Update()
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && us.state == ClickedWhileBeingSelected && us.stateWas == ClickedWhileBeingSelected {
		us.XYUpdater = nil
	}

	us.stateWas = us.state
}

func (us *UiSprite) UpdateNormal() {
	us.clicked = false

	us.updateState()

	if us.state == HoveredOver && us.stateWas != HoveredOver {

		ols := shaders.LoadOutlineShader()
		us.Shader = ols
		us.ShaderParams = make(map[string]any)
		us.ShaderParams["OutlineColor"] = [4]float64{1, 1, 0, 1}

	}

	if us.Shader != nil && (us.state != HoveredOver && us.state != Selected) {

		us.Shader = nil

	}

	if us.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && us.state == Selected {

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

	//x, y := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && us.state == ClickedWhileBeingSelected && !us.clicked {

		us.clicked = true

		/*ev := input.MouseButtonPressedUISpriteActivity{
			//filler as currently only fish food needs to generate points
			Point: &geometry.Point{X: float32(x), Y: float32(y), PType: geometry.Structure},
		}*/

		//us.EventHub.Publish(ev)

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
	us.X = us.baseX
	us.Y = us.baseY
	ev := events.UISpriteAction{}
	ev.UiSprite = us.Label
	ev.UiSpriteAction = "put back"
	us.EventHub.Publish(ev)
	us.XYUpdater = nil
}

func (us *UiSprite) updateState() {
	if !us.SpriteHovered() {
		us.state = Idle
	}

	if us.SpriteHovered() && (us.state != ClickedWhileBeingSelected && us.state != Selected) {
		us.state = HoveredOver
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && us.state == HoveredOver {
		us.state = Selected
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && us.state == Selected && us.stateWas == Selected {
		us.state = ClickedWhileBeingSelected
	}

	if us.state == ClickedWhileBeingSelected && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		us.state = Selected
	}
}

func NewUiSprite(imgs []*ebiten.Image, hub *tasks.EventHub, x, y float32, label string, screenWidth, screenHeight int) *UiSprite {

	var paramaMappa = make(map[string]any)

	uis := UiSprite{Sprite: &sprite.Sprite{X: x, Y: y}}
	uis.ShaderParams = paramaMappa
	uis.baseX = x
	uis.baseY = y
	uis.Label = label
	uis.EventHub = hub

	uis.Img = &ebiten.Image{}
	uis.Img = imgs[0]

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

	//Subs(hub, uis)

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

func Subs(hub *tasks.EventHub, uis UiSprite) {
	hub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		switch ev.ButtonText {
		case "Mode":
			uis.SwitchGameMode()
		}
	})
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

func initClickMeEffect(us *UiSprite) {
	us.EventHub.Publish(
		events.ClickMeGraphicEvent{
			X: float64(us.X), Y: float64(us.Y), SpriteWidth: float64(us.Img.Bounds().Dx())},
	)

	ols := shaders.LoadOutlineShader()

	us.Shader = ols
	us.ShaderParams["Opacity"] = float32(0.0)
	us.ShaderParams["OutlineColor"] = [4]float32{0.2, 0.7, 0.2, 1.0}
	us.UpdateShaderParams = shaders.UpdatePulseWithText
}

func turnOffClickMeEffect(us *UiSprite) {
	us.Shader = nil
	ev := events.TurnOffGraphic{X: float64(us.X), Y: float64(us.Y)}
	us.EventHub.Publish(ev)
}
