package ui

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/loader"
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
	triggered   bool
	pages       []*widget.Container
	activeIndex int
	//background    *sprite.Sprite
	buttonGraphic *graphics.SpriteGraphic
	fish          map[string]*ebiten.Image
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

	magUI.pages = append(magUI.pages, indexPage, fishPage)

	MagSubscriptions(&magUI, eHub)

	return &magUI, nil
}

func LoadFishSprites() (map[string]*ebiten.Image, error) {

	fish := make(map[string]*ebiten.Image)

	kirbensis, err := loader.LoadImageAssetAsEbitenImage("staticFish/kirbensis2")
	if err != nil {
		return nil, err
	}

	guppy, err := loader.LoadImageAssetAsEbitenImage("staticFish/guppy2")
	if err != nil {
		return nil, err
	}

	fish["Kirbensis"] = kirbensis
	fish["Guppy"] = guppy

	return fish, nil
}

func LoadFishDescriptions() map[string]string {
	descriptionMap := make(map[string]string)
	descriptionMap["Guppy"] = "Guppies are Hardy fish that comes in a variety of vibrant colors. They prefer warmer temperatures. Guppies are social and prefer 2-3 friends."
	descriptionMap["Kirbensis"] = "Kirbensis are easy to breed if they have a cave-like structure. Be cautious housing with aggressive species, very territorial, especially when mating"

	return descriptionMap
}

func (m *Magazine) Update() {
	if m.triggered {
		m.pages[m.activeIndex].Update()
		m.buttonGraphic.Update()
	}
}

func LoadMagNineSlice() (*eimage.NineSlice, *eimage.NineSlice, error) {
	bground, err := loader.LoadImageAssetAsEbitenImage("uiSprites/magazineAlt")
	if err != nil {
		return nil, nil, err
	}

	magNineSlice := eimage.NewNineSlice(
		bground, [3]int{32, bground.Bounds().Dx() - 64, 32},
		[3]int{32, bground.Bounds().Dy() - 64, 32})

	flipImg := ebiten.NewImage(bground.Bounds().Dx(), bground.Bounds().Dy())
	flipOpts := &ebiten.DrawImageOptions{}

	flipOpts.GeoM.Scale(-1, 1)

	flipImg.DrawImage(bground, flipOpts)

	flipNineSlice := eimage.NewNineSlice(
		flipImg, [3]int{32, bground.Bounds().Dx() - 64, 32},
		[3]int{32, bground.Bounds().Dy() - 64, 32})

	return magNineSlice, flipNineSlice, nil

}

func LoadMagazineIndexPage(eHub *tasks.EventHub, b image.Rectangle) (*widget.Container, error) {

	magNineSlice, _, err := LoadMagNineSlice()
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

	button := LoadSubmitButton("Tips and Tricks", eHub, 16)
	button2 := LoadSubmitButton("Fish", eHub, 12)
	button3 := LoadSubmitButton("Tank Upgrades", eHub, 12)
	button4 := LoadSubmitButton("Accessories", eHub, 12)

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

func LoadFishPages(eHub *tasks.EventHub, fishImgMap map[string]*ebiten.Image) (*widget.Container, error) {

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

	fishDescriptions := LoadFishDescriptions()

	for key, fish := range fishImgMap {

		button, err := LoadStackSpriteSelectButton(key, fish, 16, eHub)
		if err != nil {
			return nil, err
		}

		container := LoadStackedButtonWithText(button, fishDescriptions[key], eHub, "Buy: "+key)
		leftPage.AddChild(container)
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
