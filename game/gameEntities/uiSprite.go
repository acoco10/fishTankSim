package gameEntities

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type DrawableSprite interface {
	Draw(screen *ebiten.Image)
	Update()
}

type uiSpriteState uint8

const (
	Selected uiSpriteState = iota
	Hovered
	Clicked
	Idle
)

type gameMode uint8

const (
	Position gameMode = iota
	Normal
)

type UiSprite struct {
	*Sprite
	HoverImg               *ebiten.Image
	AltImg                 *ebiten.Image
	AltOffsetX, AltOffsetY float32
	*XYUpdater
	*EventHub
	state    uiSpriteState
	stateWas uiSpriteState
	gameMode
	Label                     string
	screenHeight, screenWidth int
}

func (us *UiSprite) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(us.X), float64(us.Y))
	screen.DrawImage(us.Img, &opts)
	if us.shader != nil {
		shaderOpts := &ebiten.DrawRectShaderOptions{}
		shaderOpts.Uniforms = us.shaderParams
		shaderOpts.GeoM.Translate(float64(us.X), float64(us.Y))
		shaderOpts.Images[0] = us.Img
		b := us.Img.Bounds().Max
		screen.DrawRectShader(b.X, b.Y, us.shader, shaderOpts)
		return
	}
	if us.state == Selected || us.state == Hovered {
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
		us.XYUpdater = NewUpdater(us.Sprite)
	}
	if us.XYUpdater != nil {
		us.XYUpdater.Update()
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && us.state == Clicked && us.stateWas == Clicked {
		us.XYUpdater = nil
	}

	us.stateWas = us.state

}
func (us *UiSprite) UpdateNormal() {
	us.updateState()
	if us.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && us.state == Selected {
		us.XYUpdater = NewUpdater(us.Sprite)
	}

	if us.XYUpdater != nil {
		us.XYUpdater.Update()
	}

	x, y := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && us.state == Clicked {
		if us.AltImg != nil && us.stateWas != Clicked {
			if us.AltImg != nil {
				img := us.Sprite.Img
				us.Sprite.Img = us.AltImg
				us.AltImg = img
			}
		}
		ev := MouseButtonPressed{
			Point: &Point{X: float32(x), Y: float32(y), PType: Food},
		}
		us.EventHub.Publish(ev)
	}

	if us.stateWas == Clicked && us.state != Clicked {
		if us.AltImg != nil {
			img := us.Img
			us.Sprite.Img = us.AltImg
			us.AltImg = img
		}

	}

	us.stateWas = us.state
}

func (us *UiSprite) updateState() {
	if !us.SpriteHovered() {
		us.state = Idle
	}

	if us.SpriteHovered() && (us.state != Clicked && us.state != Selected) {
		us.state = Hovered
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && us.state == Hovered {
		us.state = Selected
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && us.state == Selected && us.stateWas == Selected {
		us.state = Clicked
	}

	if us.state == Clicked && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		us.state = Selected
	}
}

func NewUiSprite(imgs []*ebiten.Image, hub *EventHub, x, y float32, label string, screenWidth, screenHeight int) *UiSprite {
	uis := UiSprite{Sprite: &Sprite{X: x, Y: y}}
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

	UiSpriteEventSub(hub, uis)

	uis.state = Idle
	uis.gameMode = Normal

	return &uis
}

func UiSpriteEventSub(hub *EventHub, uis UiSprite) {
	hub.Subscribe(ButtonClickedEvent{}, func(e Event) {
		ev := e.(ButtonClickedEvent)
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

type SavePositionData struct {
	X    float32
	Y    float32
	name string
}

func (us *UiSprite) SavePosition() SavePositionData {
	sp := SavePositionData{}
	sp.X = us.X
	sp.Y = us.Y
	sp.name = us.Label
	return sp
}

type FishFoodSprite struct {
	*UiSprite
}

func (ff *FishFoodSprite) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	if ff.state == Idle || ff.state == Clicked {
		opts.GeoM.Translate(float64(ff.X), float64(ff.Y))
		screen.DrawImage(ff.Img, &opts)
		opts.GeoM.Reset()
	} else if ff.state == Selected || ff.state == Hovered {
		if ff.HoverImg != nil {
			opts.GeoM.Translate(float64(ff.X), float64(ff.Y))
			opts.GeoM.Translate(float64(ff.AltOffsetX), float64(ff.AltOffsetY))
			screen.DrawImage(ff.HoverImg, &opts)
		}
	}
}
