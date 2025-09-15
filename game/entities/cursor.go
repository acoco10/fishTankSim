package entities

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"log"
)

type CursorState uint8

const (
	idle CursorState = iota
	idlePressed
	active
	activePressed
	menu
)

type CursorUpdater struct {
	focusEntityID     uint32
	currentPosition   image.Point
	systemPosition    image.Point
	state             CursorState
	speed             float64 //multiply by current position to reduce or increase spead of cursor
	resistanceX       float64
	resistanceY       float64
	bounds            *image.Rectangle
	stateBeforeWindow CursorState
	cursorImages      map[CursorState]*ebiten.Image
}

func (cu *CursorUpdater) AfterUpdate() {
}

func (cu *CursorUpdater) Update() {
	x, y := ebiten.CursorPosition()

	mouseX, mouseY := ebiten.CursorPosition()
	actualPosition := image.Point{mouseX, mouseY}

	// Apply resistance/smoothing to create drag feeling
	if cu.speed != 0.0 {

		// Calculate the difference between actual mouse and current virtual cursor
		deltaX := float64(mouseX) - float64(cu.currentPosition.X)
		deltaY := float64(mouseY) - float64(cu.currentPosition.Y)

		// Apply resistance: only move a fraction of the way toward the target
		// Lower speed values = more resistance (slower movement)
		cu.currentPosition.X += int(deltaX * cu.speed)
		cu.currentPosition.Y += int(deltaY * cu.speed)

		// Alternative: Use exponential smoothing for more natural feel
		// smoothing := 0.1 * cu.speed  // Adjust this value for different resistance
		// cu.currentPosition.X = int(float64(cu.currentPosition.X)*(1-smoothing) + float64(mouseX)*smoothing)
		// cu.currentPosition.Y = int(float64(cu.currentPosition.Y)*(1-smoothing) + float64(mouseY)*smoothing)
	} else {
		// No resistance - cursor follows mouse exactly
		cu.currentPosition = actualPosition
	}

	x, y = util.GetScaledCursorPosition()
	pt := image.Point{x, y}

	entity, exists := GetEntity(cu.focusEntityID)
	if !exists {
		cu.state = idle
		return
	}

	if entity.UiData != nil {
		if pt.In(entity.UiData.ActivationRect) {
			cu.state = active
		} else {
			cu.state = idle
		}
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if cu.state == idle {
			cu.state = idlePressed
		}
		if cu.state == active {
			cu.state = activePressed
		}
	} else {
		if cu.state == idlePressed {
			cu.state = idle
		}
	}

	if registry.Config.Zoom {
		cu.state = menu
	} else if cu.state == menu {
		cu.state = idle
	}

	registry.Config.Set(registry.CursorPoint, cu.currentPosition)
	if cu.bounds != nil {
		scaledX, _ := util.GetScaledCursorPosition()
		if scaledX > cu.bounds.Max.X {
			cu.currentPosition.X = cu.bounds.Max.X * registry.Config.ResolutionScalingi
		}
		if scaledX < cu.bounds.Min.X {
			cu.currentPosition.X = cu.bounds.Min.X * registry.Config.ResolutionScalingi
		}
		registry.Config.Set(registry.CursorPoint, cu.currentPosition)
	}
}

func (cu *CursorUpdater) MouseButtonJustReleased(b ebiten.MouseButton) bool {
	return false
}

func (cu *CursorUpdater) Draw(screen *ebiten.Image) {
	dopts := &ebiten.DrawImageOptions{}
	dopts.GeoM.Translate(float64(cu.currentPosition.X), float64(cu.currentPosition.Y))
	screen.DrawImage(cu.cursorImages[cu.state], dopts)
}
func (cu *CursorUpdater) AfterDraw(screen *ebiten.Image) {
}

// MouseButtonPressed returns whether mouse button b is currently pressed.
func (cu *CursorUpdater) MouseButtonPressed(b ebiten.MouseButton) bool {
	return ebiten.IsMouseButtonPressed(b) || ebiten.IsKeyPressed(ebiten.KeyEnter)
}

// MouseButtonJustPressed returns whether mouse button b has just been pressed.
// It only returns true during the first frame that the button is pressed.
func (cu *CursorUpdater) MouseButtonJustPressed(b ebiten.MouseButton) bool {
	return inpututil.IsMouseButtonJustPressed(b) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// CursorPosition returns the current cursor position.
// If you define a CursorPosition that doesn't align with a system cursor you will need to
// set the CursorDrawMode to Custom. This is because ebiten doesn't have a way to set the
// cursor location manually
func (cu *CursorUpdater) CursorPosition() (int, int) {
	return cu.currentPosition.X, cu.currentPosition.Y
}

// GetCursorImage Returns the image to use as the cursor
// EbitenUI by default will look for the following cursors:
//
//	"EWResize"
//	"NSResize"
//	"Default"
func (cu *CursorUpdater) GetCursorImage(name string) *ebiten.Image {
	return cu.cursorImages[cu.state]
}
func (cu *CursorUpdater) AftedUpdate() {
}
func (cu *CursorUpdater) ChangeSpeed(newSpeed float64) {
	registry.Config.Set(registry.CursorSpeed, newSpeed)
	cu.speed = newSpeed
}
func (cu *CursorUpdater) ResetSpeed() {
	registry.Config.Set(registry.CursorSpeed, 0.0)
	cu.speed = 0.0
}

func (cu *CursorUpdater) SetBounds(rectangle image.Rectangle) {
	cu.bounds = &rectangle
	cu.speed = 0.0
}
func (cu *CursorUpdater) ResetBounds() {
	cu.bounds = nil
}

// GetCursorOffset Returns how far from the CursorPosition to offset the cursor image.
// This is best used with cursors such as resizing.
func (cu *CursorUpdater) GetCursorOffset(name string) image.Point {
	return image.Point{}
}

func loadCursorImages() map[CursorState]*ebiten.Image {
	imgMap := make(map[CursorState]*ebiten.Image)
	img, err := util.LoadImageAssetAsEbitenImage("uiSprites/handSheet")
	if err != nil {
		log.Fatal(err)
	}

	imgMap[idle] = ebiten.NewImageFromImage(img.SubImage(image.Rect(0, 0, 32, 32)))
	imgMap[idlePressed] = ebiten.NewImageFromImage(img.SubImage(image.Rect(32, 0, 64, 32)))
	imgMap[active] = ebiten.NewImageFromImage(img.SubImage(image.Rect(64, 0, 96, 32)))
	imgMap[activePressed] = ebiten.NewImageFromImage(img.SubImage(image.Rect(96, 0, 128, 32)))
	imgMap[menu] = ebiten.NewImageFromImage(img.SubImage(image.Rect(128, 0, 160, 32)))
	return imgMap
}

func CreateCursorUpdater(hub *tasks.EventHub) *CursorUpdater {
	c := &CursorUpdater{}
	c.cursorImages = loadCursorImages()
	subs(hub, c)
	return c
}

func subs(hub *tasks.EventHub, updater *CursorUpdater) {

	hub.Subscribe(events.FocusEvent{}, func(e tasks.Event) {
		ev := e.(events.FocusEvent)
		updater.focusEntityID = ev.EntID
	})

	hub.Subscribe(events.UnFocusEvent{}, func(e tasks.Event) {
		updater.state = idle
		updater.focusEntityID = 0
	})

	hub.Subscribe(events.WindowOpened{}, func(e tasks.Event) {
		updater.state = menu
		updater.focusEntityID = 0
	})

	hub.Subscribe(events.WindowClosed{}, func(e tasks.Event) {
		updater.state = idle
		updater.focusEntityID = 0
	})

}
