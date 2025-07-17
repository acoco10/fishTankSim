package ui

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"log"
)

type ButtonType uint8

const (
	SubmitButton ButtonType = iota
	SpriteSelectButton
)

func LoadMainFishMenu(gameWidth, gameHeight int, eHub *tasks.EventHub) (*ebitenui.UI, *TextBoxUi, error) {

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
					VerticalPosition:   widget.AnchorLayoutPositionStart,
				}),
		),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Spacing(20),
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(widget.Insets{Right: 0, Left: 0, Top: 100, Bottom: 0}),
			)),
	)

	button := LoadSubmitButton("Save", eHub, 16, "")

	fishStats, err := NewTextBlock(eHub, StatsMenu)

	if err != nil {
		return nil, nil, err
	}

	fishStats.text.GetWidget().Visibility = widget.Visibility_Hide

	//notePad, err := NewTextBlock(eHub, NotePad)
	//if err != nil {

	//notePad.text.SetText("To Do:")

	buttonContainer.AddChild(fishStats)
	buttonContainer.AddChild(button)
	rootContainer.AddChild(buttonContainer)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	mag, err := LoadMagazineUiMenu(eHub, gameWidth, gameHeight)
	if err != nil {
		return nil, nil, err
	}

	MainMenuSubs(mag, &ui, eHub)

	return &ui, fishStats, nil
}

func TriggerNextDayWindow(ui *ebitenui.UI, hub *tasks.EventHub) widget.RemoveWindowFunc {

	buttonText := []string{"Yes", "vibe awhile longer"}

	ndUI, err := LoadNextOptionsMenuUI("Go to Bed?", buttonText, hub)
	if err != nil {
		log.Fatal(err, "error loading next day UI")
	}

	removeFunc := ui.AddWindow(ndUI)
	return removeFunc
}

func TriggerTextWindow(hub *tasks.EventHub, ui *ebitenui.UI) widget.RemoveWindowFunc {
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
		//widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(200)))))
		)))

	childContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(colornames.Darkmagenta)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  true,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(20),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(50)),
		)))

	headerLbl := widget.NewText(
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			})),

		widget.TextOpts.Text("Pick your first Decoration", registry.FontMap["nk57_12"], color.RGBA{R: 250, G: 160, B: 0, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionStart),
		widget.TextOpts.Insets(widget.Insets{}),
	)

	textBox, err := NewTextBlock(hub, StoreMenu)
	if err != nil {
		log.Fatal(err)
	}

	childContainer.AddChild(headerLbl)
	childContainer.AddChild(textBox)
	rootContainer.AddChild(childContainer)

	window := widget.NewWindow(
		//Set the main contents of the window
		widget.WindowOpts.Contents(rootContainer),
		//Set the titlebar for the window (Optional)
		//Set the window above everything else and block input elsewhere
		widget.WindowOpts.Modal(),
		//Set how to close the window. CLICK_OUT will close the window when clicking anywhere
		//that is not a part of the window object
		//widget.WindowOpts.CloseMode(widget.CLICK_OUT),
		//Indicates that the window is draggable. It must have a TitleBar for this to work
		//Set the window resizeable
		//Set the minimum size the window can be
		widget.WindowOpts.MinSize(200, 100),
		//Set the maximum size a window can be
		widget.WindowOpts.MaxSize(300, 300),
		//Set the callback that triggers when a move is complete
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Moved")
		}),
		//Set the callback that triggers when a resize is complete
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Resized")
		}),

		widget.WindowOpts.Location(image.Rect(500, 100, 700, 400)),
	)

	rfunc := ui.AddWindow(window)
	return rfunc

}

func TriggerImageButtonSelectWindow(screenWidth int, screenHeight int, ui *ebitenui.UI, hub *tasks.EventHub) widget.RemoveWindowFunc {

	windowWidth := 400
	windowHeight := 400

	x0 := screenWidth/2 - windowWidth/2
	y0 := screenHeight/2 - windowHeight/2
	x1 := screenWidth/2 + windowWidth/2
	y1 := screenWidth/2 + windowHeight/2

	castleImg, err := loader.LoadImageAssetAsEbitenImage("menuAssets/castleSelectButton")
	logImg, err := loader.LoadImageAssetAsEbitenImage("menuAssets/logSelectButton")

	propButton, err := LoadStackSpriteSelectButton("Castle", castleImg, 16, hub, []string{"+Hiding Spot", "-PH"})

	if err != nil {
		log.Fatal(err)
	}
	propButton2, err := LoadStackSpriteSelectButton("Log", logImg, 16, hub, []string{"+PH", "+Eco"})

	if err != nil {
		log.Fatal(err)
	}

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
		//widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(200)))))
		)))

	childContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(colornames.Darkmagenta)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  true,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(20),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(50)),
		)))

	headerLbl := widget.NewText(
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			})),

		widget.TextOpts.Text("Pick your first Decoration", registry.FontMap["nk57_12"], color.RGBA{R: 250, G: 160, B: 0, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionStart),
		widget.TextOpts.Insets(widget.Insets{}),
	)

	widge := addPickSpriteContainer(2)

	widge.AddChild(propButton, propButton2)
	window := widget.NewWindow(
		//Set the main contents of the window
		widget.WindowOpts.Contents(rootContainer),
		//Set the titlebar for the window (Optional)
		//Set the window above everything else and block input elsewhere
		widget.WindowOpts.Modal(),
		//Set how to close the window. CLICK_OUT will close the window when clicking anywhere
		//that is not a part of the window object
		//widget.WindowOpts.CloseMode(widget.CLICK_OUT),
		//Indicates that the window is draggable. It must have a TitleBar for this to work
		//Set the window resizeable
		//Set the minimum size the window can be
		widget.WindowOpts.MinSize(200, 100),
		//Set the maximum size a window can be
		widget.WindowOpts.MaxSize(300, 300),
		//Set the callback that triggers when a move is complete
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Moved")
		}),
		//Set the callback that triggers when a resize is complete
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Resized")
		}),

		widget.WindowOpts.Location(image.Rect(x0, y0, x1, y1)),
	)

	b1 := LoadSubmitButton("Confirm", hub, 16, "Confirm for prop select")

	childContainer.AddChild(headerLbl)
	childContainer.AddChild(widge)
	childContainer.AddChild(b1)
	rootContainer.AddChild(childContainer)

	removeFunc := ui.AddWindow(window)
	return removeFunc
}

func TriggerOptionWindow(headerText string, ui *ebitenui.UI, hub *tasks.EventHub) widget.RemoveWindowFunc {

	buttonText := []string{"Yes", "No"}

	ndUI, err := LoadNextOptionsMenuUI(headerText, buttonText, hub)
	if err != nil {
		log.Fatal(err, "error loading next day UI")
	}

	removeFunc := ui.AddWindow(ndUI)
	return removeFunc
}

func TriggerMagWindow(mag *Magazine, ui *ebitenui.UI, hub *tasks.EventHub) widget.RemoveWindowFunc {

	activeContainer := mag.ActiveWindow()

	window := widget.NewWindow(

		widget.WindowOpts.Contents(activeContainer),

		widget.WindowOpts.Modal(),

		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Moved")
		}),
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Resized")
		}),
		widget.WindowOpts.Location(image.Rect(10, 10, 500, 500)),
	)

	removeFunc := ui.AddWindow(window)

	return removeFunc
}

func MainMenuSubs(magazine *Magazine, ui *ebitenui.UI, hub *tasks.EventHub) {
	var nextDay bool
	var removeFunc widget.RemoveWindowFunc
	var Closeable = true

	hub.Subscribe(events.UISpriteAction{}, func(e tasks.Event) {
		ev := e.(events.UISpriteAction)
		if ev.UiSprite == "pillow" && ev.UiSpriteAction == "clicked" && !nextDay {
			nextDay = true
			removeFunc = TriggerNextDayWindow(ui, hub)
		}
		if ev.UiSprite == "magazine" {
			removeFunc = TriggerMagWindow(magazine, ui, hub)
			Closeable = true
		}
		if ev.UiSprite == "door" {
			switch ev.UiSpriteAction {
			case "clickedCamp":
				removeFunc = TriggerOptionWindow("Go to Camp?", ui, hub)
			case "clickedChores":
				removeFunc = TriggerOptionWindow("Go do your Chores?", ui, hub)
			}
			Closeable = true
		}
	})

	hub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		//subs for updating mag window
		ev := e.(events.ButtonClickedEvent)
		println(ev.ButtonText, "button event received")
		switch ev.ButtonText {
		case "Fish":
			removeFunc()
			removeFunc = TriggerMagWindow(magazine, ui, hub)
		case "Info":
			removeFunc()
			removeFunc = TriggerMagWindow(magazine, ui, hub)
		}

	})

	hub.Subscribe(tasks.TaskCompleted{}, func(e tasks.Event) {
		ev := e.(tasks.TaskCompleted)
		if ev.Task.Text == "1. Take a ph reading of your tank" {
			width, height := ebiten.WindowSize()
			removeFunc = TriggerImageButtonSelectWindow(width, height, ui, hub)
			Closeable = false
		}
	})

	hub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		if ev.ButtonText == "Vibe awhile Longer" {
			removeFunc()
			nextDay = false
			Closeable = true
		}
	})

	hub.Subscribe(events.DayOver{}, func(e tasks.Event) {
		if removeFunc != nil {
			removeFunc()
			Closeable = true
		}
		nextDay = false
	})

	hub.Subscribe(events.CloseWindow{}, func(e tasks.Event) {
		ev := e.(events.CloseWindow)
		if removeFunc != nil && Closeable || ev.OverRide {
			removeFunc()
			Closeable = true
		}
	})

}

func loadSubmitButtonImage() *widget.ButtonImage {

	img, _ := loader.LoadImageAssetAsEbitenImage("menuAssets/submitButton3")

	imgClicked, _ := loader.LoadImageAssetAsEbitenImage("menuAssets/submitButtonAlt")

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
	}
}

func loadSpriteSelectButtonImage(t string) (*widget.ButtonImage, error) {

	img, err := loader.LoadImageAssetAsEbitenImage("menuAssets/spriteOutlineButton")

	if err != nil {
		return nil, err
	}

	imgClicked, err := loader.LoadImageAssetAsEbitenImage("menuAssets/spriteOutlineButtonAlt")
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
			Idle:     hover,
			Hover:    hover,
			Pressed:  pressed,
			Disabled: pressed,
		}, nil
	} else {
		return &widget.ButtonImage{
			Idle:     idle,
			Hover:    hover,
			Pressed:  pressed,
			Disabled: pressed,
		}, nil
	}
}

func loadBackButtonImage() (*widget.ButtonImage, error) {

	img, err := loader.LoadImageAssetAsEbitenImage("menuAssets/backButton")

	if err != nil {
		return nil, err
	}

	nineSliceImage := eimage.NewNineSlice(img, [3]int{(59 - 32) / 2, 38, 59 - 32/2}, [3]int{(74 - 66) / 2, 66, (74 - 66) / 2})

	idle := nineSliceImage

	hover := nineSliceImage

	pressed := nineSliceImage

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}

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
			button.Text().Inset.Top = 0
			button.Text().Inset.Left = 0
			ev := events.ButtonEvent{
				buttonText,
				"cursor exited",
			}
			hub.Publish(ev)
		}),
	)
	return button
}

func LoadStackedImageButton(buttonText string, img *ebiten.Image, fontSize float64, hub *tasks.EventHub) {

}

func LoadStackSpriteSelectButton(buttonText string, fishImg *ebiten.Image, fontSize float64, hub *tasks.EventHub, tipText []string) (*widget.Container, error) {

	tooltipContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(widget.RowLayoutOpts.Direction(widget.DirectionVertical))),
		widget.ContainerOpts.AutoDisableChildren(),
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 230, A: 255})),
	)

	if len(tipText) > 0 {
		for _, text := range tipText {
			option1 := widget.NewText(
				widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
				widget.TextOpts.Text(fmt.Sprint(text), registry.FontMap["nk57"], colornames.Greenyellow),
				widget.TextOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 0)),
			)
			tooltipContainer.AddChild(option1)
		}
	}

	face, err := util.LoadFont(fontSize, "nk57")
	if err != nil {
		log.Fatal(err)
	}

	imgForTransform := ebiten.NewImage(fishImg.Bounds().Dx()*2, fishImg.Bounds().Dx()*2)
	dopts := &ebiten.DrawImageOptions{}
	dopts.GeoM.Scale(2, 2)

	imgForTransform.DrawImage(fishImg, dopts)

	img, err := loader.LoadImageAssetAsEbitenImage("menuAssets/spriteOutlineButton")
	if err != nil {
		return nil, err
	}

	imgClicked, err := loader.LoadImageAssetAsEbitenImage("menuAssets/spriteOutlineButtonAlt")
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
			button.Text().Inset.Top = 0
			button.Text().Inset.Left = 0
			button.GetWidget().CustomData = false
		}),
		widget.ButtonOpts.Text(buttonText, face, &widget.ButtonTextColor{
			Idle:     color.NRGBA{0, 0, 0, 0xff},
			Hover:    color.NRGBA{255, 255, 0, 255},
			Pressed:  color.NRGBA{255, 255, 0, 255},
			Disabled: color.NRGBA{255, 255, 0, 255},
		}),
		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
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
	// Put an image on top of the button, it will be centered.
	// If your image doesn't fit the button and there is no Y stretching support,
	// you may see a transparent rectangle inside the button.
	// To fix that, either use a separate button image (that can fit the image)
	// or add an appropriate stretching.
	buttonStackedLayout.AddChild(
		btnIconG,
	)

	return buttonStackedLayout, nil
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
				widget.RowLayoutOpts.Padding(widget.Insets{}),
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
				widget.RowLayoutOpts.Padding(widget.Insets{}),
			)))

	txtImg := LoadBackgroundImageForTextInput(StoreMenu)

	textarea := widget.NewTextArea(
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
		widget.TextAreaOpts.FontFace(face),

		//Tell the TextArea to show the vertical scrollbar
		//Set padding between edge of the widget and where the text is drawn
		widget.TextAreaOpts.TextPadding(
			widget.Insets{Top: 10, Right: 10, Left: 10, Bottom: 10}),
		//This sets the background images for the scroll container
		widget.TextAreaOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(txtImg),
		),
		//This sets the images to use for the sliders
		widget.TextAreaOpts.SliderOpts(
			widget.SliderOpts.Images(
				// Set the track images
				&widget.SliderTrackImage{
					Idle:  eimage.NewNineSliceColor(color.NRGBA{R: 200, G: 200, B: 200, A: 255}),
					Hover: eimage.NewNineSliceColor(color.NRGBA{R: 200, G: 200, B: 200, A: 255}),
				},
				// Set the handle images
				&widget.ButtonImage{
					Idle:    eimage.NewNineSliceColor(color.NRGBA{R: 255, G: 100, B: 100, A: 255}),
					Hover:   eimage.NewNineSliceColor(color.NRGBA{R: 255, G: 100, B: 100, A: 255}),
					Pressed: eimage.NewNineSliceColor(color.NRGBA{R: 255, G: 100, B: 100, A: 255}),
				},
			),
		),
	)

	txtContainer.AddChild(textarea)
	buyButton := LoadSubmitButton(eventName, hub, 12, "")
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
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(120, 100)),
		// specify the button's text, the font face, and the color
		//widget.ButtonOpts.Text("Hello, World!", face, &widget.ButtonTextColor{
		widget.ButtonOpts.Text(buttonText, face, &widget.ButtonTextColor{
			Idle:     color.NRGBA{0, 0, 0, 0xff},
			Hover:    color.NRGBA{255, 255, 0, 255},
			Pressed:  color.NRGBA{255, 255, 0, 255},
			Disabled: color.NRGBA{255, 255, 0, 255},
		}),
		widget.ButtonOpts.TextProcessBBCode(true),
		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
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

			button.Text().Inset.Top = 2
			button.Text().Inset.Left = 2
			hub.Publish(ev)
		}),
		//Move the text back to start on press released
		widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {
			button.Text().Inset.Top = 0
			button.Text().Inset.Left = 0
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
			button.Text().Inset.Top = 0
			button.Text().Inset.Left = 0
			ev := events.ButtonEvent{
				buttonText,
				"cursor exited",
			}
			hub.Publish(ev)
		}),
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
			button.Text().Inset.Top = 0
			button.Text().Inset.Left = 0
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
			button.Text().Inset.Top = 0
			button.Text().Inset.Left = 0
			ev := events.ButtonEvent{
				ButtonText: "back",
				EType:      "cursor exited",
			}
			hub.Publish(ev)
		}),
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

		}),

		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {

		}),

		//Move the text back to start on press released
		widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {
			if button.GetWidget().Disabled == false {
				button.GetWidget().Disabled = true
				println("button event generated for", eventText)
				ev := events.ButtonClickedEvent{
					ButtonText: eventText,
				}
				hub.Publish(ev)
			}
			button.Text().Inset.Top = 2
			button.Text().Inset.Left = -2
			button.GetWidget().Disabled = false
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

func LoadSubmitButton(buttonText string, hub *tasks.EventHub, fontSize float64, eventText string) *widget.Button {
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

		}),

		widget.ButtonOpts.PressedHandler(func(args *widget.ButtonPressedEventArgs) {

		}),

		//Move the text back to start on press released
		widget.ButtonOpts.ReleasedHandler(func(args *widget.ButtonReleasedEventArgs) {
			if button.GetWidget().Disabled == false {
				button.GetWidget().Disabled = true
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
			}
			button.Text().Inset.Top = 2
			button.Text().Inset.Left = -2
			button.GetWidget().Disabled = false
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

func LoadButton(buttonText string, hub *tasks.EventHub, fontSize float64) *widget.Button {
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
			ev := events.ButtonClickedEvent{
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
			ev := events.ButtonClickedEvent{
				ButtonText: buttonText,
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
