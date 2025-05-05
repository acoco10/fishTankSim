package debug

import (
	"encoding/json"
	"fishTankWebGame/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"os"
)

type RectState uint8

const (
	Initiated RectState = iota
	Drawn
	Off
	On
)

type Rect struct {
	X1, Y1, X2, Y2 float32
	RectState
	tag string
}

func (r *Rect) Init(tag string) {
	x, y := ebiten.CursorPosition()
	r.X1 = float32(x)
	r.X2 = float32(x)
	r.Y1 = float32(y)
	r.Y2 = float32(y)
	r.tag = tag
	r.RectState = On
}

func (r *Rect) Draw(screen *ebiten.Image) {
	if r.RectState == Initiated || r.RectState == Drawn {
		clr := color.RGBA{R: 200, G: 100, B: 100, A: 255}
		vector.StrokeRect(screen, r.X1, r.Y1, r.X2-r.X1, r.Y2-r.Y1, 2, clr, false)
	}
}

func (r *Rect) Update() {
	if r.RectState == On {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			r.RectState = Initiated
		}
	}
	if r.RectState == Initiated {
		x, y := ebiten.CursorPosition()
		r.X2 = float32(x)
		r.Y2 = float32(y)
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			r.RectState = Drawn
		}
	}
	if r.RectState == Drawn {
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			r.Save()
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			r.RectState = Off
		}
	}
}

func (r *Rect) Save() error {

	existingPos, err := LoadCollisions()
	if err != nil {
		return err
	}

	datMap := make(map[string]Rect)
	datMap[r.tag] = *r

	for key, value := range existingPos {
		datMap[key] = value
	}

	outputSave, err := json.Marshal(datMap)
	if err != nil {
		return err
	}

	println(
		outputSave)

	err = os.WriteFile("../assets/data/collisionPosition.json", outputSave, 999)
	if err != nil {
		return err
	}
	return nil
}

func LoadCollisions() (map[string]Rect, error) {
	colDat, err := assets.DataDir.ReadFile("data/collisionPosition.json")
	if err != nil {
		return map[string]Rect{}, err
	}

	datMap := make(map[string]Rect)

	err = json.Unmarshal(colDat, &datMap)

	if err != nil {
		return map[string]Rect{}, err
	}

	return datMap, nil
}
