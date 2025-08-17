package ui

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
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

func loadSubmitButtonImage() *widget.ButtonImage {

	img, _ := util.LoadImageAssetAsEbitenImage("menuAssets/submitButton3")

	imgClicked, _ := util.LoadImageAssetAsEbitenImage("menuAssets/submitButtonAlt")

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
		Disabled:     pressed,
	}
}

func LoadStackSpriteSelectButton(buttonText string, fishImg *ebiten.Image, fontSize float64, hub *tasks.EventHub, imgScale float64) (*widget.Container, error) {

	face, err := util.LoadFont(fontSize, "nk57")
	if err != nil {
		log.Fatal(err)
	}

	imgForTransform := ebiten.NewImage(fishImg.Bounds().Dx()*int(imgScale), fishImg.Bounds().Dx()*int(imgScale))

	dopts := &ebiten.DrawImageOptions{}
	dopts.GeoM.Scale(imgScale, imgScale)

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
		widget.ButtonOpts.WidgetOpts(widget.WidgetOpts.MinSize(250, 250)),
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
