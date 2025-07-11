package ui

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
	"log"
	"strings"
)

type Magazine struct {
	triggered     bool
	pages         []*widget.Container
	activeIndex   int
	buttonGraphic *graphics.SpriteGraphic
	fish          [12]*ebiten.Image
}

func (m *Magazine) ActiveWindow() *widget.Container {
	return m.pages[m.activeIndex]
}

func (m *Magazine) Trigger() {
	m.triggered = true
}

func LoadMagazineUiMenu(eHub *tasks.EventHub, screenWidth int, screenHeight int) (*Magazine, error) {
	bground, err := loader.LoadImageAssetAsEbitenImage("uiSprites/magazineAlt")
	if err != nil {
		return nil, err
	}

	buttonGraphicImg, err := loader.LoadImageAssetAsEbitenImage("menuAssets/arrowButton")
	if err != nil {
		return nil, err
	}

	b := bground.Bounds()
	x := float32(screenWidth-b.Dx()) / 2
	y := float32(screenHeight-b.Dy()) / 2

	//s := sprite.Sprite{Img: bground, X: x, Y: y}
	buttonSprite := sprite.Sprite{Img: buttonGraphicImg, X: x + float32(b.Dx()-10), Y: y + float32(b.Dy()-10)}
	buttonGraphic := graphics.NewFadeInSprite(buttonSprite)

	indexPage, err := LoadMagazineIndexPage(eHub, b)
	if err != nil {
		return nil, err
	}

	magUI := Magazine{}
	//magUI.background = &s
	magUI.activeIndex = 0
	magUI.buttonGraphic = buttonGraphic

	fish, err := LoadFishSprites()
	if err != nil {
		log.Fatal("Fish catalogue image not found:", err)
	}
	magUI.fish = fish

	fishPage, err := LoadFishPages(eHub, magUI.fish)
	if err != nil {
		return nil, err
	}

	infoPage1, err := LoadInfoPages(eHub)

	magUI.pages = append(magUI.pages, indexPage, fishPage, infoPage1)

	MagSubscriptions(&magUI, eHub)

	return &magUI, nil
}

func LoadFishSprites() ([12]*ebiten.Image, error) {

	fish := [12]*ebiten.Image{}

	kirbensis, err := loader.LoadImageAssetAsEbitenImage("staticFish/kirbensis2")
	if err != nil {
		return [12]*ebiten.Image{}, err
	}

	guppy, err := loader.LoadImageAssetAsEbitenImage("staticFish/guppy2")
	if err != nil {
		return [12]*ebiten.Image{}, err
	}

	goldFish, err := loader.LoadFishSprite(entities.Fish, 2)
	if err != nil {
		return [12]*ebiten.Image{}, err
	}

	mollyFish, err := loader.LoadFishSprite(entities.MollyFish, 2)
	if err != nil {
		return [12]*ebiten.Image{}, err
	}

	fish[0] = kirbensis
	fish[1] = guppy
	fish[2] = goldFish.GetFirstFrameAsStaticImage()
	fish[3] = mollyFish.GetFirstFrameAsStaticImage()

	return fish, nil
}

func LoadFishDescriptions() ([4]string, [4]string) {
	descriptionMap := [4]string{}
	descriptionMap[0] = "Guppies are Hardy fish that comes in a variety of vibrant colors. They prefer warmer temperatures. Guppies are social and prefer 2-3 friends."
	descriptionMap[1] = "Kirbensis are easy to breed if they have a cave-like structure. Be cautious housing with aggressive species, very territorial, especially when mating"
	descriptionMap[2] = "One of the earliest fish to be kept as pets and specifically bred. They were first kept in China over 1000 years ago"
	descriptionMap[3] = "Molly fish are a friendly, active fish that come in lots of colors, They are very social and can do well in a tank with more than one species"

	names := [4]string{"Guppy", "Kirbensis", "Fish", "MollyFish"}
	return descriptionMap, names
}

func (m *Magazine) Update() {
	if m.triggered {
		m.pages[m.activeIndex].Update()
		m.buttonGraphic.Update()
	}
}

func LoadMagNineSlice() (*eimage.NineSlice, *eimage.NineSlice, error) {
	bground, err := loader.LoadImageAssetAsEbitenImage("uiSprites/magazineLeft")
	if err != nil {
		return nil, nil, err
	}

	flipImg, err := loader.LoadImageAssetAsEbitenImage("uiSprites/magazineRight")
	if err != nil {
		return nil, nil, err
	}

	magNineSlice := eimage.NewNineSlice(
		bground, [3]int{32, bground.Bounds().Dx() - 64, 32},
		[3]int{32, bground.Bounds().Dy() - 64, 32})

	flipNineSlice := eimage.NewNineSlice(
		flipImg, [3]int{32, bground.Bounds().Dx() - 64, 32},
		[3]int{32, bground.Bounds().Dy() - 64, 32})

	return magNineSlice, flipNineSlice, nil

}

func LoadMagazineIndexPage(eHub *tasks.EventHub, b image.Rectangle) (*widget.Container, error) {

	magNineSlice, rightNineSlice, err := LoadMagNineSlice()
	if err != nil {
		return nil, err
	}

	headerText := "Index"

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Spacing(0, 40),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:   20,
					Right: 50,
					Left:  50,
				},
				),
			)),
	)

	leftContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(400, 500),
		),
		widget.ContainerOpts.BackgroundImage(magNineSlice),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top: 20},
			),
			widget.RowLayoutOpts.Spacing(20),
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		),
		),
	)

	face, err := util.LoadFont(24, "nk57")

	if err != nil {
		return nil, err
	}

	headerLbl := widget.NewText(
		widget.TextOpts.Text(headerText, face, color.RGBA{R: 60, G: 160, B: 200, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		widget.TextOpts.Insets(widget.Insets{
			Top: 20,
		}),
	)

	buttonContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.RowLayoutData{
					Position: widget.RowLayoutPositionCenter,
				}),
		),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Spacing(20),
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Padding(widget.Insets{Top: 20}),
			)),
	)

	rightContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(400, 500),
		),
		widget.ContainerOpts.BackgroundImage(rightNineSlice),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top: 20},
			),
			widget.RowLayoutOpts.Spacing(20),
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		),
		),
	)

	button := LoadSubmitButton("Info", eHub, 16)
	button2 := LoadSubmitButton("Fish", eHub, 16)
	button3 := LoadSubmitButton("Tank Upgrades", eHub, 16)
	button4 := LoadSubmitButton("Accessories", eHub, 16)

	buttonContainer.AddChild(button)
	buttonContainer.AddChild(button2)
	buttonContainer.AddChild(button3)
	buttonContainer.AddChild(button4)

	leftContainer.AddChild(headerLbl)
	leftContainer.AddChild(buttonContainer)

	rootContainer.AddChild(leftContainer)
	rootContainer.AddChild(rightContainer)

	// construct the UI

	return rootContainer, nil
}

func LoadInfoPages(eHub *tasks.EventHub) (*widget.Container, error) {

	headerText := "GoldFish"
	face := registry.FontMap["nk57_12"]

	magNineSlice, flippedMagNineSlice, err := LoadMagNineSlice()
	if err != nil {
		return nil, err
	}

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Spacing(0, 0),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:   20,
					Right: 50,
					Left:  50,
				},
				),
			)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.GridLayoutData{
					HorizontalPosition: widget.GridLayoutPositionCenter,
					VerticalPosition:   widget.GridLayoutPositionCenter,
				}),
		),
	)

	headerLbl := widget.NewText(
		widget.TextOpts.Text(headerText, face, color.RGBA{R: 60, G: 160, B: 200, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
			}),
		),
		widget.TextOpts.Insets(widget.Insets{
			Top: 20,
		}),
	)

	leftPage := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(400, 500),
		),
		widget.ContainerOpts.BackgroundImage(magNineSlice),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.GridLayoutData{
					HorizontalPosition: widget.GridLayoutPositionCenter,
					VerticalPosition:   widget.GridLayoutPositionCenter,
				}),
		),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(10),
				widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
			),
		),
	)

	rightPage := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(400, 500),
		),
		widget.ContainerOpts.BackgroundImage(flippedMagNineSlice),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.GridLayoutData{
					HorizontalPosition: widget.GridLayoutPositionCenter,
					VerticalPosition:   widget.GridLayoutPositionCenter,
				}),
		),

		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(10),
				widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
			),
		),
	)

	leftPage.AddChild(headerLbl)

	rootContainer.AddChild(leftPage)
	rootContainer.AddChild(rightPage)

	// construct the UI

	return rootContainer, nil
}

func LoadFishPages(eHub *tasks.EventHub, fishImages [12]*ebiten.Image) (*widget.Container, error) {

	magNineSlice, flippedMagNineSlice, err := LoadMagNineSlice()
	if err != nil {
		return nil, err
	}

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Spacing(0, 0),
				widget.GridLayoutOpts.Padding(widget.Insets{
					Top:   20,
					Right: 50,
					Left:  50,
				},
				),
			)),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.GridLayoutData{
					HorizontalPosition: widget.GridLayoutPositionCenter,
					VerticalPosition:   widget.GridLayoutPositionCenter,
				}),
		),
	)

	leftPage := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(400, 500),
		),
		widget.ContainerOpts.BackgroundImage(magNineSlice),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.GridLayoutData{
					HorizontalPosition: widget.GridLayoutPositionCenter,
					VerticalPosition:   widget.GridLayoutPositionCenter,
				}),
		),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(10),
				widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
			),
		),
	)

	rightPage := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(400, 500),
		),
		widget.ContainerOpts.BackgroundImage(flippedMagNineSlice),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.GridLayoutData{
					HorizontalPosition: widget.GridLayoutPositionCenter,
					VerticalPosition:   widget.GridLayoutPositionCenter,
				}),
		),

		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(10),
				widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(30)),
			),
		),
	)

	fishDescriptions, names := LoadFishDescriptions()

	for i, fish := range fishImages[0:4] {

		button, err := LoadStackSpriteSelectButton(names[i], fish, 16, eHub)
		if err != nil {
			return nil, err
		}

		container := LoadStackedButtonWithText(button, fishDescriptions[i], eHub, "Buy: "+names[i])
		if i < 3 {
			leftPage.AddChild(container)
		}
		if i >= 3 {
			rightPage.AddChild(container)
		}
	}

	rootContainer.AddChild(leftPage)
	rootContainer.AddChild(rightPage)

	// construct the UI

	return rootContainer, nil
}

func MagSubscriptions(magUi *Magazine, eHub *tasks.EventHub) {
	eHub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		println(ev.ButtonText, "button event received")
		switch ev.ButtonText {
		//button text = event published cases
		case "Fish":
			magUi.activeIndex = 1
		case "Info":
			magUi.activeIndex = 2
		}
		// text processing for buy events
		if strings.HasPrefix(ev.ButtonText, "Buy:") {
			// Extract the part after "Buy"

			itemName := strings.TrimSpace(ev.ButtonText[len("Buy:"):])
			itemName = util.LowCase(itemName)

			pev := events.BuyAttempt{
				Name: itemName,
				Cost: 1,
				Item: "fish",
			}
			eHub.Publish(pev)
		}

	})

}
