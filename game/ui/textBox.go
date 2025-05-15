package ui

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/acoco10/fishTankWebGame/game/gameEntities"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	text2 "github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"log"
)

var n int

type TextBoxType uint

const (
	StatsMenu TextBoxType = iota
	NotePad
)

type TextBoxUi struct {
	*widget.Container
	text     *widget.TextArea
	triggerd bool
	eventhub *gameEntities.EventHub
}

func LoadBackgroundImageForTextInput(boxType TextBoxType) *widget.ScrollContainerImage {
	switch boxType {
	case StatsMenu:
		img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/statsMenuBackground.png")
		if err != nil {
			log.Printf("background image for text output container not loading due to: %q", err.Error())
		}

		nineSliceImage := eimage.NewNineSlice(img, [3]int{8, 66 - 16, 8}, [3]int{8, 32 - 16, 8})

		wImage := widget.ScrollContainerImage{
			Idle: nineSliceImage,
			Mask: eimage.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
		}

		return &wImage
	case NotePad:
		/*img, _, err := ebitenutil.NewImageFromFileSystem(assets.ImagesDir, "images/menuAssets/whiteBoardMain.png")
		if err != nil {
			log.Printf("background image for text output container not loading due to: %q", err.Error())
		}*/

		//nineSliceImage := eimage.NewNineSlice(img, [3]int{32, 180 - 64, 32}, [3]int{32, 135 - 64, 32})

		wImage := widget.ScrollContainerImage{
			Idle: eimage.NewNineSliceColor(color.NRGBA{100, 100, 100, 10}),
			Mask: eimage.NewNineSliceColor(color.NRGBA{100, 100, 100, 255}),
		}

		return &wImage
	}
	return nil
}

func NewTextBlock(hub *gameEntities.EventHub, tp TextBoxType) (*TextBoxUi, error) {
	t := &TextBoxUi{}

	t.eventhub = hub

	img := LoadBackgroundImageForTextInput(tp)
	rootContainer, textArea, err := NewTextBlockContainer(hub, img, tp)
	if err != nil {
		return nil, err
	}
	t.text = textArea
	t.Container = rootContainer
	t.subs(tp)

	return t, nil
}

func LoadLayoutData(tp TextBoxType) widget.AnchorLayoutData {
	switch tp {
	case StatsMenu:
		layoutData := widget.AnchorLayoutData{
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
			HorizontalPosition: widget.AnchorLayoutPositionStart,
			StretchHorizontal:  false,
			StretchVertical:    false,
		}
		return layoutData
	case NotePad:
		layoutData := widget.AnchorLayoutData{
			VerticalPosition:   widget.AnchorLayoutPositionCenter,
			HorizontalPosition: widget.AnchorLayoutPositionEnd,
			StretchHorizontal:  false,
			StretchVertical:    false,
		}
		return layoutData
	}
	layoutData := widget.AnchorLayoutData{
		VerticalPosition:   widget.AnchorLayoutPositionCenter,
		HorizontalPosition: widget.AnchorLayoutPositionEnd,
		StretchHorizontal:  false,
		StretchVertical:    false,
	}
	return layoutData
}

func LoadPadding(tp TextBoxType) (widget.Insets, widget.Insets) {
	var w widget.Insets
	var w2 widget.Insets
	switch tp {
	case StatsMenu:
		w = widget.Insets{Right: 0, Left: 50, Top: 0, Bottom: 200}
		w2 = widget.Insets{Right: 10, Left: 10, Top: 10, Bottom: 10}
	case NotePad:
		w = widget.Insets{Right: 10, Left: 0, Top: 0, Bottom: 0}
		w2 = widget.Insets{Right: 10, Left: 10, Top: 10, Bottom: 0}
	}

	return w, w2
}

func LoadMinSize(tp TextBoxType) (int, int) {
	var w, h int
	switch tp {
	case StatsMenu:
		w = 160
		h = 120
	case NotePad:
		w = 180
		h = 135
	}

	return w, h
}

func LoadFontByType(tp TextBoxType) (text2.Face, color.Color, error) {
	var face text2.Face
	var clr color.Color
	switch tp {
	case StatsMenu:
		lFace, err := LoadFont(10, "nk57")
		if err != nil {
			return face, clr, err
		}
		face = lFace
		clr = color.White
	case NotePad:
		lFace, err := LoadFont(16, "rockSalt")
		if err != nil {
			return face, clr, err
		}
		face = lFace
		clr = color.Black
	}
	return face, clr, nil
}

func NewTextBlockContainer(hub *gameEntities.EventHub, backGroundImg *widget.ScrollContainerImage, tp TextBoxType) (*widget.Container, *widget.TextArea, error) {
	t := &TextBoxUi{}

	t.eventhub = hub

	face, textClr, err := LoadFontByType(tp)

	if err != nil {
		return nil, nil, err
	}
	w, h := LoadMinSize(tp)
	layoutData := LoadLayoutData(tp)
	padding, textPadding := LoadPadding(tp)
	/*layoutDataRoot := widget.AnchorLayoutData{
		VerticalPosition:   widget.AnchorLayoutPositionCenter,
		HorizontalPosition: widget.AnchorLayoutPositionEnd,
		StretchHorizontal:  false,
		StretchVertical:    false,
	}*/

	// construct a new container that serves as the root of the UI hierarchy
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(layoutData),
		),

		widget.ContainerOpts.Layout(
			widget.NewAnchorLayout(widget.AnchorLayoutOpts.Padding(padding))),
		// the container will use an anchor layout to layout its single child widget
	)

	text := []string{"fish stats go here"}

	// construct a textarea
	textarea := widget.NewTextArea(
		widget.TextAreaOpts.ContainerOpts(
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(layoutData),
				//Set the layout data for the textarea
				//including a max height to ensure the scroll bar is visible
				//Set the minimum size for the widget
				widget.WidgetOpts.MinSize(w, h),
			),
		),

		widget.TextAreaOpts.ControlWidgetSpacing(2),
		widget.TextAreaOpts.ProcessBBCode(true),
		widget.TextAreaOpts.FontColor(textClr),
		widget.TextAreaOpts.FontFace(face),

		widget.TextAreaOpts.Text(text[0]),
		//Tell the TextArea to show the vertical scrollbar
		//Set padding between edge of the widget and where the text is drawn
		widget.TextAreaOpts.TextPadding(textPadding),
		//This sets the background images for the scroll container
		widget.TextAreaOpts.ScrollContainerOpts(
			widget.ScrollContainerOpts.Image(backGroundImg),
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
	return rootContainer, textarea, nil

}

func (t *TextBoxUi) subs(tp TextBoxType) {
	switch tp {
	case StatsMenu:
		t.eventhub.Subscribe(gameEntities.SendData{}, func(e gameEntities.Event) {
			ev := e.(gameEntities.SendData)
			if ev.DataFor == "statsMenu" {
				switch ev.Data {
				case "fish deselect":
					t.text.GetWidget().Visibility = widget.Visibility_Hide
				default:
					t.UpdateTextArea(ev.Data)
					t.text.GetWidget().Visibility = widget.Visibility_Show
				}
			}
		})

	case NotePad:
		t.eventhub.Subscribe(gameEntities.SendData{}, func(e gameEntities.Event) {
			ev := e.(gameEntities.SendData)
			if ev.DataFor == "whiteBoard" {
				t.AppendTextArea(ev.Data)
			}
		})

		t.eventhub.Subscribe(gameEntities.TaskCompleted{}, func(e gameEntities.Event) {
			//ev := e.(gameEntities.TaskCompleted)
		})
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

func (t *TextBoxUi) AppendTextArea(text string) {
	t.text.AppendText("\n" + text)
}

func (t *TextBoxUi) Trigger() {
	t.triggerd = true
}
