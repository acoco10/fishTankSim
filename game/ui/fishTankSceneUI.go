package ui

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"image"
	"image/color"
	"log"
	"strconv"
	"strings"
)

type ButtonType uint8

const (
	SubmitButton ButtonType = iota
	SpriteSelectButton
)

type WindowType uint8

const (
	None WindowType = iota
	PHGuess
	ChooseProp
	GoToBed
	Door
	MagazineState
)

type MainMenuData struct {
	removeFunc             []widget.RemoveWindowFunc
	nextDay                bool
	currentWindowCloseable bool
	day                    int
	dayType                string
	eventHub               *tasks.EventHub
	windowOpen             WindowType
	textInputReference     *widget.TextInput
}

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

	fishStats, err := NewTextBlock(eHub, StatsMenu)

	if err != nil {
		return nil, nil, err
	}

	fishStats.text.GetWidget().Visibility = widget.Visibility_Hide

	//notePad, err := NewTextBlock(eHub, NotePad)
	//if err != nil {

	//notePad.text.SetText("To Do:")

	buttonContainer.AddChild(fishStats)
	rootContainer.AddChild(buttonContainer)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	mag, err := LoadMagazineUiMenu(eHub, gameWidth, gameHeight)
	if err != nil {
		return nil, nil, err
	}
	mainDat := &MainMenuData{eventHub: eHub}
	MainMenuSubs(mainDat, mag, &ui, eHub)

	return &ui, fishStats, nil
}

func MakePHmenu(hub *tasks.EventHub) (*widget.Container, *widget.TextInput) {

	face, err := util.LoadFont(32, "reglisseOutline") //white center text

	face2, err := util.LoadFont(32, "reglisseOutlined") //black outline

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
		//widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(200)))))
		)))

	bImage, err := loadOptionsMenuInputImage()
	if err != nil {
		log.Fatal(err)
	}

	childContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(bImage),
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

	buttonRow := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			//Define number of columns in the grid
			widget.GridLayoutOpts.Columns(3),
			//onlt one row so row spacing second input doesnt really matter
			//widget.GridLayoutOpts.Spacing(),
			// DefaultStretch values will be used when extra columns/rows are used
			// out of the ones defined on the normal Stretch
			widget.GridLayoutOpts.DefaultStretch(false, true),
			//Define how to stretch the rows and columns.
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false}),
		),
		))

	buttonRow2 := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	textContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	headContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	headerLbl := widget.NewText(
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			})),

		widget.TextOpts.Text("Guess the PH!", face, color.RGBA{R: 255, G: 255, B: 255, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.Insets(widget.Insets{}),
	)

	headerLblOutline := widget.NewText(
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			})),

		widget.TextOpts.Text("Guess the PH!", face2, color.RGBA{R: 0, G: 0, B: 0, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.Insets(widget.Insets{}),
	)

	inputFace := registry.FontMap["nk57_24"]

	img, _ := loadTextInputImage()

	textInput := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			//Set the layout information to center the textbox in the parent
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),

			widget.WidgetOpts.MinSize(50, 10),
		),
		//Set the Idle and Disabled background image for the text input
		//If the NineSlice image has a minimum size, the widget will use that or
		// widget.WidgetOpts.MinSize; whichever is greater
		widget.TextInputOpts.Image(img),

		//Set the font face and size for the widget
		widget.TextInputOpts.Face(inputFace),

		//Set the colors for the text and caret
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          color.NRGBA{0, 0, 50, 255},
			Disabled:      color.NRGBA{0, 0, 20, 100},
			Caret:         color.NRGBA{0, 0, 50, 255},
			DisabledCaret: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		}),

		//Set how much padding there is between the edge of the input and the text
		widget.TextInputOpts.Padding(widget.NewInsetsSimple(5)),

		//Set the font and width of the caret
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(inputFace, 0),
		),

		//This text is displayed if the input is empty
		widget.TextInputOpts.Placeholder("7.0"),

		//This is called when the user hits the "Enter" key.
		//There are other options that can configure this behavior
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			ev := entities.SendData{
				DataFor: "ph guess",
				Data:    args.InputText,
			}

			hub.Publish(ev)
		}),

		//This is called whenever there is a change to the text
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
		}),
	)

	b1 := LoadOutlineTextButtonNoBg("<", hub)

	b2 := LoadOutlineTextButtonNoBg(">", hub)

	b3 := LoadOutlineTextButtonSubmitBg("Guess", hub, "")

	textInput.SetText("7.0")
	textContainer.AddChild(textInput)

	buttonRow.AddChild(b1, textContainer, b2)
	buttonRow2.AddChild(b3)

	headContainer.AddChild(headerLbl, headerLblOutline)

	childContainer.AddChild(headContainer, buttonRow, buttonRow2)
	rootContainer.AddChild(childContainer)

	return rootContainer, textInput
}

func makeWindow(rtContainer *widget.Container, y int, width int, height int) *widget.Window {

	windowWidth := 400
	windowHeight := 400
	if width != 0 {
		windowWidth = width
	}
	if height != 0 {
		windowHeight = height
	}

	var y0 int

	if y == 0 {
		y0 = registry.Config.ResolutionHeight/3 - windowHeight/2
	} else {
		y0 = y
	}
	x0 := registry.Config.ResolutionWidth/2 - windowWidth/2

	window := widget.NewWindow(
		//Set the main contents of the window
		widget.WindowOpts.Contents(rtContainer),
		//Set the titlebar for the window (Optional)
		//Set the window above everything else and block input elsewhere
		widget.WindowOpts.Modal(),
		//Set how to close the window. CLICK_OUT will close the window when clicking anywhere
		//that is not a part of the window object
		//widget.WindowOpts.CloseMode(widget.CLICK_OUT),
		//Indicates that the window is draggable. It must have a TitleBar for this to work
		//Set the window resizeable
		//Set the minimum size the window can be
		widget.WindowOpts.MinSize(windowWidth, windowHeight),
		//Set the maximum size a window can be
		widget.WindowOpts.MaxSize(windowWidth, windowHeight),
		//Set the callback that triggers when a move is complete
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Moved")
		}),
		//Set the callback that triggers when a resize is complete
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Resized")
		}),

		widget.WindowOpts.Location(image.Rect(x0, y0, x0+windowWidth, y0+windowHeight)),
	)
	return window
}

func triggerPHwindow(ui *ebitenui.UI, hub *tasks.EventHub) (widget.RemoveWindowFunc, *widget.TextInput) {
	cont, txtInput := MakePHmenu(hub)
	window := makeWindow(cont, registry.Config.ResolutionHeight/2+100, 400, 130)
	rfunc := ui.AddWindow(window)
	return rfunc, txtInput
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

func TriggerTextWindow(hub *tasks.EventHub, ui *ebitenui.UI, header string, inputText []string) widget.RemoveWindowFunc {

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
		//widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(200)))))
		)))

	bimg, err := loadOptionsMenuInputImage()
	if err != nil {
		log.Fatal(err)
	}

	childContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(bimg),
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

	btnContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	headerLbl := widget.NewText(
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			})),

		widget.TextOpts.Text(header, registry.FontMap["nk57_24"], color.RGBA{R: 250, G: 160, B: 0, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionStart),
		widget.TextOpts.Insets(widget.Insets{}),
	)

	var lines []*widget.Text
	for _, line := range inputText {
		textCon := widget.NewText(
			widget.TextOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				})),

			widget.TextOpts.Text(line, registry.FontMap["nk57_24"], color.RGBA{R: 250, G: 160, B: 0, A: 255}),
			widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionStart),
			widget.TextOpts.Insets(widget.Insets{}),
		)
		lines = append(lines, textCon)
	}

	childContainer.AddChild(headerLbl)

	for _, line := range lines {
		childContainer.AddChild(line)
	}

	b1 := LoadOutlineTextButtonSubmitBg("Lets Mow!", hub, "")
	btnContainer.AddChild(b1)
	childContainer.AddChild(btnContainer)

	rootContainer.AddChild(childContainer)
	wind := makeWindow(rootContainer, 100, 800, 400)

	rfunc := ui.AddWindow(wind)
	return rfunc

}

func TriggerImageButtonSelectWindow(ui *ebitenui.UI, hub *tasks.EventHub) widget.RemoveWindowFunc {

	castleImg, err := util.LoadImageAssetAsEbitenImage("menuAssets/castleSelectButton")
	logImg, err := util.LoadImageAssetAsEbitenImage("menuAssets/logSelectButton")

	propButton, err := LoadStackSpriteSelectButtonWithToolTip("Castle", castleImg, 16, hub, []string{"+Hiding Spot", "-PH"})

	if err != nil {
		log.Fatal(err)
	}
	propButton2, err := LoadStackSpriteSelectButtonWithToolTip("Log", logImg, 16, hub, []string{"+PH", "+Eco"})

	if err != nil {
		log.Fatal(err)
	}

	image, err := loadOptionsMenuInputImage()

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
		//widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(200)))))
		)))

	childContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image),
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

	header := MakeOutlineText("Pick your first Tank Decoration!")

	widge := addPickSpriteContainer(2)

	widge.AddChild(propButton, propButton2)

	y0 := registry.Config.ResolutionHeight/2 - int(50*registry.Config.ResolutionScalingF) - 100

	window := makeWindow(rootContainer, y0, 600, 0)

	bContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	b1 := LoadOutlineTextButtonSubmitBg("Confirm", hub, "Confirm for prop select")

	childContainer.AddChild(header)
	childContainer.AddChild(widge)
	bContainer.AddChild(b1)
	childContainer.AddChild(bContainer)
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

func MainMenuAddWindow(hub *tasks.EventHub) {
	ev := events.AddWindow{}
	hub.Publish(ev)
}

func SendWindowOpened(hub *tasks.EventHub) {
	ev := events.WindowOpened{}
	hub.Publish(ev)
}

func (m *MainMenuData) addRemoveFunc(windowFunc widget.RemoveWindowFunc) {
	m.removeFunc = append(m.removeFunc, windowFunc)
}
func (m *MainMenuData) RemoveFunc(args []string) {
	for _, f := range m.removeFunc {
		f()
	}
	ev2 := events.WindowClosed{}
	m.eventHub.Publish(ev2)
	for _, arg := range args {
		if arg == "maintainState" {
			return
		}
	}
	m.windowOpen = None
}

func closeWindow(hub *tasks.EventHub, data *MainMenuData) {
	data.RemoveFunc([]string{})
	ev2 := events.WindowClosed{}
	hub.Publish(ev2)
}

func MainMenuSubs(mainDat *MainMenuData, magazine *Magazine, ui *ebitenui.UI, hub *tasks.EventHub) {

	buttonTextSubs(mainDat, magazine, ui, hub)
	uiSpriteActionSubs(mainDat, magazine, ui, hub)

	hub.Subscribe(tasks.TaskCompleted{}, func(e tasks.Event) {
		ev := e.(tasks.TaskCompleted)
		if ev.Task.Text == "1. Take a ph reading of your tank" && mainDat.day == 1 {
			{
				SendWindowOpened(hub)
				mainDat.addRemoveFunc(TriggerImageButtonSelectWindow(ui, hub))
				mainDat.windowOpen = 0
			}
		}
	})

	hub.Subscribe(events.CloseWindow{}, func(e tasks.Event) {
		switch mainDat.windowOpen {
		case PHGuess:
			return
		case ChooseProp:
			return
		case MagazineState:
			mainDat.RemoveFunc([]string{})
			magazine.activeIndex = 0
		case GoToBed:
			mainDat.RemoveFunc([]string{})
		case Door:
			mainDat.RemoveFunc([]string{})
		default:
			mainDat.RemoveFunc([]string{})
		}
	})

	hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		ev := e.(events.NewDay)
		mainDat.day = ev.Day
		mainDat.dayType = ev.DayType
	})

	hub.Subscribe(events.DayOver{}, func(e tasks.Event) {
		if mainDat.removeFunc != nil {
			mainDat.RemoveFunc([]string{})
		}
	})

}

func buttonTextSubs(mainDat *MainMenuData, magazine *Magazine, ui *ebitenui.UI, hub *tasks.EventHub) {
	hub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		switch mainDat.windowOpen {
		case MagazineState:
			MagazineWindowButtonEvent(mainDat, magazine, ui, ev, hub)
		case ChooseProp:
			if ev.ButtonText == "Confirm for prop select" {
				closeWindow(hub, mainDat)
			}
		case Door:
			DoorWindowButtonEvent(mainDat, magazine, ui, ev, hub)
		case PHGuess:
			PhGuessWindowButtonEvent(mainDat, ev, hub)
		case GoToBed:
			if ev.ButtonText == "Vibe awhile Longer" {
				mainDat.RemoveFunc([]string{})
			}
		}
	})
}

func uiSpriteActionSubs(mainDat *MainMenuData, magazine *Magazine, ui *ebitenui.UI, hub *tasks.EventHub) {
	hub.Subscribe(events.UISpriteAction{}, func(e tasks.Event) {
		ev := e.(events.UISpriteAction)
		if ev.UiSprite == "pillow" && ev.UiSpriteAction == "clicked" {
			SendWindowOpened(hub)
			mainDat.windowOpen = GoToBed
			mainDat.addRemoveFunc(TriggerNextDayWindow(ui, hub))
		}

		if ev.UiSprite == "magazine" {
			mainDat.windowOpen = MagazineState
			SendWindowOpened(hub)
			mainDat.addRemoveFunc(TriggerMagWindow(magazine, ui, hub))
		}

		if ev.UiSprite == "door" {
			mainDat.windowOpen = Door
			switch mainDat.dayType {
			case "Camp":
				SendWindowOpened(hub)
				mainDat.addRemoveFunc(TriggerOptionWindow("Go to Camp?", ui, hub))
			case "Chores":
				SendWindowOpened(hub)
				mainDat.addRemoveFunc(TriggerOptionWindow("Go do your Chores?", ui, hub))
			}
		}

		if ev.UiSprite == "phreader" && ev.UiSpriteAction == "ph reading" {
			SendWindowOpened(hub)
			mainDat.windowOpen = PHGuess
			rfunc, tarea := triggerPHwindow(ui, hub)
			mainDat.textInputReference = tarea
			mainDat.addRemoveFunc(rfunc)
		}
	})
}

func PhGuessWindowButtonEvent(mainDat *MainMenuData, ev events.ButtonClickedEvent, hub *tasks.EventHub) {
	switch ev.ButtonText {
	case ">":
		curr := mainDat.textInputReference.GetText()
		phVal, err := strconv.ParseFloat(curr, 64)
		if err != nil {
			log.Printf("Error parsing '%s': %v", curr, err) // This will show you exactly what's in the string
		}
		phVal += 0.1
		phValStr := fmt.Sprintf("%.1f", phVal)
		mainDat.textInputReference.SetText(phValStr)
	case "<":
		curr := mainDat.textInputReference.GetText()
		phVal, err := strconv.ParseFloat(curr, 64)
		if err != nil {
			log.Printf("Error parsing '%s': %v", curr, err) // This will show you exactly what's in the string
		}
		phVal -= 0.1
		phValStr := fmt.Sprintf("%.1f", phVal)
		mainDat.textInputReference.SetText(phValStr)
	case "Guess":
		curr := mainDat.textInputReference.GetText()
		phVal, err := strconv.ParseFloat(curr, 64)
		if err != nil {
			log.Printf("Error parsing '%s': %v", curr, err) // This will show you exactly what's in the string
		}
		ev2 := events.PHGuess{
			Guess: phVal,
		}
		hub.Publish(ev2)
		mainDat.RemoveFunc([]string{})

	}
}

func MagazineWindowButtonEvent(mainDat *MainMenuData, magazine *Magazine, ui *ebitenui.UI, ev events.ButtonClickedEvent, hub *tasks.EventHub) {
	switch ev.ButtonText {
	case "Fish":
		mainDat.RemoveFunc([]string{"maintainState"})
		mainDat.addRemoveFunc(TriggerMagWindow(magazine, ui, hub))
	case "Info":
		mainDat.RemoveFunc([]string{"maintainState"})
		mainDat.addRemoveFunc(TriggerMagWindow(magazine, ui, hub))
	case "Helpers":
		mainDat.RemoveFunc([]string{"maintainState"})
		mainDat.addRemoveFunc(TriggerMagWindow(magazine, ui, hub))
	case "ph+", "ph-":
		closeWindow(hub, mainDat)
		magazine.activeIndex = 0
	}

	if strings.HasPrefix(ev.ButtonText, "Buy:") {
		closeWindow(hub, mainDat)
		magazine.activeIndex = 0
	}

}

func DoorWindowButtonEvent(mainDat *MainMenuData, magazine *Magazine, ui *ebitenui.UI, ev events.ButtonClickedEvent, hub *tasks.EventHub) {
	switch ev.ButtonText {
	case "Go to Bed?: Yes":
		closeWindow(hub, mainDat)
	case "Go do your Chores?: Yes", "Go do your Chores?: No":
		mainDat.RemoveFunc([]string{})
	case "Go to Camp?: Yes", "Go to Camp?: No":
		mainDat.RemoveFunc([]string{})
	}
}
