package util

import (
	"encoding/json"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
	"image"
	"log"
	"os"
)

type RectState uint8

const (
	Initiated RectState = iota
	Drawn
	On
)

type Rect struct {
	*image.Rectangle
	RectState
	tag                string
	cursorMarkerRadius float32
	interPoints        []image.Point
	StateMachine
	eventHub    *tasks.EventHub
	drawMeDaddy bool
}

type StateMachine struct {
	States       map[RectState]StateHandler
	CurrentState RectState
}

func (s *StateMachine) Transition(rect *Rect) {
	if s.States[s.CurrentState].TransitionFunc != nil {
		s.States[s.CurrentState].TransitionFunc(rect)
	}
	s.CurrentState = s.States[s.CurrentState].TransitionTo
}

type StateHandler struct {
	Updater        func(entity *Rect)
	TransitionTo   RectState
	TransitionFunc func(rect *Rect)
}

func (r *Rect) GivePoint(pt image.Point) {
	r.interPoints = append(r.interPoints, pt)
}

func (r *Rect) Init(tag string, hub *tasks.EventHub) {

	state1 := StateHandler{Updater: updateRectInitState, TransitionTo: Initiated, TransitionFunc: transitionFromIntToDraw}
	state2 := StateHandler{Updater: updateRectInitiated, TransitionTo: Drawn}
	state3 := StateHandler{Updater: updateRectDrawn, TransitionTo: On}

	states := map[RectState]StateHandler{
		On:        state1,
		Initiated: state2,
		Drawn:     state3,
	}

	r.tag = tag
	sm := StateMachine{States: states}
	r.eventHub = hub
	r.StateMachine = sm
	r.CurrentState = On
	r.subscribe()

}

func updateRectInitState(rect *Rect) {
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		rect.eventHub.Publish(events.DebugTextInput{LastText: rect.tag})
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && rect.drawMeDaddy {
		rect.Transition(rect)
	}
}

func transitionFromIntToDraw(rect *Rect) {
	x, y := GetScaledCursorPosition()
	rect.Rectangle.Min = image.Point{X: int(x), Y: int(y)}
	rect.Rectangle.Min = image.Point{X: int(x), Y: int(y)}
	rect.Rectangle.Max = image.Point{X: int(x), Y: int(y)}
	rect.Rectangle.Max = image.Point{X: int(x), Y: int(y)}
}

func updateRectInitiated(rect *Rect) {
	x, y := GetScaledCursorPosition()
	rect.Max = image.Point{X: x, Y: y}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		rect.Transition(rect)
	}
}

func updateRectDrawn(rect *Rect) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		rect.Transition(rect)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		err := rect.Save()
		if err != nil {
			println("cannot save current debug rect")
			log.Fatal(err)
		}
		rect.Transition(rect)
	}
}

func (r *Rect) Draw(screen *ebiten.Image) {
	if r.CurrentState == Initiated || r.CurrentState == Drawn {
		StrokeRectFromImageRect(*r.Rectangle, screen, colornames.Darkmagenta, false)
	}
}

func (r *Rect) Update() error {
	r.StateMachine.States[r.CurrentState].Updater(r)
	return nil
}

func (r *Rect) subscribe() {
	r.eventHub.Subscribe(events.DebugTextEntered{}, func(e tasks.Event) {
		ev := e.(events.DebugTextEntered)
		r.drawMeDaddy = true
		r.tag = ev.InputText
	})
}

func (r *Rect) Save() error {

	closePT := ClosestPoint(image.Point{r.Min.X, r.Min.Y}, r.interPoints)

	//theoretically the collision should always be within the bounds of the sprite
	offSetX := r.Min.X - closePT.X
	offSetY := r.Min.Y - closePT.Y

	width := r.Rectangle.Dx()
	height := r.Rectangle.Dy()

	r.Rectangle.Min.X = offSetX
	r.Rectangle.Max.X = offSetX + width
	r.Rectangle.Min.Y = offSetY
	r.Rectangle.Max.Y = offSetY + height

	//positive distance from the image origin for the collision

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

	err = os.WriteFile("assets/data/collisionPositionTest.json", outputSave, 999)
	if err != nil {
		return err
	}
	return nil
}

func LoadCollisions() (map[string]image.Rectangle, error) {
	colDat, err := assets.DataDir.ReadFile("data/collisionPositionTest.json")
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
