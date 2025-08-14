package customCursor

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
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
	stateBeforeWindow CursorState
	cursorImages      map[CursorState]*ebiten.Image
}

func (cu *CursorUpdater) Update() {
	x, y := ebiten.CursorPosition()
	cu.currentPosition = image.Point{x, y}

	x, y = util.GetScaledCursorPosition()
	pt := image.Point{x, y}

	entity, exists := entities.GetEntity(cu.focusEntityID)
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

	hub.Subscribe(events.Focus{}, func(e tasks.Event) {
		ev := e.(events.Focus)
		updater.focusEntityID = ev.EntID
	})

	hub.Subscribe(events.UnFocus{}, func(e tasks.Event) {
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
