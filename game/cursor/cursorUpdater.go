package cursorUpdater

import (
	"fishTankWebGame/assets"
	"github.com/ebitenui/ebitenui/input"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	"io/fs"
	"log"
)

type CursorUpdater struct {
	currentPosition image.Point
	systemPosition  image.Point
	statusX         int
	statusY         int
	cursorImages    map[string]*ebiten.Image
	counter         int
	countdown       int
}

func CreateCursorUpdater() *CursorUpdater {
	cu := CursorUpdater{}
	X, Y := ebiten.CursorPosition()
	cu.currentPosition = image.Point{X, Y}
	cu.cursorImages = make(map[string]*ebiten.Image)
	cu.cursorImages[input.CURSOR_DEFAULT] = loadNormalCursorImage()
	cu.cursorImages["pressed"] = loadPressedCursorImage()
	cu.cursorImages["statusBar"] = loadNormalCursorImage()
	cu.countdown = 0
	return &cu
}

// Called every Update call from Ebiten
// Note that before this is called the current cursor shape is reset to DEFAULT every cycle

func (cu *CursorUpdater) Update() {

	if cu.countdown > 0 {
		cu.countdown--
	}

	X, Y := ebiten.CursorPosition()

	diffX := cu.systemPosition.X - X
	diffY := cu.systemPosition.Y - Y

	cu.currentPosition.X -= diffX
	cu.currentPosition.Y -= diffY

	/*if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		cu.currentPosition.X -= 10
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		cu.currentPosition.X += 10
	*/

	cu.systemPosition = image.Point{X, Y}

}
func (cu *CursorUpdater) Draw(screen *ebiten.Image) {
}
func (cu *CursorUpdater) AfterDraw(screen *ebiten.Image) {
}

// MouseButtonPressed returns whether mouse button b is currently pressed.
func (cu *CursorUpdater) MouseButtonPressed(b ebiten.MouseButton) bool {
	input.SetCursorImage("clicked", cu.cursorImages["pressed"])
	return ebiten.IsMouseButtonPressed(b)
}

func (cu *CursorUpdater) MouseButtonJustReleased(b ebiten.MouseButton) bool {

	return inpututil.IsMouseButtonJustReleased(b)
}

// MouseButtonJustPressed returns whether mouse button b has just been pressed.
// It only returns true during the first frame that the button is pressed.
func (cu *CursorUpdater) MouseButtonJustPressed(b ebiten.MouseButton) bool {
	return inpututil.IsMouseButtonJustPressed(b)
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
	return cu.cursorImages[name]
}

// GetCursorOffset Returns how far from the CursorPosition to offset the cursor image.
// This is best used with cursors such as resizing.
func (cu *CursorUpdater) GetCursorOffset(name string) image.Point {
	return image.Point{}
}

// Layout implements gameScenes.

func loadNormalCursorImage() *ebiten.Image {
	f, err := assets.ImagesDir.Open("images/fishFoodCursor.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(f fs.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	i, _, _ := ebitenutil.NewImageFromReader(f)
	return ebiten.NewImageFromImage(i.SubImage(image.Rect(0, 0, 171, 135)))
	//(64, 0, 87, 16)
}

func loadHoverCursorImage() *ebiten.Image {
	f, err := assets.ImagesDir.Open("images/fishFoodCursor.png")
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	i, _, _ := ebitenutil.NewImageFromReader(f)
	return ebiten.NewImageFromImage(i.SubImage(image.Rect(0, 0, 171, 135)))
}

func loadPressedCursorImage() *ebiten.Image {
	f, err := assets.ImagesDir.Open("images/fishFoodCursorClicked.png")
	if err != nil {
		return nil
	}
	defer f.Close()
	i, _, _ := ebitenutil.NewImageFromReader(f)
	return i
}
