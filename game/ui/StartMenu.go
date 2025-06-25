package ui

import (
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/acoco10/fishTankWebGame/shaders"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
	"log"
)

type startMenuState uint8

const (
	title startMenuState = iota
	fishSelected
	propSelected
)

type StartMenu struct {
	*ebitenui.UI
	state               startMenuState
	screenWidth         int
	screenHeight        int
	root                *widget.Container
	TextInputContainer  *widget.Container
	TextInput           *widget.TextInput
	TextInputButton     *widget.Button
	DrawOptions         map[string]drawables.Drawable
	SelectSpritesToDraw []drawables.Drawable
	eventHub            *tasks.EventHub
	fishButtons         map[string]*widget.Button
	propButtons         []*widget.Container
	fishButtonContainer *widget.Container
	propButtonContainer *widget.Container
	selectContainer     *widget.Container
}

func LoadStartMenu(hub *tasks.EventHub, screenWidth int, screenHeight int) (*StartMenu, error) {
	headerFontSize := 54.0

	s := StartMenu{}
	s.eventHub = hub
	s.screenHeight = screenHeight
	s.screenWidth = screenWidth

	err := LoadStartMenuUI(&s, headerFontSize)
	if err != nil {
		return &s, err
	}

	selectSprites, err := loader.LoadAnimatedSelectSprites(screenWidth, screenHeight)

	if err != nil {
		return nil, err
	}

	s.DrawOptions = selectSprites

	for _, sp := range selectSprites {
		s.SelectSpritesToDraw = append(s.SelectSpritesToDraw, sp)
	}

	img, err := loader.LoadImageAssetAsEbitenImage("menuAssets/backButton")
	if err != nil {
		return nil, err
	}

	backSprite := &sprite.Sprite{Img: img, X: float32(screenWidth/2 - 100 - img.Bounds().Dx()), Y: float32(screenHeight/2 - (img.Bounds().Dy())/2)}
	s.DrawOptions["Back"] = backSprite

	s.subs()
	loader.StartScreenCoordinatePositioner(s.screenHeight, s.screenWidth, s.DrawOptions, 12.0, 54)
	return &s, nil
}

func LoadStartMenuUI(startMenu *StartMenu, headerFontSize float64) error {

	headerText := "Pick Your Starter Fish!"

	face, err := util.LoadFont(headerFontSize, "nk57")

	if err != nil {
		return err
	}

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	childContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    true,
				Padding:            widget.Insets{Top: startMenu.screenHeight / 5},
			}),
		),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(20),
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		)))

	pickFishContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		// the container will use a plain color as its background
		// the container will use an anchor layout to layout its single child widget
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			//Define number of columns in the grid
			widget.GridLayoutOpts.Columns(2),
			//Define how much padding to inset the child content
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			//Define how far apart the rows and columns should be
			widget.GridLayoutOpts.Spacing(20, 10),
			// DefaultStretch values will be used when extra columns/rows are used
			// out of the ones defined on the normal Stretch
			widget.GridLayoutOpts.DefaultStretch(false, true),
			//Define how to stretch the rows and columns.
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false}),
		),
		),
	)

	headerContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	headerLbl := widget.NewText(
		widget.TextOpts.Text(headerText, face, color.RGBA{R: 250, G: 160, B: 0, A: 255}),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
		widget.TextOpts.Insets(widget.Insets{
			Left:  30,
			Right: 10,
		}),
	)

	headerContainer.AddChild(headerLbl)

	b1 := LoadSpriteSelectButton("Goldfish", startMenu.eventHub, 16)
	b2 := LoadSpriteSelectButton("Common Molly", startMenu.eventHub, 16)
	b3 := LoadSpriteSelectButton("Select", startMenu.eventHub, 16)

	fishButtonMap := make(map[string]*widget.Button)

	pickFishContainer.AddChild(
		b1, b2,
	)

	propButs := addPickSpriteContainer(2)

	fishButtonMap["Goldfish"] = b1
	fishButtonMap["Common Molly"] = b2
	fishButtonMap["Selected Button"] = b3
	childContainer.AddChild(headerContainer)
	childContainer.AddChild(pickFishContainer)
	childContainer.AddChild(propButs)

	rootContainer.AddChild(
		childContainer)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	startMenu.fishButtons = fishButtonMap
	startMenu.fishButtonContainer = pickFishContainer
	startMenu.propButtonContainer = propButs
	//startMenu.buttonContainer = childContainer
	startMenu.UI = &ui
	startMenu.root = childContainer

	return nil
}

func (s *StartMenu) subs() {

	s.eventHub.Subscribe(events.ButtonEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonEvent)
		if ev.EType == "cursor exited" {
			if ev.ButtonText != "Select" {
				if len(s.SelectSpritesToDraw) > 1 {
					s.DrawOptions[ev.ButtonText].(*sprite.AnimatedSprite).UnLoadShader()
				}
			}
		}
	})

	s.eventHub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		switch ev.ButtonText {
		case "Submit":
			s.TextInputButton.Press()
			s.TextInput.Focus(false)
		case "Back":
			s.Back()
		default:
			s.SpriteSelected(ev.ButtonText)
		}
	})

	s.eventHub.Subscribe(events.ButtonEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonEvent)
		if ev.EType == "cursor entered" {
			if ev.ButtonText != "Select" {
				if s.DrawOptions[ev.ButtonText] != nil {
					s.DrawOptions[ev.ButtonText].(*sprite.AnimatedSprite).LoadShader(registry.ShaderMap["Outline"])
				}
			}
		}
	})
}

func (s *StartMenu) SpriteSelected(tx string) {
	switch s.state {
	case title:
		s.state = fishSelected
		s.fishButtonContainer.RemoveChildren()
		s.fishButtonContainer.AddChild(s.fishButtons["Selected Button"])

		/*s.fishButtons["Selected Button"].GetWidget().Visibility = widget.Visibility_Show*/
		s.fishButtons["Selected Button"].Text().Label = tx
		//s.fishButtons["Selected Button"].Press()

		s.SelectSpritesToDraw = []drawables.Drawable{}

		selectedFish := s.DrawOptions[tx].(*sprite.AnimatedSprite)
		offset := 120 - selectedFish.SpriteWidth*4

		selectedFish.X = float32(s.screenWidth/2 - (selectedFish.SpriteWidth)*2 - offset/2)

		ols := shaders.LoadOutlineShader()
		selectedFish.Shader = ols

		s.SelectSpritesToDraw = append(s.SelectSpritesToDraw, selectedFish, s.DrawOptions["Back"])

		if s.propButtons != nil {
		}

		propButtons, err := propSelectChild(s.eventHub)
		if err != nil {
			log.Fatal("error initiating prop container", err)
		}
		s.propButtons = propButtons
		for _, but := range propButtons {
			s.propButtonContainer.AddChild(but)
		}
		moveBack(s.DrawOptions["Back"], s.state)

	case fishSelected:
		s.propButtonContainer.RemoveChildren()
		s.propButtonContainer.AddChild(s.propButtons[0])
		textinputContainer, _, textInputButton, err := NewTextInput(s.eventHub)
		if err != nil {
			log.Fatal("text input dun fucked up", err)
		}
		s.TextInputContainer = textinputContainer
		s.TextInputButton = textInputButton
		s.root.AddChild(textinputContainer)
		s.state = propSelected
	}

}

func (s *StartMenu) Back() {
	switch s.state {
	case fishSelected:
		s.state = title
		s.propButtonContainer.RemoveChildren()
		s.fishButtonContainer.RemoveChildren()
		s.fishButtonContainer.AddChild(s.fishButtons["Goldfish"])
		s.fishButtonContainer.AddChild(s.fishButtons["Common Molly"])
		loader.StartScreenCoordinatePositioner(s.screenHeight, s.screenWidth, s.DrawOptions, 12.0, 54)
		s.SelectSpritesToDraw = []drawables.Drawable{}
		s.SelectSpritesToDraw = append(s.SelectSpritesToDraw, s.DrawOptions["Common Molly"], s.DrawOptions["Goldfish"])

	case propSelected:

		s.state = fishSelected
		s.propButtonContainer.RemoveChildren()

		for _, but := range s.propButtons {
			s.propButtonContainer.AddChild(but)
		}

		s.TextInputContainer.RemoveChildren()

	}
}

func moveBack(backButton drawables.Drawable, state startMenuState) {
	switch state {
	case propSelected:
		backButton.(*sprite.Sprite).X -= 100
	case fishSelected:

		backButton.(*sprite.Sprite).X -= 100
		backButton.(*sprite.Sprite).Y += 150

	}

}

func (s *StartMenu) ResetSpritePositions(imageWidth, height int, fishOptions map[string]drawables.Drawable) {
	midpoint := image.Point{X: imageWidth / 2, Y: height / 2}
	orderKeys := []string{"Goldfish", "Common Molly"}
	spacing := 20
	yOffset := 20
	i := 0

	for _, key := range orderKeys {
		fish := fishOptions[key].(*sprite.AnimatedSprite)
		minSize := 120
		imgWidth := fish.Img.Bounds().Dx()
		widthAndBuffer := minSize
		offSet := widthAndBuffer - imgWidth

		if i == 0 {
			fish.X = float32(midpoint.X - imgWidth - 20/2 - offSet/2)
		} else {
			fish.X = float32(midpoint.X + spacing/2 + offSet/2)
		}

		fish.Y = float32(midpoint.Y - fish.Img.Bounds().Dy() - yOffset)
		i++
	}
}

func propSelectChild(eventHub *tasks.EventHub) ([]*widget.Container, error) {
	// no prop icons yet
	/*castleImg, err := loader.LoadImageAssetAsEbitenImage("tankProps/castleProp")
	if err != nil {
		return nil, err
	}*/

	b1, err := LoadStackSpriteSelectButton("Castle", ebiten.NewImage(10, 10), 16, eventHub)

	if err != nil {
		return nil, err
	}

	b2, err := LoadStackSpriteSelectButton("Log", ebiten.NewImage(10, 10), 16, eventHub)

	if err != nil {
		return nil, err
	}
	buttons := []*widget.Container{b1, b2}

	return buttons, nil
}

func addPickSpriteContainer(nButtons int) *widget.Container {
	pickPropContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),

		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(nButtons),
			widget.GridLayoutOpts.Padding(widget.NewInsetsSimple(10)),
			widget.GridLayoutOpts.Spacing(20, 10),
			widget.GridLayoutOpts.DefaultStretch(false, true),
			widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false}),
		),
		),
	)

	return pickPropContainer
}
