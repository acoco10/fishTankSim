package ui

import (
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
)

func LoadMowingUI(eHub *tasks.EventHub) *ebitenui.UI {
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
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				}),
		),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Spacing(20),
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(widget.Insets{Right: 0, Left: 0, Top: 100, Bottom: 0}),
			)),
	)

	button := LoadSubmitButton("Continue", eHub, 16)
	//choreButton.GetWidget().Visibility = widget.Visibility_Hide
	//choreButton.GetWidget().Disabled = true
	//modeButton := LoadSubmitButton("Mode", eHub, 12)

	//button2 := LoadSubmitButton("Mute Music", eHub, 12)
	//button3 := LoadSubmitButton("Mute Sounds", eHub, 12)

	//buttonContainer.AddChild(button2)
	//buttonContainer.AddChild(button3)
	//buttonContainer.AddChild(modeButton)

	//notePad, err := NewTextBlock(eHub, NotePad)
	//if err != nil {

	//notePad.text.SetText("To Do:")

	buttonContainer.AddChild(button)
	rootContainer.AddChild(buttonContainer)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	return &ui

}
