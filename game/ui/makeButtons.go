package ui

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	image "image"
	"image/color"
	"log"
)

func LoadMuteButton(buttonText string, hub *tasks.EventHub, fontSize float64) *widget.Button {
	//load a generic button labeled with button text string that will send a button clicked event to event hub
	buttonImage, err := loadSpriteSelectButtonImage(buttonText)
	if err != nil {
		log.Fatal(err)
	}

	face, err := util.LoadFont(fontSize, "nk57")
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
		widget.ButtonOpts.Text(buttonText, &face, &widget.ButtonTextColor{
			Idle:    color.NRGBA{0, 0, 0, 0xff},
			Hover:   color.NRGBA{255, 255, 0, 255},
			Pressed: color.NRGBA{255, 255, 0, 255},
		}),
		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(&widget.Insets{
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
			button.GetWidget().CustomData = false
		}),
		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if button.GetWidget().Disabled == false {
				ev := events.ButtonClickedEvent{
					buttonText,
				}
				hub.Publish(ev)
			}

		}),

		// add a handler that reacts to entering the button with the cursor
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
			//If we moved the Text because we clicked on this button previously, move the text down and right
			if button.GetWidget().Disabled == false {
				ev := events.ButtonEvent{
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
			ev := events.ButtonEvent{
				buttonText,
				"cursor exited",
			}
			hub.Publish(ev)
		}),
	)
	return button
}

func LoadStackSpriteSelectButtonWithToolTip(buttonText string, fishImg *ebiten.Image, fontSize float64, hub *tasks.EventHub, tipText []string) (*widget.Container, error) {

	var tooltipContainer *widget.Container

	if len(tipText) > 0 {
		tooltipContainer = MakeToolTipContainer(tipText)
	}

	face, err := util.LoadFont(fontSize, "nk57")
	if err != nil {
		log.Fatal(err)
	}

	imgForTransform := ebiten.NewImage(fishImg.Bounds().Dx()*2, fishImg.Bounds().Dx()*2)
	dopts := &ebiten.DrawImageOptions{}
	dopts.GeoM.Scale(2, 2)

	imgForTransform.DrawImage(fishImg, dopts)

	img, err := util.LoadImageAssetAsEbitenImage("menuAssets/spriteOutlineButton")
	if err != nil {
		return nil, err
	}

	imgClicked, err := util.LoadImageAssetAsEbitenImage("menuAssets/spriteOutlineButtonAlt")
	if err != nil {
		return nil, err
	}

	nineSliceImage := eimage.NewNineSlice(img, [3]int{16, 32, 16}, [3]int{16, 48, 16})

	nineSliceImageClicked := eimage.NewNineSlice(imgClicked, [3]int{16, 32, 16}, [3]int{16, 48, 16})

	idle := nineSliceImage

	hover := nineSliceImageClicked

	pressed := nineSliceImageClicked

	btnimg := &widget.ButtonImage{
		Idle:     idle,
		Hover:    hover,
		Pressed:  pressed,
		Disabled: pressed,
	}

	buttonStackedLayout := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),

		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
		})),
	)
	btnIconG := widget.NewGraphic(
		widget.GraphicOpts.Images(&widget.GraphicImage{
			Idle:     imgForTransform,
			Disabled: imgForTransform,
		},
		),
	)
	// construct a pressable button
	var button *widget.Button

	button = widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(btnimg),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			btnIconG.GetWidget().Disabled = !btnIconG.GetWidget().Disabled

		}),
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			button.GetWidget().CustomData = true
			button.KeepPressedOnExit = true
		}),
		//Move the text back to start on press released
		widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {
			button.GetWidget().Disabled = false
			button.GetWidget().CustomData = false
		}),
		widget.ButtonOpts.Text(buttonText, &face, &widget.ButtonTextColor{
			Idle:     color.NRGBA{0, 0, 0, 0xff},
			Hover:    color.NRGBA{255, 255, 0, 255},
			Pressed:  color.NRGBA{255, 255, 0, 255},
			Disabled: color.NRGBA{255, 255, 0, 255},
		}),
		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(&widget.Insets{
			Left:   10,
			Right:  10,
			Top:    100,
			Bottom: 10,
		}),

		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(120, 100),
			widget.WidgetOpts.ToolTip(widget.NewToolTip(
				widget.ToolTipOpts.Content(tooltipContainer),
				//widget.WidgetToolTipOpts.Delay(1*time.Second),
				widget.ToolTipOpts.Offset(image.Point{-5, 5}),
				widget.ToolTipOpts.Position(widget.TOOLTIP_POS_WIDGET),
				//When the Position is set to TOOLTIP_POS_WIDGET, you can configure where it opens with the optional parameters below
				//They will default to what you see below if you do not provide them
				widget.ToolTipOpts.AnchorOriginHorizontal(widget.TOOLTIP_ANCHOR_END),
				widget.ToolTipOpts.AnchorOriginVertical(widget.TOOLTIP_ANCHOR_END),
				widget.ToolTipOpts.ContentOriginHorizontal(widget.TOOLTIP_ANCHOR_END),
				widget.ToolTipOpts.ContentOriginVertical(widget.TOOLTIP_ANCHOR_START),
			))),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if button.GetWidget().Disabled == false {
				button.GetWidget().Disabled = true
				ev := events.ButtonClickedEvent{
					buttonText,
				}
				hub.Publish(ev)
			}
		}),
	)

	buttonStackedLayout.AddChild(button)
	buttonStackedLayout.AddChild(
		btnIconG,
	)

	return buttonStackedLayout, nil
}

func LoadOutlineTextButtonNotDisableOnPress(buttonImg *widget.ButtonImage, buttonText string, hub *tasks.EventHub, fontSize float64) *widget.Container {

	face, err := util.LoadFont(fontSize, "reglisseOutline") //white center text

	face2, err := util.LoadFont(fontSize, "reglisseOutlined") //black outline

	if err != nil {
		log.Fatal("error loading new font")
	}

	buttonStackedLayout := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	// construct a pressable button
	var button *widget.Button
	txt := buttonText

	headerLbl := widget.NewText(
		widget.TextOpts.Text(txt, &face, colornames.Aliceblue),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart))

	headerLblOutline := widget.NewText(
		widget.TextOpts.Text(txt, &face2, colornames.Black),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
	)

	button = widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImg),

		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {

		}),
		//Move the text back to start on press released
		widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {

		}),
		widget.ButtonOpts.Text(txt, &face, &widget.ButtonTextColor{Idle: color.Black}),
		widget.ButtonOpts.TextPadding(&widget.Insets{
			Left:   2,
			Right:  2,
			Top:    1,
			Bottom: 1,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			ev := events.ButtonClickedEvent{buttonText}
			hub.Publish(ev)
		}),
	)

	buttonStackedLayout.AddChild(button)
	buttonStackedLayout.AddChild(headerLbl, headerLblOutline)

	return buttonStackedLayout
}

func LoadOutlineTextButtonSubmitBg(buttonText string, hub *tasks.EventHub, eventText string) *widget.Container {
	img := loadSubmitButtonImage()
	but := LoadOutlineTextButton(img, buttonText, hub, eventText)
	return but
}

func LoadOutlineTextColoreBg(clr color.RGBA, buttonText string, hub *tasks.EventHub, eventText string) *widget.Container {
	img := LoadColoredImageButton(clr)
	but := LoadOutlineTextButton(img, buttonText, hub, eventText)
	return but
}

func LoadOutlineTextButtonNoBg(buttonText string, hub *tasks.EventHub) *widget.Container {
	img := loadClearButtonImage()
	but := LoadOutlineTextButtonNotDisableOnPress(img, buttonText, hub, 35)
	return but
}

func LoadTextButtonNoBg(buttonText string, hub *tasks.EventHub, eventText string) *widget.Button {
	img := loadClearButtonImage()
	but := LoadDefualtTextButton(img, buttonText, hub, eventText)
	return but
}

func LoadDefualtTextButton(buttonImg *widget.ButtonImage, buttonText string, hub *tasks.EventHub, eventText string) *widget.Button {
	face := registry.FontMap["RockSalt"]

	// construct a pressable button
	var button *widget.Button
	txt := buttonText

	button = widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImg),

		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {

		}),
		//Move the text back to start on press released
		widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {

		}),
		widget.ButtonOpts.Text(txt, &face, &widget.ButtonTextColor{Idle: color.Black}),
		widget.ButtonOpts.TextPadding(&widget.Insets{
			Left:   18,
			Right:  18,
			Top:    1,
			Bottom: 1,
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if eventText == "" {
				println("button event generated for", buttonText)
				ev := events.ButtonClickedEvent{
					ButtonText: buttonText,
				}
				hub.Publish(ev)
			} else {
				println("button event generated for", eventText)
				ev := events.ButtonClickedEvent{
					ButtonText: eventText,
				}
				hub.Publish(ev)
			}

		}))

	return button
}

func LoadOutlineTextButton(buttonImg *widget.ButtonImage, buttonText string, hub *tasks.EventHub, eventText string) *widget.Container {

	face := registry.FontMap["RockSalt_16"]

	baseInsets := &widget.Insets{
		Left:   18,
		Right:  18,
		Top:    6,
		Bottom: 6,
	}

	offSetInset := &widget.Insets{
		Left:   20,
		Right:  16,
		Top:    8,
		Bottom: 4,
	}

	dropInset := &widget.Insets{
		Left:   20,
		Right:  16,
		Top:    8,
		Bottom: 4,
	}

	dropInsetOffset := &widget.Insets{
		Left:   22,
		Right:  14,
		Top:    10,
		Bottom: 2,
	}

	buttonStackedLayout := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout()),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	// construct a pressable button
	var button *widget.Button
	txt := buttonText

	text1 := widget.NewText(
		widget.TextOpts.Text(txt, &face, colornames.White),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
		widget.TextOpts.Padding(baseInsets))

	text2 := widget.NewText(
		widget.TextOpts.Text(txt,
			&face, colornames.Black),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
		widget.TextOpts.Padding(dropInset),
	)

	button = widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImg),

		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			button.Text().SetPadding(&widget.Insets{Top: 2, Left: 2})

			widget.ButtonOpts.TextPadding(offSetInset)

			text1.SetPadding(offSetInset)
			text2.SetPadding(dropInsetOffset)
		}),
		//Move the text back to start on press released
		widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {
			button.Text().SetPadding(&widget.Insets{Top: 0, Left: 0})
			widget.ButtonOpts.TextPadding(baseInsets)
			text1.SetPadding(baseInsets)
			text2.SetPadding(dropInset)
		}),
		widget.ButtonOpts.Text(txt, &face, &widget.ButtonTextColor{Idle: color.Black}),
		widget.ButtonOpts.TextPadding(baseInsets),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {

			if eventText == "" {
				println("button event generated for", buttonText)
				ev := events.ButtonClickedEvent{
					ButtonText: buttonText,
				}
				hub.Publish(ev)
			} else {
				println("button event generated for", eventText)
				ev := events.ButtonClickedEvent{
					ButtonText: eventText,
				}
				hub.Publish(ev)
			}
		}),
	)

	buttonStackedLayout.AddChild(button)
	buttonStackedLayout.AddChild(text2, text1)

	return buttonStackedLayout
}

func LoadSubmitButton(buttonText string, hub *tasks.EventHub, eventText string) *widget.Button {
	//load a generic button labeled with button text string that will send a button clicked event to event hub
	buttonImage := loadSubmitButtonImage()

	face := registry.FontMap["RockSalt_16"]
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
		widget.ButtonOpts.Text(buttonText, &face, &widget.ButtonTextColor{
			Idle:    color.NRGBA{0, 0, 0, 255},
			Hover:   color.NRGBA{255, 255, 100, 255},
			Pressed: color.NRGBA{0, 0, 100, 255},
		}),
		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(&widget.Insets{
			Left:   10,
			Right:  10,
			Top:    10,
			Bottom: 10,
		}),
		//Move the text down and right on press
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if eventText == "" {
				println("button event generated for", buttonText)
				ev := events.ButtonClickedEvent{
					ButtonText: buttonText,
				}
				hub.Publish(ev)
			} else {
				println("button event generated for", eventText)
				ev := events.ButtonClickedEvent{
					ButtonText: eventText,
				}
				hub.Publish(ev)
			}
		}),

		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {

		}),

		//Move the text back to start on press released
		widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) {
		}),

		// add a handler that reacts to exiting the button with the cursor
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) {
			//ResetVls the Text inset if the cursor is no longer over the button

		}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(113, 27)),
	)
	return button
}

func LoadSubmitButtonAltEvent(buttonText string, hub *tasks.EventHub, fontSize float64, eventText string) *widget.Button {
	//load a generic button labeled with button text string that will send a button clicked event to event hub
	buttonImage := loadSubmitButtonImage()

	face, err := util.LoadFont(fontSize, "nk57")
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
		widget.ButtonOpts.Text(buttonText, &face, &widget.ButtonTextColor{
			Idle:    color.NRGBA{0, 0, 100, 0xff},
			Hover:   color.NRGBA{0, 0, 100, 255},
			Pressed: color.NRGBA{0, 0, 100, 255},
		}),
		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(&widget.Insets{
			Left:   30,
			Right:  30,
			Top:    3,
			Bottom: 3,
		}),
		//Move the text down and right on press
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			println("button event generated for", eventText)
			ev := events.ButtonClickedEvent{
				ButtonText: eventText,
			}
			hub.Publish(ev)
		}),

		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			widget.ButtonOpts.TextPadding(&widget.Insets{
				Left: 35,
				Top:  5,
			})
		}),

		//Move the text back to start on press released
		widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {
		}),
		widget.ButtonOpts.CursorMovedHandler(func(args *widget.ButtonHoverEventArgs) {
		}),

		// add a handler that reacts to exiting the button with the cursor
		widget.ButtonOpts.CursorExitedHandler(func(args *widget.ButtonHoverEventArgs) {
			//ResetVls the Text inset if the cursor is no longer over the button

		}),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(113, 27)),
	)
	return button
}

func LoadBackButton(hub *tasks.EventHub) *widget.Button {
	buttonImage, err := loadBackButtonImage()
	if err != nil {
		log.Fatal(err)
	}

	var button = &widget.Button{}

	button = widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(40, 68)),
		// specify the button's text, the font face, and the color
		//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		//Move the text down and right on press
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			button.GetWidget().CustomData = true
			button.KeepPressedOnExit = true
		}),
		//Move the text back to start on press released
		widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {

			button.GetWidget().CustomData = false
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if button.GetWidget().Disabled == false {
				ev := events.ButtonClickedEvent{
					ButtonText: "back",
				}
				hub.Publish(ev)
			}

		}),

		// add a handler that reacts to entering the button with the cursor
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
			//If we moved the Text because we clicked on this button previously, move the text down and right
			if button.GetWidget().Disabled == false {
				ev := events.ButtonEvent{
					ButtonText: "backButton",
					EType:      "cursor entered",
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

			ev := events.ButtonEvent{
				ButtonText: "back",
				EType:      "cursor exited",
			}
			hub.Publish(ev)
		}),
	)
	return button
}

func LoadStackedButtonWithText(StackedButton *widget.Container, Description string, hub *tasks.EventHub, eventName string) *widget.Container {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Position: widget.RowLayoutPositionCenter}),
		),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Spacing(20),
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			)))

	face, err := util.LoadFont(12, "nk57")
	if err != nil {
		log.Fatal(err)
	}

	txtContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Position: widget.RowLayoutPositionCenter}),
		),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Spacing(20),
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			)))

	txtImg := LoadBackgroundImageForTextInput(StoreMenu)

	textarea := widget.NewTextArea(
		widget.TextAreaOpts.ScrollContainerImage(txtImg),
		widget.TextAreaOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{Position: widget.RowLayoutPositionCenter}),
				//Set the layout data for the textarea
				//including a max height to ensure the scroll bar is visible
				//Set the minimum size for the widget
				widget.WidgetOpts.MinSize(250, 90),
			),
		),

		widget.TextAreaOpts.ControlWidgetSpacing(2),
		widget.TextAreaOpts.ProcessBBCode(true),
		widget.TextAreaOpts.FontColor(color.Black),
		widget.TextAreaOpts.FontFace(&face),

		//Tell the TextArea to show the vertical scrollbar
		//Set padding between edge of the widget and where the text is drawn
		widget.TextAreaOpts.TextPadding(
			widget.Insets{Top: 10, Right: 10, Left: 10, Bottom: 10}),
		//This sets the background images for the scroll container

		//This sets the images to use for the sliders

	)

	txtContainer.AddChild(textarea)
	buyButton := LoadSubmitButton(eventName, hub, "")
	txtContainer.AddChild(buyButton)

	rootContainer.AddChild(StackedButton)
	rootContainer.AddChild(txtContainer)
	AppendTextArea(Description, textarea, 35)

	return rootContainer
}

func LoadSpriteSelectButton(buttonText string, hub *tasks.EventHub, fontSize float64) *widget.Button {
	//load a generic button labeled with button text string that will send a button clicked event to event hub

	buttonImage, err := loadSpriteSelectButtonImage(buttonText)
	if err != nil {
		log.Fatal(err)
	}

	face, err := util.LoadFont(fontSize, "nk57")
	if err != nil {
		log.Fatal(err)
	}

	var button = &widget.Button{}

	button = widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(200, 200)),
		// specify the button's text, the font face, and the color
		//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
		widget.ButtonOpts.Text(buttonText, &face, &widget.ButtonTextColor{
			Idle:     color.NRGBA{0, 0, 0, 0xff},
			Hover:    color.NRGBA{255, 255, 0, 255},
			Pressed:  color.NRGBA{255, 255, 0, 255},
			Disabled: color.NRGBA{255, 255, 0, 255},
		}),
		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(&widget.Insets{
			Left:   10,
			Right:  10,
			Top:    100,
			Bottom: 10,
		}),
		//Move the text down and right on press
		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {
			ev := events.ButtonClickedEvent{
				ButtonText: buttonText,
			}

			hub.Publish(ev)
		}),
		//Move the text back to start on press released
		widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {

			button.GetWidget().CustomData = false
			button.GetWidget().Disabled = false
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {

		}),

		// add a handler that reacts to entering the button with the cursor
		widget.ButtonOpts.CursorEnteredHandler(func(args *widget.ButtonHoverEventArgs) {
			//If we moved the Text because we clicked on this button previously, move the text down and right
			if button.GetWidget().Disabled == false {
				ev := events.ButtonEvent{
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

			ev := events.ButtonEvent{
				buttonText,
				"cursor exited",
			}
			hub.Publish(ev)
		}),
	)
	return button
}

func LoadColoredImageButton(clr color.RGBA) *widget.ButtonImage {
	nineSlice := eimage.NewNineSliceColor(clr)
	btnimg := &widget.ButtonImage{
		Idle:     nineSlice,
		Hover:    nineSlice,
		Pressed:  nineSlice,
		Disabled: nineSlice,
	}

	return btnimg
}
