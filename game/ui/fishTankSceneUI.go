package ui

import (
	"bytes"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/loaders"
	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
)

type ButtonType uint8

const (
	SubmitButton ButtonType = iota
	SpriteSelectButton
)

func LoadMainFishMenu(gameWidth, gameHeight int, eHub *events.EventHub) (*ebitenui.UI, *TextBoxUi, error) {

	rootContainer := widget.NewContainer(
		//widget.ContainerOpts.BackgroundImage(nineSliceImage),
		// the container will use a plain color as its background
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	buttonContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionStart,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				}),
		),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Spacing(20),
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(widget.Insets{Right: 0, Left: 50, Top: 100, Bottom: 0}),
			)),
	)

	button := LoadSubmitButton("Save", eHub, 12)
	//modeButton := LoadSubmitButton("Mode", eHub, 12)

	button2 := LoadSubmitButton("Mute Music", eHub, 12)
	button3 := LoadSubmitButton("Mute Sounds", eHub, 12)

	buttonContainer.AddChild(button)
	buttonContainer.AddChild(button2)
	buttonContainer.AddChild(button3)
	//buttonContainer.AddChild(modeButton)

	fishStats, err := NewTextBlock(eHub, StatsMenu)

	if err != nil {
		return nil, nil, err
	}

	fishStats.text.GetWidget().Visibility = widget.Visibility_Hide

	notePad, err := NewTextBlock(eHub, NotePad)
	if err != nil {
		return nil, nil, err
	}

	notePad.text.SetText("To Do:")

	rootContainer.AddChild(fishStats)
	rootContainer.AddChild(notePad)
	rootContainer.AddChild(buttonContainer)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	return &ui, fishStats, nil
}

func loadSubmitButtonImage() (*widget.ButtonImage, error) {

	img, err := loaders.LoadImageAssetAsEbitenImage("menuAssets/submitButton3")

	if err != nil {
		return nil, err
	}

	imgClicked, err := loaders.LoadImageAssetAsEbitenImage("menuAssets/submitButtonAlt")

	if err != nil {
		return nil, err
	}

	nineSliceImage := eimage.NewNineSlice(img, [3]int{9, img.Bounds().Dx() - 18, 9}, [3]int{8, 9, 10})

	nineSliceImageClicked := eimage.NewNineSlice(imgClicked, [3]int{9, img.Bounds().Dx() - 18, 9}, [3]int{8, 9, 10})

	idle := nineSliceImage

	hover := nineSliceImage

	pressed := nineSliceImageClicked

	return &widget.ButtonImage{
		Idle:         idle,
		Hover:        hover,
		Pressed:      pressed,
		PressedHover: pressed,
		Disabled:     eimage.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 200, A: 255}),
	}, nil
}

func loadSpriteSelectButtonImage(t string) (*widget.ButtonImage, error) {

	img, err := loaders.LoadImageAssetAsEbitenImage("menuAssets/spriteOutlineButton")

	if err != nil {
		return nil, err
	}

	imgClicked, err := loaders.LoadImageAssetAsEbitenImage("menuAssets/spriteOutlineButtonAlt")

	if err != nil {
		return nil, err
	}

	nineSliceImage := eimage.NewNineSlice(img, [3]int{16, 32, 16}, [3]int{16, 48, 16})

	nineSliceImageClicked := eimage.NewNineSlice(imgClicked, [3]int{16, 32, 16}, [3]int{16, 48, 16})

	idle := nineSliceImage

	hover := nineSliceImageClicked

	pressed := nineSliceImageClicked
	if t == "Selected Button" {
		return &widget.ButtonImage{
			Idle:    hover,
			Hover:   hover,
			Pressed: pressed,
		}, nil
	} else {
		return &widget.ButtonImage{
			Idle:    idle,
			Hover:   hover,
			Pressed: pressed,
		}, nil
	}
}

func LoadFont(size float64, fontName string) (text.Face, error) {
	var font []byte
	switch fontName {
	case "nk57":
		loadedFont, err := assets.FontsDir.ReadFile("fonts/nk57.otf")
		if err != nil {
			return nil, err
		}
		font = loadedFont
	case "rockSalt":
		loadedFont, err := assets.FontsDir.ReadFile("fonts/RockSalt.ttf")
		if err != nil {
			return nil, err
		}
		font = loadedFont
	}
	s, err := text.NewGoTextFaceSource(bytes.NewReader(font))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}

func LoadMuteButton(buttonText string, hub *events.EventHub, fontSize float64) *widget.Button {
	//load a generic button labeled with button text string that will send a button clicked event to event hub
	buttonImage, err := loadSpriteSelectButtonImage(buttonText)
	if err != nil {
		log.Fatal(err)
	}

	face, err := LoadFont(fontSize, "nk57")
	if err != nil {
		log.Fatal(err)
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
			Idle:    color.NRGBA{0, 0, 0, 0xff},
			Hover:   color.NRGBA{255, 255, 0, 255},
			Pressed: color.NRGBA{255, 255, 0, 255},
		}),
		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    100,
			Bottom: 10,
		}),
		//Move the text down and right on press
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			button.GetWidget().CustomData = true
			button.KeepPressedOnExit = true
		}),
		//Move the text back to start on press released
		widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {
			button.Text().Inset.Top = 0
			button.Text().Inset.Left = 0
			button.GetWidget().CustomData = false
		}),
		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if button.GetWidget().Disabled == false {
				ev := ButtonClickedEvent{
					buttonText,
				}
				hub.Publish(ev)
			}

		}),

		// add a handler that reacts to entering the button with the cursor
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
			//If we moved the Text because we clicked on this button previously, move the text down and right
			if button.GetWidget().Disabled == false {
				ev := ButtonEvent{
					buttonText,
					"cursor entered",
				}
				hub.Publish(ev)
			}

		}),

		// add a handler that reacts to moving the cursor on the button
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) {
		}),

		// add a handler that reacts to exiting the button with the cursor
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) {
			//ResetVls the Text inset if the cursor is no longer over the button
			button.Text().Inset.Top = 0
			button.Text().Inset.Left = 0
			ev := ButtonEvent{
				buttonText,
				"cursor exited",
			}
			hub.Publish(ev)
		}),
	)
	return button
}

func LoadSpriteSelectButton(buttonText string, hub *events.EventHub, fontSize float64) *widget.Button {
	//load a generic button labeled with button text string that will send a button clicked event to event hub
	buttonImage, err := loadSpriteSelectButtonImage(buttonText)
	if err != nil {
		log.Fatal(err)
	}

	face, err := LoadFont(fontSize, "nk57")
	if err != nil {
		log.Fatal(err)
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
			Idle:    color.NRGBA{0, 0, 0, 0xff},
			Hover:   color.NRGBA{255, 255, 0, 255},
			Pressed: color.NRGBA{255, 255, 0, 255},
		}),
		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    100,
			Bottom: 10,
		}),
		//Move the text down and right on press
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			button.GetWidget().CustomData = true
			button.KeepPressedOnExit = true
		}),
		//Move the text back to start on press released
		widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {
			button.Text().Inset.Top = 0
			button.Text().Inset.Left = 0
			button.GetWidget().CustomData = false
		}),
		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if button.GetWidget().Disabled == false {
				ev := ButtonClickedEvent{
					buttonText,
				}
				hub.Publish(ev)
			}

		}),

		// add a handler that reacts to entering the button with the cursor
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
			//If we moved the Text because we clicked on this button previously, move the text down and right
			if button.GetWidget().Disabled == false {
				ev := ButtonEvent{
					buttonText,
					"cursor entered",
				}
				hub.Publish(ev)
			}

		}),

		// add a handler that reacts to moving the cursor on the button
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) {
		}),

		// add a handler that reacts to exiting the button with the cursor
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) {
			//ResetVls the Text inset if the cursor is no longer over the button
			button.Text().Inset.Top = 0
			button.Text().Inset.Left = 0
			ev := ButtonEvent{
				buttonText,
				"cursor exited",
			}
			hub.Publish(ev)
		}),
	)
	return button
}
func LoadSubmitButton(buttonText string, hub *events.EventHub, fontSize float64) *widget.Button {
	//load a generic button labeled with button text string that will send a button clicked event to event hub
	buttonImage, err := loadSubmitButtonImage()
	if err != nil {
		log.Fatal(err)
	}

	face, err := LoadFont(fontSize, "nk57")
	if err != nil {
		log.Fatal(err)
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
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Position: widget.RowLayoutPositionCenter}),
		),
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
		widget.ButtonOpts.Text(buttonText, face, &widget.ButtonTextColor{
			Idle:    color.NRGBA{0, 0, 100, 0xff},
			Hover:   color.NRGBA{0, 0, 100, 255},
			Pressed: color.NRGBA{0, 0, 100, 255},
		}),
		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    3,
			Bottom: 3,
		}),
		//Move the text down and right on press
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {

			println("button event generated for", buttonText)
			button.Text().Inset.Top = 2
			button.Text().Inset.Left = -2
			ev := ButtonClickedEvent{
				ButtonText: buttonText,
			}
			hub.Publish(ev)
		}),
		//Move the text back to start on press released
		widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {
			button.Text().Inset.Top = 0
			button.Text().Inset.Left = 0
			button.GetWidget().CustomData = false
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) {
		}),

		// add a handler that reacts to exiting the button with the cursor
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) {
			//ResetVls the Text inset if the cursor is no longer over the button
			button.Text().Inset.Top = 0
			button.Text().Inset.Left = 0
		}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(113, 27)),
	)
	return button
}

func LoadButton(buttonText string, hub *events.EventHub, fontSize float64) *widget.Button {
	//load a generic button labeled with button text string that will send a button clicked event to event hub
	buttonImage, err := loadSubmitButtonImage()
	if err != nil {
		log.Fatal(err)
	}

	face, err := LoadFont(fontSize, "nk57")
	if err != nil {
		log.Fatal(err)
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
			Idle:    color.NRGBA{0, 0, 0, 0xff},
			Hover:   color.NRGBA{0, 255, 128, 255},
			Pressed: color.NRGBA{255, 0, 0, 255},
		}),
		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    3,
			Bottom: 3,
		}),
		//Move the text down and right on press
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			println("button event generated for", buttonText)

			ev := ButtonClickedEvent{
				ButtonText: buttonText,
			}
			hub.Publish(ev)
		}),
		//Move the text back to start on press released
		widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {
			button.Text().Inset.Top = 0
			button.Text().Inset.Left = 0
			button.GetWidget().CustomData = false
		}),
		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("button event generated for", buttonText)

			ev := ButtonClickedEvent{
				ButtonText: buttonText,
			}

			hub.Publish(ev)

		}),

		// add a handler that reacts to entering the button with the cursor
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
			println("button", buttonText, "is hovered")
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
			//ResetVls the Text inset if the cursor is no longer over the button
			button.Text().Inset.Top = 0
			button.Text().Inset.Left = 0
		}),
	)
	return button
}

func LoadHeader(headerText string, face text.Face) *widget.Container {
	headerContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionStart,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	headerLbl := widget.NewText(
		widget.TextOpts.Text(headerText, face, color.RGBA{R: 0, G: 160, B: 0, A: 255}),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionEnd,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.TextOpts.Insets(widget.Insets{
			Left:   30,
			Right:  10,
			Top:    50,
			Bottom: 600,
		}),
	)

	headerContainer.AddChild(headerLbl)

	return headerContainer
}
