package ui

import (
	"github.com/acoco10/fishTankWebGame/game/drawables"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"
	"image"
	"log"
)

const (
	ScreenWidth  = 960
	ScreenHeight = 540
)

const (
	WindowSizeWidth  = 1080
	WindowSizeHeight = 1920
)

type startMenuState uint8

const (
	title startMenuState = iota
	fishSelected
	nameInput
)

type StartMenu struct {
	*ebitenui.UI
	state               startMenuState
	screenWidth         int
	screenHeight        int
	root                *widget.Container
	TextInputContainer  *widget.Container
	TextInput           *widget.TextInput
	DrawOptions         map[string]drawables.Drawable
	SelectSpritesToDraw []drawables.Drawable
	eventHub            *tasks.EventHub
	fishButtons         map[string]*widget.Container
	fishButtonContainer *widget.Container
	selectContainer     *widget.Container
}

func LoadStartMenu(hub *tasks.EventHub, resolutionScaling float64) (*StartMenu, error) {
	headerFontSize := 64.0

	s := StartMenu{}
	s.eventHub = hub
	s.screenHeight = ScreenWidth
	s.screenWidth = ScreenHeight

	err := LoadStartMenuUI(&s, headerFontSize, resolutionScaling)
	if err != nil {
		return &s, err
	}

	s.subs()
	loader.StartScreenCoordinatePositioner(s.screenHeight, s.screenWidth, s.DrawOptions, 12.0, 54)
	return &s, nil
}

func LoadStartMenuUI(startMenu *StartMenu, headerFontSize float64, resolutionScalar float64) error {

	headerText := "Pick Your Starter Fish!"

	face, err := util.LoadFont(headerFontSize, "reglisseOutline") //white center text

	face2, err := util.LoadFont(headerFontSize, "reglisseOutlined") //black outline

	if err != nil {
		println("error loading new font")
		return err
	}

	borders := eimage.NewBorder(4, 4, 4, 4, colornames.Lightgoldenrodyellow)
	nineSliceImage := eimage.NewAdvancedNineSliceColor(colornames.Lightcoral, borders)

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
				Padding: widget.Insets{
					Top:    int(float64(registry.Config.ResolutionHeight) * resolutionScalar / 10),
					Bottom: int(float64(registry.Config.ResolutionHeight) * resolutionScalar / 10),
					Left:   int(float64(registry.Config.ResolutionWidth) * resolutionScalar / 10),
					Right:  int(float64(registry.Config.ResolutionWidth) * resolutionScalar / 10)}, // base container distance from the top of the screen

			}),
		),
		widget.ContainerOpts.BackgroundImage(nineSliceImage),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			//MAIN VARIABLE FOR CHANGING ROW SPACING
			widget.RowLayoutOpts.Spacing(int(30*resolutionScalar)),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)),
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
			//onlt one row so row spacing second input doesnt really matter
			widget.GridLayoutOpts.Spacing(int(20*resolutionScalar), int(10*resolutionScalar)),
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
			widget.WidgetOpts.LayoutData(widget.RowLayoutOpts)),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	headerLbl := widget.NewText(
		widget.TextOpts.Text(headerText, face, colornames.Aliceblue),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
		widget.TextOpts.Insets(widget.Insets{50, 50, 50, 50}),
	)

	headerLblOutline := widget.NewText(
		widget.TextOpts.Text(headerText, face2, colornames.Black),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
		widget.TextOpts.Insets(widget.Insets{50, 50, 50, 50}),
	)

	headerContainer.AddChild(headerLbl)
	headerContainer.AddChild(headerLblOutline)

	goldFish, err := entities.LoadFishSprite(entities.Fish, 2)
	if err != nil {
		return err
	}

	mollyFish, err := entities.LoadFishSprite(entities.MollyFish, 2)
	if err != nil {
		return err
	}

	goldFishImg := goldFish["swimming"].GetFirstFrameAsStaticImage()
	mFishImg := mollyFish["swimming"].GetFirstFrameAsStaticImage()

	b1, err := LoadStackSpriteSelectButton("Goldfish", goldFishImg, float64(12*resolutionScalar), startMenu.eventHub, 4.0)
	if err != nil {
		return err
	}

	b2, err := LoadStackSpriteSelectButton("Common Molly", mFishImg, float64(12*resolutionScalar), startMenu.eventHub, 4.0)
	if err != nil {
		return err
	}

	fishButtonMap := make(map[string]*widget.Container)

	pickFishContainer.AddChild(
		b1, b2,
	)

	fishButtonMap["Goldfish"] = b1
	fishButtonMap["Common Molly"] = b2

	childContainer.AddChild(headerContainer)
	childContainer.AddChild(pickFishContainer)

	rootContainer.AddChild(
		childContainer)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	startMenu.fishButtons = fishButtonMap
	startMenu.fishButtonContainer = pickFishContainer
	//startMenu.buttonContainer = childContainer
	startMenu.UI = &ui
	startMenu.root = childContainer

	return nil
}

func (s *StartMenu) subs() {

	s.eventHub.Subscribe(events.ButtonEvent{}, func(e tasks.Event) {
		//ev := e.(events.ButtonEvent)

	})

	s.eventHub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		if ev.ButtonText == "Common Molly" || ev.ButtonText == "Goldfish" {
			s.SpriteSelected(ev.ButtonText)
		}
	})

	s.eventHub.Subscribe(events.ButtonEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonEvent)
		if ev.EType == "cursor entered" {
			if ev.ButtonText != "Select" {
				sp, ok := s.DrawOptions[ev.ButtonText].(*sprite.AnimatedSprite)
				if ok {
					sp.LoadShader(registry.ShaderMap["Outline"])
				}
			}
		}
	})
}

func (s *StartMenu) SpriteSelected(tx string) {

	s.state = fishSelected
	s.fishButtonContainer.RemoveChildren()

	s.fishButtonContainer.AddChild(s.fishButtons[tx])

	//moveBack(s.DrawOptions["Back"], s.state)

	textinputContainer, textInput, _, err := NewTextInput(s.eventHub, "Give her a name!")
	if err != nil {
		log.Fatal("text input dun fucked up", err)
	}

	s.TextInputContainer = textinputContainer
	s.root.AddChild(textinputContainer)
	textInput.Focus(true)
}

func (s *StartMenu) Back() {
	switch s.state {
	case fishSelected:
		s.state = title
		s.fishButtonContainer.RemoveChildren()
		s.fishButtonContainer.AddChild(s.fishButtons["Goldfish"])
		s.fishButtonContainer.AddChild(s.fishButtons["Common Molly"])
		loader.StartScreenCoordinatePositioner(s.screenHeight, s.screenWidth, s.DrawOptions, 12.0, 54)
		s.SelectSpritesToDraw = []drawables.Drawable{}
		s.SelectSpritesToDraw = append(s.SelectSpritesToDraw, s.DrawOptions["Common Molly"], s.DrawOptions["Goldfish"])
	}
}

func moveBack(backButton drawables.Drawable, state startMenuState) {
	switch state {
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
