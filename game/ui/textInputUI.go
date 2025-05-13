package ui

import (
	"github.com/acoco10/fishTankWebGame/game/gameEntities"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"image/color"
)

func NewTextInput(ehub *gameEntities.EventHub) (*widget.Container, *widget.TextInput, *widget.Button, error) {

	img, err := loadTextInputImage()

	if err != nil {
		return nil, nil, nil, err
	}

	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		// the container will use a row layout to layout the textinput widgets
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(2),
			widget.GridLayoutOpts.Spacing(20, 20),
			widget.GridLayoutOpts.Padding(widget.Insets{
				Top:    0,
				Left:   0,
				Right:  0,
				Bottom: 100,
			}),
		)))

	face, err := LoadFont(20)
	if err != nil {
		return nil, nil, nil, err
	}

	btnContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 50)),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	textContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 10)),

		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	textInput := widget.NewTextInput(

		widget.TextInputOpts.WidgetOpts(
			//Set the layout information to center the textbox in the parent
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),

			widget.WidgetOpts.MinSize(200, 10),
		),
		//Set the Idle and Disabled background image for the text input
		//If the NineSlice image has a minimum size, the widget will use that or
		// widget.WidgetOpts.MinSize; whichever is greater
		widget.TextInputOpts.Image(img),

		//Set the font face and size for the widget
		widget.TextInputOpts.Face(face),

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
			widget.CaretOpts.Size(face, 2),
		),

		//This text is displayed if the input is empty
		widget.TextInputOpts.Placeholder("Give her a name!"),

		//This is called when the user hits the "Enter" key.
		//There are other options that can configure this behavior
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			ev := gameEntities.SendData{
				DataFor: "Name Input",
				Data:    args.InputText,
			}

			ehub.Publish(ev)
		}),

		//This is called whenever there is a change to the text
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
		}),
	)

	b3 := LoadSubmitButton("Submit", ehub, 16)

	textContainer.AddChild(textInput)
	btnContainer.AddChild(b3)

	rootContainer.AddChild(textContainer)
	rootContainer.AddChild(btnContainer)

	return rootContainer, textInput, b3, nil
}

func loadTextInputImage() (*widget.TextInputImage, error) {
	img, err := gameEntities.LoadImageAssetAsEbitenImage("menuAssets/textInputBox")

	if err != nil {
		return nil, err
	}

	nineSliceImage := eimage.NewNineSlice(img, [3]int{16, img.Bounds().Dx() - 32, 16}, [3]int{16, img.Bounds().Dy() - 16 - 3, 3})

	textInputImg := widget.TextInputImage{
		Idle: nineSliceImage,
	}

	return &textInputImg, nil
}
