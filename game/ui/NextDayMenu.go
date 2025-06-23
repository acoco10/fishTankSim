package ui

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
)

type NextDayMenu struct {
	*ebitenui.UI
	eventHub  *tasks.EventHub
	Triggered bool
}

func LoadNextDayMenuUI(hub *tasks.EventHub) (*widget.Window, error) {

	headerText := "Go To Bed?"

	face, err := util.LoadFont(24, "nk57")

	if err != nil {
		return nil, err
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

	ButtonContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.GridLayoutOpts.Spacing(20, 10),
			widget.GridLayoutOpts.DefaultStretch(false, true),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false}),
		),
		),
	)

	headerLbl := widget.NewText(
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			})),

		widget.TextOpts.Text(headerText, face, color.RGBA{R: 250, G: 160, B: 0, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionStart),
		widget.TextOpts.Insets(widget.Insets{}),
	)

	//headerContainer.AddChild(headerLbl)

	b1 := LoadSubmitButton("Go To Bed", hub, 12)
	b2 := LoadSubmitButton("Vibe awhile Longer", hub, 12)

	ButtonContainer.AddChild(
		b1, b2,
	)

	childContainer.AddChild(headerLbl)
	childContainer.AddChild(ButtonContainer)

	rootContainer.AddChild(childContainer)

	// construct the UI

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
		widget.WindowOpts.Location(image.Rect(100, 100, 500, 500)),
	)
	return window, nil
}

func (nd *NextDayMenu) Update() {
	nd.UI.Update()
}

func (nd *NextDayMenu) Draw(screen *ebiten.Image) {
	if nd.Triggered {
		nd.UI.Draw(screen)
	}
}

func (nd *NextDayMenu) subs() {

	nd.eventHub.Subscribe(events.UISpriteAction{}, func(e tasks.Event) {
		ev := e.(events.UISpriteAction)
		if ev.UiSprite == "pillow" && ev.UiSpriteAction == "clicked" {
			nd.Triggered = true
		}
	})

	nd.eventHub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		if ev.ButtonText == "Vibe awhile Longer" {
			nd.Triggered = false
		}
	})

}
