package ui

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"image"
	"image/color"
	"log"
	"strconv"
	"strings"
	"time"
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
	StoreMagazineState
	DebugText
	TutorialWindow
	JournalMagazineState
)

type MainMenuData struct {
	removeFunc             []widget.RemoveWindowFunc
	nextDay                bool
	currentWindowCloseable bool
	day                    int
	dayType                string
	storeMagazine          *Magazine
	journal                *Journal
	eventHub               *tasks.EventHub
	windowOpen             WindowType
	textInputReference     *widget.TextInput
	tempParam              string
}

func LoadMainFishMenu(gameWidth, gameHeight int, eHub *tasks.EventHub) (*ebitenui.UI, error) {

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
				widget.RowLayoutOpts.Padding(&widget.Insets{Right: 0, Left: 0, Top: 100, Bottom: 0}),
			)),
	)

	fishStats, err := NewTextBlock(eHub, StatsMenu)

	if err != nil {
		return nil, err
	}

	fishStats.text.GetWidget().Visibility = widget.Visibility_Hide
	buttonContainer.AddChild(fishStats)
	rootContainer.AddChild(buttonContainer)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	mag, err := CreateFishMagazine(eHub)
	if err != nil {
		return nil, fmt.Errorf("error creating mag:%s", err)
	}

	journal := CreateJournal(eHub)
	mainDat := &MainMenuData{eventHub: eHub, journal: journal}
	MainMenuSubs(mainDat, mag, &ui, eHub)

	return &ui, nil
}

func MakePHmenu(hub *tasks.EventHub) (*widget.Container, *widget.TextInput) {

	face := registry.FontMap["RockSalt"] //white center text

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

	/*	textRow := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			}),
		))*/

	buttonRow := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			//Define number of columns in the grid
			widget.GridLayoutOpts.Columns(4),
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

		widget.TextOpts.Text("Guess the PH!", &face, color.RGBA{R: 255, G: 255, B: 255, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	)

	headerLblOutline := widget.NewText(
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			})),

		widget.TextOpts.Text("Guess the PH!", &face, color.RGBA{R: 0, G: 0, B: 0, A: 150}),
		widget.TextOpts.Padding(&widget.Insets{Left: 2, Top: 2}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	)

	img, _ := loadTextInputImage()

	textInput := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			//Set the layout information to center the textbox in the parent
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),

			widget.WidgetOpts.MinSize(120, 10),
		),
		//Set the Idle and Disabled background image for the text input
		//If the NineSlice image has a minimum size, the widget will use that or
		// widget.WidgetOpts.MinSize; whichever is greater
		widget.TextInputOpts.Image(img),

		//Set the font face and size for the widget
		widget.TextInputOpts.Face(&face),

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

	b1 := LoadOutlineTextButtonSubmitBg("+0.1", hub, "")

	b2 := LoadOutlineTextButtonSubmitBg("-0.1", hub, "")

	b3 := LoadOutlineTextButtonSubmitBg("-1.0", hub, "")

	b4 := LoadOutlineTextButtonSubmitBg("+1.0", hub, "")

	b5 := LoadOutlineTextButtonSubmitBg("Guess", hub, "")

	textInput.SetText("7.0")
	textContainer.AddChild(textInput)

	buttonRow.AddChild(b2, b3, b4, b1)
	buttonRow2.AddChild(b5)

	headContainer.AddChild(headerLblOutline, headerLbl)

	childContainer.AddChild(headContainer, textInput, buttonRow, buttonRow2)
	rootContainer.AddChild(childContainer)

	return rootContainer, textInput
}

func MakeDebugTextInput(lastText string, hub *tasks.EventHub) (*widget.Container, *widget.TextInput) {

	face := registry.FontMap["RockSalt_24"]

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
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(100)),
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

		widget.TextOpts.Text("Collision Name", &face, color.RGBA{R: 255, G: 255, B: 255, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
	)

	headerLblOutline := widget.NewText(
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			})),

		widget.TextOpts.Text("Collision Name", &face, color.RGBA{R: 0, G: 0, B: 0, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.Padding(&widget.Insets{Left: 2, Top: 2}),
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

			widget.WidgetOpts.MinSize(300, 10),
		),
		//Set the Idle and Disabled background image for the text input
		//If the NineSlice image has a minimum size, the widget will use that or
		// widget.WidgetOpts.MinSize; whichever is greater
		widget.TextInputOpts.Image(img),

		//Set the font face and size for the widget
		widget.TextInputOpts.Face(&inputFace),

		//Set the colors for the text and caret
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          color.NRGBA{0, 0, 50, 255},
			Disabled:      color.NRGBA{0, 0, 20, 100},
			Caret:         color.NRGBA{0, 0, 50, 255},
			DisabledCaret: color.NRGBA{R: 200, G: 0, B: 0, A: 255},
		}),

		//Set how much padding there is between the edge of the input and the text
		widget.TextInputOpts.Padding(widget.NewInsetsSimple(5)),

		//This text is displayed if the input is empty
		widget.TextInputOpts.Placeholder(lastText),

		//This is called when the user hits the "Enter" key.
		//There are other options that can configure this behavior

		//This is called whenever there is a change to the text
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Changed: ", args.InputText)

		}),
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			hub.Publish(events.ButtonClickedEvent{ButtonText: "Submit"})
		}),
	)

	b3 := LoadOutlineTextButtonSubmitBg("Submit", hub, "")

	textContainer.AddChild(textInput)

	buttonRow.AddChild(textContainer)
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
	window := makeWindow(cont, registry.Config.ResolutionHeight/2+190, 500, 150)
	rfunc := ui.AddWindow(window)
	return rfunc, txtInput
}

func triggerDebugTextInputWindow(lastInput string, ui *ebitenui.UI, hub *tasks.EventHub) (widget.RemoveWindowFunc, *widget.TextInput) {
	cont, txtInput := MakeDebugTextInput(lastInput, hub)
	window := makeWindow(cont, registry.Config.ResolutionHeight/2+100, 400, 130)
	rfunc := ui.AddWindow(window)
	return rfunc, txtInput
}

func TriggerNextDayWindow(ui *ebitenui.UI, hub *tasks.EventHub) widget.RemoveWindowFunc {

	buttonText := []string{"Yes", "Not yet"}

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
				Position: widget.RowLayoutPositionStart,
				Stretch:  false,
			})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))
	face := registry.FontMap["RockSalt"]
	headerLbl := widget.NewText(
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			})),

		widget.TextOpts.Text(header, &face, color.RGBA{R: 0, G: 0, B: 0, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionStart),
	)
	var lines []*widget.Text
	for _, line := range inputText {
		textCon := widget.NewText(
			widget.TextOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position: widget.RowLayoutPositionStart,
				})),

			widget.TextOpts.Text(line, &face, color.RGBA{R: 0, G: 0, B: 0, A: 255}),
			widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
		)
		lines = append(lines, textCon)
	}
	b1 := LoadTextButtonNoBg("X", hub, "closeTutorial")
	btnContainer.AddChild(b1)
	childContainer.AddChild(btnContainer)

	childContainer.AddChild(headerLbl)

	for _, line := range lines {
		childContainer.AddChild(line)
	}

	rootContainer.AddChild(childContainer)
	wind := makeWindow(rootContainer, 100, 1000, 400)

	rfunc := ui.AddWindow(wind)
	return rfunc

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
func TriggerJournalWindow(mag *Journal, ui *ebitenui.UI, hub *tasks.EventHub) widget.RemoveWindowFunc {

	activeContainer := mag.ActiveWindow()

	window := widget.NewWindow(

		widget.WindowOpts.Contents(activeContainer),

		widget.WindowOpts.Modal(),
		widget.WindowOpts.Location(image.Rect(250, 100, ScreenWidth*2-250, ScreenHeight*2-150)),
		widget.WindowOpts.MinSize(ScreenWidth*2, ScreenHeight*2), // Force consistent sizing
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Moved")
		}),
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Resized")
		}),
	)

	fmt.Printf("Window location before UI add: %+v\n", window.GetContainer().GetWidget().Rect)
	removeFunc := ui.AddWindow(window)
	fmt.Printf("Window location after UI add: %+v\n", window.GetContainer().GetWidget().Rect)
	go func() {
		time.Sleep(100 * time.Millisecond) // Let one render cycle happen
		fmt.Printf("After render - Root container rect: %+v\n", activeContainer.GetWidget().Rect)
	}()
	return removeFunc
}
func TriggerMagWindow(mag *Magazine, ui *ebitenui.UI, hub *tasks.EventHub) widget.RemoveWindowFunc {

	activeContainer := mag.ActiveWindow()

	window := widget.NewWindow(

		widget.WindowOpts.Contents(activeContainer),

		widget.WindowOpts.Modal(),
		widget.WindowOpts.Location(image.Rect(250, 100, ScreenWidth*2-250, ScreenHeight*2-150)),
		widget.WindowOpts.MinSize(ScreenWidth*2, ScreenHeight*2), // Force consistent sizing
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Moved")
		}),
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Resized")
		}),
	)

	fmt.Printf("Window location before UI add: %+v\n", window.GetContainer().GetWidget().Rect)
	removeFunc := ui.AddWindow(window)
	fmt.Printf("Window location after UI add: %+v\n", window.GetContainer().GetWidget().Rect)
	go func() {
		time.Sleep(100 * time.Millisecond) // Let one render cycle happen
		fmt.Printf("After render - Root container rect: %+v\n", activeContainer.GetWidget().Rect)
	}()
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
	for _, arg := range args {
		if arg == "maintainState" {
			return
		}
	}
	ev2 := events.WindowClosed{}
	switch m.windowOpen {
	case JournalMagazineState:
		ev2.Window = string(entities.GrandpasJournal)
	case StoreMagazineState:
		ev2.Window = string(entities.Magazine)
	case PHGuess:
		ev2.Window = "phGuess"
	default:
		ev2.Window = "test"

	}
	m.eventHub.Publish(ev2)
	m.windowOpen = None
}

func closeWindow(hub *tasks.EventHub, data *MainMenuData) {
	data.RemoveFunc([]string{})
}

func MainMenuSubs(mainDat *MainMenuData, magazine *Magazine, ui *ebitenui.UI, hub *tasks.EventHub) {

	buttonTextSubs(mainDat, magazine, ui, hub)
	uiSpriteActionSubs(mainDat, magazine, ui, hub)

	hub.Subscribe(events.CloseWindow{}, func(e tasks.Event) {
		switch mainDat.windowOpen {
		case PHGuess:
			return
		case ChooseProp:
			return
		case StoreMagazineState:
			mainDat.RemoveFunc([]string{})
			magazine.activeIndex = Index
		default:
			mainDat.RemoveFunc([]string{})
		}
	})

	hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		/*ev := e.(events.NewDay)
			mainDat.day = ev.Day
			mainDat.dayType = ev.DayType
			if ev.Day == 1 {
				howToPlayText := []string{"There are a few important daily tasks for fish care",
					"1. Monitoring PH levels and Temperature.",
					"  - Each fish species likes different conditions and will \n" +
						"  become stressed and unhealthy if things are too \nfar from their ideal",
					"2. Cleaning up debris.",
					"  - use your skimmer to clean out any debris chunks\n that will be hard to filter",
					"3. Feeding your fish.",
					"  - Don't over feed your fish, wasted food\n will make the tank get dirtier faster"}
				rfunc := TriggerTextWindow(hub, ui, "Grandpa's Journal: GoldFish care 101", howToPlayText)
				mainDat.addRemoveFunc(rfunc)
				mainDat.windowOpen = TutorialWindow
			}

		})*/
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
		case StoreMagazineState, JournalMagazineState:
			MagazineWindowButtonEvent(mainDat, magazine, ui, ev, hub)
		case PHGuess:
			PhGuessWindowButtonEvent(mainDat, ev, hub)
		case GoToBed:
			if ev.ButtonText == "Go to Bed?: Not yet" {
				mainDat.RemoveFunc([]string{})
			}
		case DebugText:
			if ev.ButtonText == "Submit" {
				mainDat.RemoveFunc([]string{})
				sendEv := events.DebugTextEntered{InputText: mainDat.textInputReference.GetText(), For: mainDat.tempParam}
				mainDat.eventHub.Publish(sendEv)
			}
		case TutorialWindow:
			if ev.ButtonText == "closeTutorial" {
				mainDat.RemoveFunc([]string{})
			}
		}
	})
}

func uiSpriteActionSubs(mainDat *MainMenuData, magazine *Magazine, ui *ebitenui.UI, hub *tasks.EventHub) {
	hub.Subscribe(events.UISpriteAction{}, func(e tasks.Event) {
		ev := e.(events.UISpriteAction)
		if ev.UiSprite == "lightSwitch" && ev.UiSpriteAction == "clicked" {
			SendWindowOpened(hub)
			mainDat.windowOpen = GoToBed
			mainDat.addRemoveFunc(TriggerNextDayWindow(ui, hub))
		}

		if ev.UiSprite == "magazine" && ev.UiSpriteAction == "picked up" {
			mainDat.windowOpen = StoreMagazineState
			SendWindowOpened(hub)
			mainDat.addRemoveFunc(TriggerMagWindow(magazine, ui, hub))
		}

		if ev.UiSprite == string(entities.GrandpasJournal) && ev.UiSpriteAction == "picked up" {
			mainDat.windowOpen = JournalMagazineState
			SendWindowOpened(hub)
			mainDat.addRemoveFunc(TriggerJournalWindow(mainDat.journal, ui, hub))
		}

		if ev.UiSprite == "phreader" && ev.UiSpriteAction == "ph reading" {
			SendWindowOpened(hub)
			mainDat.windowOpen = PHGuess
			rfunc, tarea := triggerPHwindow(ui, hub)
			mainDat.textInputReference = tarea
			mainDat.addRemoveFunc(rfunc)
		}

	})
	hub.Subscribe(events.DebugTextInput{}, func(e tasks.Event) {
		ev := e.(events.DebugTextInput)
		SendWindowOpened(hub)
		mainDat.windowOpen = DebugText
		rfunc, tarea := triggerDebugTextInputWindow(ev.LastText, ui, mainDat.eventHub)
		tarea.Focus(true)
		mainDat.textInputReference = tarea
		mainDat.addRemoveFunc(rfunc)
		mainDat.tempParam = ev.For
	})
}

func PhGuessWindowButtonEvent(mainDat *MainMenuData, ev events.ButtonClickedEvent, hub *tasks.EventHub) {
	switch ev.ButtonText {
	case "+0.1":
		curr := mainDat.textInputReference.GetText()
		phVal, err := strconv.ParseFloat(curr, 64)
		if err != nil {
			log.Printf("Error parsing '%s': %v", curr, err) // This will show you exactly what's in the string
		}
		phVal += 0.1
		phValStr := fmt.Sprintf("%.1f", phVal)
		mainDat.textInputReference.SetText(phValStr)
	case "-0.1":
		curr := mainDat.textInputReference.GetText()
		phVal, err := strconv.ParseFloat(curr, 64)
		if err != nil {
			log.Printf("Error parsing '%s': %v", curr, err) // This will show you exactly what's in the string
		}
		phVal -= 0.1
		phValStr := fmt.Sprintf("%.1f", phVal)
		mainDat.textInputReference.SetText(phValStr)
	case "+1.0":
		curr := mainDat.textInputReference.GetText()
		phVal, err := strconv.ParseFloat(curr, 64)
		if err != nil {
			log.Printf("Error parsing '%s': %v", curr, err) // This will show you exactly what's in the string
		}
		phVal += 1.0
		phValStr := fmt.Sprintf("%.1f", phVal)
		mainDat.textInputReference.SetText(phValStr)
	case "-1.0":
		curr := mainDat.textInputReference.GetText()
		phVal, err := strconv.ParseFloat(curr, 64)
		if err != nil {
			log.Printf("Error parsing '%s': %v", curr, err) // This will show you exactly what's in the string
		}
		phVal -= 1.0
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

	mainDat.RemoveFunc([]string{"maintainState"})
	if mainDat.windowOpen == StoreMagazineState {
		mainDat.addRemoveFunc(TriggerMagWindow(magazine, ui, hub))
	}
	if mainDat.windowOpen == JournalMagazineState {
		mainDat.addRemoveFunc(TriggerJournalWindow(mainDat.journal, ui, hub))
	}

	if strings.HasPrefix(ev.ButtonText, "Buy:") {
		closeWindow(hub, mainDat)
		magazine.activeIndex = Index
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
