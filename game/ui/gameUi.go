package ui

import (
	"bytes"
	"fishTankWebGame/assets"
	"fishTankWebGame/game/gameEntities"
	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
)

func LoadMenu(gameWidth, gameHeight int, eHub *gameEntities.EventHub) *ebitenui.UI {
	/*img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/statsMenuBackground.png")
	if err != nil {
		return nil
	}
	*/
	//nineSliceImage := eimage.NewNineSlice(img, [3]int{8, 66 - 16, 8}, [3]int{8, 32 - 16, 8})
	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		//widget.ContainerOpts.BackgroundImage(nineSliceImage),
		// the container will use a plain color as its background
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(20),
			widget.RowLayoutOpts.Padding(
				widget.Insets{Right: 0, Left: 0, Top: 200, Bottom: 20}),
		),
		),
	)

	button := LoadButton("Save", eHub)
	modeButton := LoadButton("Mode", eHub)

	rootContainer.AddChild(button)
	rootContainer.AddChild(modeButton)
	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	return &ui
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := eimage.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 0})

	hover := eimage.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := eimage.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

func LoadFont(size float64) (text.Face, error) {
	loadedFont, err := assets.FontsDir.ReadFile("fonts/nk57.otf")
	if err != nil {
		return nil, err
	}
	s, err := text.NewGoTextFaceSource(bytes.NewReader(loadedFont))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}

func LoadButton(buttonText string, hub *gameEntities.EventHub) *widget.Button {
	//load a generic button labeled with button text string that will send a button clicked event to event hub
	buttonImage, err := loadButtonImage()
	if err != nil {
		log.Fatal()
	}

	face, err := LoadFont(12)
	if err != nil {
		log.Fatal()
	}

	var button *widget.Button

	button = widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
		widget.ButtonOpts.Text(buttonText, face, &widget.ButtonTextColor{
			Idle:    color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
			Hover:   color.NRGBA{0, 255, 128, 255},
			Pressed: color.NRGBA{255, 0, 0, 255},
		}),
		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),
		//Move the text down and right on press
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			button.Text().Inset.Top = 4
			button.Text().Inset.Left = 4
			button.GetWidget().CustomData = true
		}),
		//Move the text back to start on press released
		widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {
			button.Text().Inset.Top = 0
			button.Text().Inset.Left = 0
			button.GetWidget().CustomData = false
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {

			ev := gameEntities.ButtonClickedEvent{
				buttonText,
			}

			hub.Publish(ev)
		}),

		// add a handler that reacts to entering the button with the cursor
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
			//If we moved the Text because we clicked on this button previously, move the text down and right
			if button.GetWidget().CustomData == true {
				button.Text().Inset.Top = 4
				button.Text().Inset.Left = 4
			}
		}),

		// add a handler that reacts to moving the cursor on the button
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) {
		}),

		// add a handler that reacts to exiting the button with the cursor
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) {
			//Reset the Text inset if the cursor is no longer over the button
			button.Text().Inset.Top = 0
			button.Text().Inset.Left = 0
		}),
	)
	return button
}
