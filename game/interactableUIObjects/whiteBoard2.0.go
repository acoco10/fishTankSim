package interactableUIObjects

import (
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/textEffects"
	"github.com/hajimehoshi/ebiten/v2"
)

var ColorScaleSlice = []ebiten.ColorScale{
	ebiten.ColorScale{},
	ebiten.ColorScale{},
	ebiten.ColorScale{},
}

type WhiteBoardState uint8

const (
	Drawing WhiteBoardState = iota
	NotDrawing
	Clear
)

type WhiteBoardSprite2 struct {
	*UiSprite
	state            WhiteBoardState
	dstImg           *sprite.Sprite
	textBeingWritten *textEffects.TextWithShader
	taskIndex        int
	colorCounter     []ebiten.ColorScale
}

func (w *WhiteBoardSprite2) Init() {
	dstImg := ebiten.NewImage(w.Img.Bounds().Dx(), w.Img.Bounds().Dy())
	w.dstImg = &sprite.Sprite{}
	w.dstImg.Img = dstImg
	w.dstImg.ShaderParams = make(map[string]any)
	w.taskIndex = 1
	w.state = NotDrawing
	w.ShaderParams = make(map[string]any)
	w.Shader = registry.ShaderMap["Outline"]
	subscribe(w)

	ColorScaleSlice[0].Scale(1, 0, 0, 1)
	ColorScaleSlice[1].Scale(0, 1, 0, 1)
	ColorScaleSlice[1].Scale(0, 0, 1, 1)

}

func (w *WhiteBoardSprite2) Update() {

	w.Sprite.Update()

	switch w.state {
	case Drawing:
		w.ShaderParams["OutlineColor"] = []float64{0, 0, 1, 1}
	case NotDrawing:
		w.ShaderParams["OutlineColor"] = []float64{0, 1, 0, 1}
	case Clear:
		w.ShaderParams["OutlineColor"] = []float64{1, 0, 0, 1}
	}

	if w.textBeingWritten == nil {
		return
	}

	w.dstImg.Update()

	if w.textBeingWritten.IsFullyDrawn() {
		w.state = NotDrawing
		return
	}

	w.textBeingWritten.Update()

}

func (w *WhiteBoardSprite2) Draw(screen *ebiten.Image) {

	if w.state == Drawing {
		w.textBeingWritten.Draw(w.dstImg.Img)
		w.dstImg.Draw(w.Img)
		w.dstImg.Img.Clear()
	}

	w.Sprite.Draw(screen)
}

func subscribe(w *WhiteBoardSprite2) {
	w.EventHub.Subscribe(tasks.TaskCreated{}, func(e tasks.Event) {
		println("Recieving Event")
		ev := e.(tasks.TaskCreated)
		task := ev.Task
		w.updateTextBeingDrawn(task.Text)
	})
}

func (w *WhiteBoardSprite2) updateTextBeingDrawn(text string) {

	yInset := float64(w.taskIndex * 10)
	insets := [2]float64{10, yInset}

	txt := textEffects.NewTextWithMarkerShader(text, w.dstImg.Img.Bounds(), insets, ColorScaleSlice[1])
	w.textBeingWritten = txt
	w.state = Drawing
	w.taskIndex++

}
