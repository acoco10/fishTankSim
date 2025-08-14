package util

import (
	"encoding/json"
	"fmt"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
	"image"
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
	*image.Rectangle
	RectState
	tag                string
	cursorMarkerRadius float32
}

func (r *Rect) Init(tag string) {
	x, y := GetScaledCursorPosition()
	r.Rectangle.Min = image.Point{X: int(x), Y: int(y)}
	r.Rectangle.Min = image.Point{X: int(x), Y: int(y)}
	r.tag = tag
	r.RectState = On
}

func (r *Rect) Draw(screen *ebiten.Image) {
	if r.RectState == Initiated || r.RectState == Drawn {
		StrokeRectFromImageRect(*r.Rectangle, screen, colornames.Darkmagenta)
	}
	if r.RectState == On {
		StrokeRectFromImageRect(*r.Rectangle, screen, colornames.Orangered)
	}
}

func (r *Rect) Update() error {

	r.cursorMarkerRadius++

	if r.cursorMarkerRadius == 3 {
		r.cursorMarkerRadius = 0
	}

	if r.RectState == On {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			r.RectState = Initiated
		}
	}
	if r.RectState == Initiated {
		x, y := GetScaledCursorPosition()
		r.Max = image.Point{X: int(x), Y: int(y)}
		if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			r.RectState = Drawn
		}
	}
	if r.RectState == Drawn {
		if inpututil.IsKeyJustPressed(ebiten.KeyS) {
			err := r.Save()
			if err != nil {
				return err
			}
			fmt.Printf("Successfully saved rectangle = %s\n", r.tag)
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			r.RectState = Off
		}
	}
	if r.RectState == Off && inpututil.IsKeyJustPressed(ebiten.KeyN) {
		r.RectState = On
	}
	return nil
}

func (r *Rect) Save() error {

	existingPos, err := LoadCollisions()
	if err != nil {
		return err
	}

	datMap := make(map[string]image.Rectangle)

	datMap[r.tag] = *r.Rectangle

	for key, value := range existingPos {
		if key != r.tag {
			datMap[key] = value
		}
	}

	outputSave, err := json.Marshal(datMap)
	if err != nil {
		return err
	}

	println(
		outputSave)

	err = os.WriteFile("assets/data/collisionPosition.json", outputSave, 999)
	if err != nil {
		return err
	}
	return nil
}

func LoadCollisions() (map[string]image.Rectangle, error) {
	colDat, err := assets.DataDir.ReadFile("data/collisionPosition.json")
	if err != nil {
		return map[string]image.Rectangle{}, err
	}

	datMap := make(map[string]image.Rectangle)

	err = json.Unmarshal(colDat, &datMap)

	if err != nil {
		return map[string]image.Rectangle{}, err
	}

	return datMap, nil
}
