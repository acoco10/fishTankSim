package ui

import (
	"github.com/acoco10/fishTankWebGame/game/gameEntities"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"image/color"
	"log"
)

type StartMenu struct {
	*ebitenui.UI
	TextInputContainer  *widget.Container
	TextInput           *widget.TextInput
	TextInputButton     *widget.Button
	SelectSpriteOptions map[string]*gameEntities.AnimatedSprite
	SelectedAnimations  map[string]*gameEntities.AnimatedSprite
	SelectSpritesToDraw []*gameEntities.AnimatedSprite
	eventHub            *gameEntities.EventHub
	fishButtons         map[string]*widget.Button
}

func LoadStartMenu(hub *gameEntities.EventHub) (*StartMenu, error) {
	s := StartMenu{}
	s.eventHub = hub

	err := LoadStartMenuUI(&s)
	if err != nil {
		return &s, err
	}

	selectSprites, err := LoadAnimatedSelectSprites()
	if err != nil {
		return nil, err
	}

	s.SelectSpriteOptions = selectSprites

	s.SelectedAnimations, err = loadSelectedAnimations()
	if err != nil {
		return nil, err
	}

	for _, sprite := range selectSprites {
		s.SelectSpritesToDraw = append(s.SelectSpritesToDraw, sprite)
	}

	s.subs()

	return &s, nil
}

func LoadStartMenuUI(startMenu *StartMenu) error {

	headerText := "Pick Your Starter Fish!"

	face, err := LoadFont(54, "nk57")

	if err != nil {
		return err
	}

	rootContainer := widget.NewContainer(

		widget.ContainerOpts.Layout(widget.NewRowLayout(
			//Define how much padding to inset the child content
			widget.RowLayoutOpts.Spacing(100),
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Padding(
				widget.Insets{Right: 0, Left: 0, Top: 30, Bottom: 30},
			),
		),
		),

		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	childContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()))

	pickFishContainer := widget.NewContainer(

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
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			VerticalPosition:   widget.AnchorLayoutPositionEnd,
			HorizontalPosition: widget.AnchorLayoutPositionCenter,
			StretchHorizontal:  false,
			StretchVertical:    false,
		}),
		))

	headerContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 10)),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout()),
	)

	headerLbl := widget.NewText(
		widget.TextOpts.Text(headerText, face, color.RGBA{R: 250, G: 160, B: 0, A: 255}),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				VerticalPosition:   widget.AnchorLayoutPositionStart,
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				StretchHorizontal:  false,
				StretchVertical:    false,
			}),
		),
		widget.TextOpts.Insets(widget.Insets{
			Left:   30,
			Right:  10,
			Top:    50,
			Bottom: 600,
		}),
	)

	textinputContainer, textInputBox, textInputButton, err := NewTextInput(startMenu.eventHub)
	if err != nil {
		return err
	}

	textinputContainer.GetWidget().Disabled = true
	textinputContainer.GetWidget().Visibility = widget.Visibility_Hide

	headerContainer.AddChild(headerLbl)

	b1 := LoadSpriteSelectButton("Goldfish", startMenu.eventHub, 12)
	b2 := LoadSpriteSelectButton("Common Molly", startMenu.eventHub, 12)
	b3 := LoadSpriteSelectButton("Select", startMenu.eventHub, 16)

	fishButtonMap := make(map[string]*widget.Button)

	pickFishContainer.AddChild(
		b1, b2,
	)

	fishButtonMap["Goldfish"] = b1
	fishButtonMap["Common Molly"] = b2
	fishButtonMap["Selected Button"] = b3

	childContainer.AddChild(pickFishContainer)
	rootContainer.AddChild(
		headerContainer, childContainer, textinputContainer, b3)

	b3.GetWidget().Visibility = widget.Visibility_Hide

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	startMenu.TextInput = textInputBox
	startMenu.TextInputContainer = textinputContainer
	startMenu.fishButtons = fishButtonMap
	startMenu.TextInputButton = textInputButton
	startMenu.UI = &ui

	return nil
}

func (s *StartMenu) subs() {
	s.eventHub.Subscribe(gameEntities.ButtonClickedEvent{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.ButtonClickedEvent)
		if ev.ButtonText == "Common Molly" {

			s.HideAndDisableSelectButtons()

			s.fishButtons["Selected Button"].GetWidget().Visibility = widget.Visibility_Show
			s.fishButtons["Selected Button"].Text().Label = "Common Molly"

			s.fishButtons["Selected Button"].Press()

			s.SelectSpritesToDraw = []*gameEntities.AnimatedSprite{}
			s.SelectSpriteOptions["Common Molly"].X -= 70

			copyAs := *s.SelectSpriteOptions["Common Molly"]
			ols := gameEntities.LoadOutlineShader()
			copyAs.LoadShader(ols)
			s.SelectSpritesToDraw = append(s.SelectSpritesToDraw, &copyAs)
		}
	})

	s.eventHub.Subscribe(gameEntities.ButtonClickedEvent{}, func(e gameEntities.Event) {
		ev := e.(gameEntities.ButtonClickedEvent)
		if ev.ButtonText == "Goldfish" {

			s.HideAndDisableSelectButtons()

			s.fishButtons["Selected Button"].GetWidget().Visibility = widget.Visibility_Show
			s.fishButtons["Selected Button"].Text().Label = "Goldfish"

			s.fishButtons["Selected Button"].Press()
			s.SelectSpritesToDraw = []*gameEntities.AnimatedSprite{}
			s.SelectSpriteOptions["Goldfish"].X += 70
			s.SelectSpritesToDraw = append(s.SelectSpritesToDraw, s.SelectSpriteOptions["Goldfish"])

			copyAs := *s.SelectSpriteOptions["Goldfish"]
			ols := gameEntities.LoadOutlineShader()
			copyAs.LoadShader(ols)
			s.SelectSpritesToDraw = append(s.SelectSpritesToDraw, &copyAs)
		}
	})
}

func (s *StartMenu) HideAndDisableSelectButtons() {
	s.fishButtons["Common Molly"].GetWidget().Visibility = widget.Visibility_Hide
	s.fishButtons["Goldfish"].GetWidget().Visibility = widget.Visibility_Hide

	s.fishButtons["Common Molly"].GetWidget().Disabled = true
	s.fishButtons["Goldfish"].GetWidget().Disabled = true
}

func LoadAnimatedSelectSprites() (map[string]*gameEntities.AnimatedSprite, error) {

	fishOptions := make(map[string]*gameEntities.AnimatedSprite)

	fishSprite := gameEntities.LoadFishSprite(gameEntities.Fish, 1)

	//scaledImg := ebiten.NewImage(fishSprite.Img.Bounds().Dx()*4, fishSprite.Img.Bounds().Dy()*4)

	//dopts := &ebiten.DrawImageOptions{}
	//dopts.GeoM.Scale(4, 4)
	//scaledImg.DrawImage(fishSprite.Img, dopts)
	//fishSprite.Img = scaledImg

	//fishSprite.SpriteWidth = fishSprite.SpriteWidth * 4
	//fishSprite.SpriteHeight = fishSprite.SpriteHeight * 4

	fishSprite.X = 355
	fishSprite.Y = 260

	mollyFishSprite := gameEntities.LoadFishSprite(gameEntities.MollyFish, 1)
	mollyFishSprite.Animation.SpeedInTPS = 10

	mollyFishSprite.X = 490
	mollyFishSprite.Y = 260

	fishOptions["Common Molly"] = mollyFishSprite
	fishOptions["Goldfish"] = fishSprite

	return fishOptions, nil
}

func loadSelectedAnimations() (map[string]*gameEntities.AnimatedSprite, error) {
	fishOptions := make(map[string]*gameEntities.AnimatedSprite)

	mollyFishSprite, err := gameEntities.LoadFishSpriteAltAnimations(gameEntities.MollyFish)
	if err != nil {
		log.Fatal(err)
	}
	mollyFishSprite.Animation.SpeedInTPS = 4

	mollyFishSprite.X = 420
	mollyFishSprite.Y = 260

	fishOptions["mollyFish"] = mollyFishSprite

	return fishOptions, nil
}

func (s *StartMenu) PlaySelectedAnimation(selected string) {
	s.SelectSpritesToDraw = []*gameEntities.AnimatedSprite{s.SelectedAnimations[selected]}
}
