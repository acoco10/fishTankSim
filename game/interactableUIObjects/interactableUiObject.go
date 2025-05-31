package interactableUIObjects

import (
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/geometry"
	"github.com/acoco10/fishTankWebGame/game/input"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
)

type UISpriteLabel string

const (
	Records    UISpriteLabel = "records"
	FishFood   UISpriteLabel = "fishFood"
	FishBook   UISpriteLabel = "book"
	WhiteBoard UISpriteLabel = "whiteBoard"
	Plant      UISpriteLabel = "plant"
)

type uiSpriteState uint8

const (
	Selected uiSpriteState = iota
	HoveredOver
	ClickedWhileBeingSelected
	Idle
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
	*events.EventHub
	state    uiSpriteState
	stateWas uiSpriteState
	gameMode
	clicked                   bool
	Label                     string
	screenHeight, screenWidth int
}

func (us *UiSprite) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(us.X), float64(us.Y))
	screen.DrawImage(us.Img, &opts)

	if us.Shader != nil {
		shaderOpts := &ebiten.DrawRectShaderOptions{}
		shaderOpts.Uniforms = us.ShaderParams
		shaderOpts.GeoM.Translate(float64(us.X), float64(us.Y))
		shaderOpts.Images[0] = us.Img
		b := us.Img.Bounds().Max
		screen.DrawRectShader(b.X, b.Y, us.Shader, shaderOpts)
		return
	}

	if us.state == Selected || us.state == HoveredOver {
		if us.HoverImg != nil {
			opts.GeoM.Translate(float64(us.AltOffsetX), float64(us.AltOffsetY))
			screen.DrawImage(us.HoverImg, &opts)
		}
	}

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

	if us.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && us.state == Selected {

		if us.XYUpdater == nil {
			ev := sprite.UISpriteAction{}
			ev.UiSprite = us.Label
			ev.UiSpriteAction = "picked up"
			us.EventHub.Publish(ev)
		}
		us.XYUpdater = sprite.NewUpdater(us.Sprite)
	}

	if us.XYUpdater != nil {
		us.XYUpdater.Update()
	}

	x, y := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && us.state == ClickedWhileBeingSelected && !us.clicked {
		us.clicked = true
		ev := input.MouseButtonPressed{
			Point: &geometry.Point{X: float32(x), Y: float32(y), PType: geometry.Food},
		}
		us.EventHub.Publish(ev)
	}

	baseDis := math.Hypot(float64(us.X-us.baseX), float64(us.Y-us.baseY))

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if us.state == Selected && baseDis < 100 && us.stateWas == Selected {
			us.returnToBase()
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
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
	ev := sprite.UISpriteAction{}
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
		if us.Y < 200 {
			us.state = ClickedWhileBeingSelected
		}
	}

	if us.state == ClickedWhileBeingSelected && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		us.state = Selected
	}
}

func NewUiSprite(imgs []*ebiten.Image, hub *events.EventHub, x, y float32, label string, screenWidth, screenHeight int) *UiSprite {

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

/*func Subs(hub *events.EventHub, uis UiSprite) {
	hub.Subscribe(ui.ButtonClickedEvent{}, func(e events.Event) {
		ev := e.(ui.ButtonClickedEvent)
		switch ev.ButtonText {
		case "Mode":
			uis.SwitchGameMode()
		}
	})

}*/

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
