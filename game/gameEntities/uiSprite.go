package gameEntities

import (
	"encoding/json"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
	"os"
)

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

type UpdateSprite interface {
	Draw(screen *ebiten.Image)
	Update()
}

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
	label string
}

func (us *UiSprite) Draw(screen *ebiten.Image) {

	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(us.X), float64(us.Y))
	screen.DrawImage(us.Img, &opts)

	if us.state == Selected || us.state == Hovered {
		opts.GeoM.Translate(float64(us.AltOffsetX), float64(us.AltOffsetY))
		screen.DrawImage(us.HoverImg, &opts)
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
	if us.SpriteHovered() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && us.state == Selected {
		us.XYUpdater = NewUpdater(us.Sprite)
	}
	if us.XYUpdater != nil {
		us.XYUpdater.Update()
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && us.state == Clicked && us.stateWas == Clicked {
		us.XYUpdater = nil
		us.savePosition()
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
			img := us.Sprite.Img
			us.Sprite.Img = us.AltImg
			us.AltImg = img
		}
		ev := MouseButtonPressed{
			Point: &Point{X: float32(x), Y: float32(y), PType: Food},
		}
		us.EventHub.Publish(ev)
	}

	if us.stateWas == Clicked && us.state != Clicked {
		img := us.Img
		us.Sprite.Img = us.AltImg
		us.AltImg = img

	}

	us.stateWas = us.state
}

func (us *UiSprite) updateState() {

	if us.SpriteHovered() && (us.state != Clicked && us.state != Selected) {
		us.state = Hovered
	}

	if !us.SpriteHovered() {
		us.state = Idle
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

func NewUiSprite(imgs []*ebiten.Image, hub *EventHub, x, y float32, label string) *UiSprite {
	uis := UiSprite{Sprite: &Sprite{X: x, Y: y}}
	uis.label = label
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

	uis.state = Idle
	uis.gameMode = Normal

	hub.Subscribe(ButtonClickedEvent{}, func(e Event) {
		ev := e.(ButtonClickedEvent)
		switch ev.ButtonText {
		case "Mode":
			uis.SwitchGameMode()
		}
	})

	return &uis
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

func (us *UiSprite) savePosition() {
	sp := SavePositionData{}
	sp.X = us.X
	sp.Y = us.Y
	sp.name = us.label

	outputSave, err := json.Marshal(sp)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("../assets/data/spritePosition.json", outputSave, 999)
	if err != nil {
		log.Fatal(err)
	}

}
