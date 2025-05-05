package ui

import (
	"fishTankWebGame/assets"
	"fishTankWebGame/game/gameEntities"
	"fmt"
	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

type TextBoxUi struct {
	*ebitenui.UI
	text     *widget.TextArea
	triggerd bool
	eventhub *gameEntities.EventHub
}

func NewTextBlocKMenu(hub *gameEntities.EventHub) (*TextBoxUi, error) {
	t := &TextBoxUi{}
	t.eventhub = hub
	face, err := LoadFont(10)
	if err != nil {
		return nil, err
	}

	img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/statsMenuBackground.png")
	if err != nil {
		return nil, err
	}

	nineSliceImage := eimage.NewNineSlice(img, [3]int{8, 66 - 16, 8}, [3]int{8, 32 - 16, 8})

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(
				widget.Insets{Right: 670, Left: 0, Top: 100, Bottom: 0},
			),
		)),
	)

	text := []string{"fish stats go here"}

	// construct a textarea
	textarea := widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				//Set the layout data for the textarea
				//including a max height to ensure the scroll bar is visible
				widget.WidgetOpts.LayoutData(widget.RowLayoutData{
					Position:  widget.RowLayoutPositionCenter,
					MaxWidth:  200,
					MaxHeight: 600,
				}),
				//Set the minimum size for the widget
				widget.WidgetOpts.MinSize(100, 100),
			),
		),
		widget.TextAreaOpts.ControlWidgetSpacing(2),
		widget.TextAreaOpts.ProcessBBCode(true),
		widget.TextAreaOpts.FontColor(color.White),
		widget.TextAreaOpts.FontFace(face),

		widget.TextAreaOpts.Text(text[0]),
		//Tell the TextArea to show the vertical scrollbar
		//Set padding between edge of the widget and where the text is drawn
		widget.TextAreaOpts.TextPadding(widget.Insets{Right: 10, Left: 10, Top: 10, Bottom: 10}),
		//This sets the background images for the scroll container
		widget.TextAreaOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(&widget.ScrollContainerImage{
				Idle: nineSliceImage,
				Mask: eimage.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
			}),
		),
		//This sets the images to use for the sliders
		widget.TextAreaOpts.SliderOpts(
			widget.SliderOpts.Images(
				// Set the track images
				&widget.SliderTrackImage{
					Idle:  eimage.NewNineSliceColor(color.NRGBA{200, 200, 200, 255}),
					Hover: eimage.NewNineSliceColor(color.NRGBA{200, 200, 200, 255}),
				},
				// Set the handle images
				&widget.ButtonImage{
					Idle:    eimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
					Hover:   eimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
					Pressed: eimage.NewNineSliceColor(color.NRGBA{255, 100, 100, 255}),
				},
			),
		),
	)

	t.text = textarea
	//Add text to the end of the textarea
	//textarea.AppendText("\nLast Row")
	//Add text to the beginning of the textarea
	//textarea.PrependText("First Row\n")
	//Replace the current text with the new value
	//textarea.SetText("New Value!")
	//Retrieve the current value of the text area text
	fmt.Println(textarea.GetText())
	// add the textarea as a child of the container
	rootContainer.AddChild(textarea)

	ui := ebitenui.UI{
		Container: rootContainer,
	}

	t.UI = &ui

	t.subs()

	return t, nil
}

func (t *TextBoxUi) subs() {
	t.eventhub.Subscribe(gameEntities.SendData{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.SendData)
		if ev.DataFor == "statsMenu" {
			t.UpdateTextArea(ev.Data)
			t.Trigger()
		}
	})
}

func (t *TextBoxUi) Draw(screen *ebiten.Image) {
	if t.triggerd {
		t.UI.Draw(screen)
	}
}

func (t *TextBoxUi) RequestData(target any) {
	ev := &gameEntities.RequestData{}
	ev.RequestFor = target
	t.eventhub.Publish(ev)
}

func (t *TextBoxUi) UpdateTextArea(text string) {
	t.text.SetText(text)
}

func (t *TextBoxUi) Trigger() {
	t.triggerd = true
}
