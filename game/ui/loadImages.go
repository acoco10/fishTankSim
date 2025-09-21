package ui

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

func loadBackButtonImage() (*widget.ButtonImage, error) {

	img, err := util.LoadImageAssetAsEbitenImage("menuAssets/backButton")

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

func loadHighlightSubmitButtonImage() *widget.ButtonImage {

	img, _ := util.LoadImageAssetAsEbitenImage("menuAssets/submitButtonHighlight")
	img2, _ := util.LoadImageAssetAsEbitenImage("menuAssets/submitButtonDisabled")

	imgClicked, _ := util.LoadImageAssetAsEbitenImage("menuAssets/submitButtonAlt")

	nineSliceImage := eimage.NewNineSlice(img, [3]int{9, img.Bounds().Dx() - 18, 9}, [3]int{8, 9, 10})

	nineSliceImageClicked := eimage.NewNineSlice(imgClicked, [3]int{9, img.Bounds().Dx() - 18, 9}, [3]int{8, 9, 10})

	nineSliceDisabled := eimage.NewNineSlice(img2, [3]int{9, img.Bounds().Dx() - 18, 9}, [3]int{8, 9, 10})

	idle := nineSliceImage

	hover := nineSliceImage

	pressed := nineSliceImageClicked

	disabled := nineSliceDisabled

	return &widget.ButtonImage{
		Idle:         idle,
		Hover:        hover,
		Pressed:      pressed,
		PressedHover: pressed,
		Disabled:     disabled,
	}
}

func loadSubmitButtonImage() *widget.ButtonImage {

	img, _ := util.LoadImageAssetAsEbitenImage("menuAssets/submitButton3")
	img2, _ := util.LoadImageAssetAsEbitenImage("menuAssets/submitButtonDisabled")

	imgClicked, _ := util.LoadImageAssetAsEbitenImage("menuAssets/submitButtonAlt")

	nineSliceImage := eimage.NewNineSlice(img, [3]int{12, img.Bounds().Dx() - 24, 12}, [3]int{8, 10, 10})

	nineSliceImageClicked := eimage.NewNineSlice(imgClicked, [3]int{12, img.Bounds().Dx() - 24, 12}, [3]int{8, 10, 10})

	nineSliceDisabled := eimage.NewNineSlice(img2, [3]int{9, img.Bounds().Dx() - 18, 9}, [3]int{8, 9, 10})

	idle := nineSliceImage

	hover := nineSliceImage

	pressed := nineSliceImageClicked

	disabled := nineSliceDisabled

	return &widget.ButtonImage{
		Idle:         idle,
		Hover:        hover,
		Pressed:      pressed,
		PressedHover: pressed,
		Disabled:     disabled,
	}
}

func LoadStackSpriteSelectButton(buttonText string, fishImg *ebiten.Image, hub *tasks.EventHub, params map[string]any) (*widget.Container, error) {

	face := registry.FontMap["RockSalt_18"]

	imgScale := params["imageScale"].(int)

	imgForTransform := ebiten.NewImage(params["minWidth"].(int), params["minHeight"].(int))

	dopts := &ebiten.DrawImageOptions{}
	dopts.GeoM.Scale(float64(imgScale), float64(imgScale))
	height := fishImg.Bounds().Dy() * imgScale
	width := fishImg.Bounds().Dx() * imgScale

	dopts.GeoM.Translate(float64(params["minWidth"].(int)/2-width/2), float64(params["minHeight"].(int)/2-height/2))
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
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(params["minWidth"].(int), params["minHeight"].(int))),
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

func LoadGraphic(image *ebiten.Image, params map[string]any, text *widget.Text) *widget.Container {

	imgScale := params["imageScale"].(int)

	imgForTransform := ebiten.NewImage(params["minWidth"].(int), params["minHeight"].(int))

	dopts := &ebiten.DrawImageOptions{}
	dopts.GeoM.Scale(float64(imgScale), float64(imgScale))
	height := image.Bounds().Dy() * imgScale
	width := image.Bounds().Dx() * imgScale

	dopts.GeoM.Translate(float64(params["minWidth"].(int)/2-width/2), float64(params["minHeight"].(int)/2-height/2))
	imgForTransform.DrawImage(image, dopts)

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Spacing(10, 0),
				widget.GridLayoutOpts.Padding(&widget.Insets{}),
				widget.GridLayoutOpts.Stretch([]bool{true, true}, []bool{true, true}),
			),
		),
	)

	graph := widget.NewGraphic(
		widget.GraphicOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.GridLayoutData{HorizontalPosition: widget.GridLayoutPositionCenter})),
		widget.GraphicOpts.Images(&widget.GraphicImage{
			Idle:     imgForTransform,
			Disabled: imgForTransform,
		},
		),
	)

	rootContainer.AddChild(graph)
	rootContainer.AddChild(text)

	return rootContainer
}

func loadSpriteSelectButtonImage(t string) (*widget.ButtonImage, error) {

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

func loadOptionsMenuInputImage() (*eimage.NineSlice, error) {
	img, err := util.LoadImageAssetAsEbitenImage("menuAssets/opMenu")

	if err != nil {
		return nil, err
	}

	nineSliceImage := eimage.NewNineSlice(img,
		[3]int{11, img.Bounds().Dx() - 22, 11}, // left, middle, right
		[3]int{11, img.Bounds().Dy() - 22, 11})
	return nineSliceImage, nil
}

func loadClearButtonImage() *widget.ButtonImage {

	nineSliceImage := eimage.NewNineSliceColor(color.RGBA{0, 0, 0, 0})

	return &widget.ButtonImage{
		Idle:         nineSliceImage,
		Hover:        nineSliceImage,
		Pressed:      nineSliceImage,
		PressedHover: nineSliceImage,
		Disabled:     nineSliceImage,
	}
}
