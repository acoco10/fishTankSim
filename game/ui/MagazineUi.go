package ui

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/loaders"
	"github.com/acoco10/fishTankWebGame/game/sprite"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/ebitenui/ebitenui"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"log"
)

type Magazine struct {
	pages         []*ebitenui.UI
	activeIndex   int
	background    *sprite.Sprite
	buttonGraphic *graphics.SpriteGraphic
	fish          map[string]*ebiten.Image
}

func LoadMagazineUiMenu(eHub *tasks.EventHub, screenWidth int, screenHeight int) (*Magazine, error) {
	bground, err := loaders.LoadImageAssetAsEbitenImage("uiSprites/magazineAlt")
	if err != nil {
		return nil, err
	}

	buttonGraphicImg, err := loaders.LoadImageAssetAsEbitenImage("menuAssets/arrowButton")
	if err != nil {
		return nil, err
	}

	b := bground.Bounds()
	x := float32(screenWidth-b.Dx()) / 2
	y := float32(screenHeight-b.Dy()) / 2

	s := sprite.Sprite{Img: bground, X: x, Y: y}
	buttonSprite := sprite.Sprite{Img: buttonGraphicImg, X: x + float32(b.Dx()-10), Y: y + float32(b.Dy()-10)}
	buttonGraphic := graphics.NewFadeInSprite(buttonSprite)

	indexPage, err := LoadMagazineIndexPage(eHub, b)
	if err != nil {
		return nil, err
	}

	magUI := Magazine{}
	magUI.background = &s
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

	kirbensis, err := loaders.LoadImageAssetAsEbitenImage("staticFish/kirbensis2")
	if err != nil {
		return nil, err
	}

	guppy, err := loaders.LoadImageAssetAsEbitenImage("staticFish/guppy2")
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

func (m *Magazine) Draw(screen *ebiten.Image) {
	m.pages[m.activeIndex].Draw(screen)
	m.buttonGraphic.Draw(screen)

	//400, 500

	//opts := &ebiten.DrawImageOptions{}

	/*	if m.activeIndex == 1 {
		for _, fish := range m.fish {
			b := fish.Bounds()
			x := float64(400/3 - b.Dx() + 50 + 130)
			y := float64(500/3 - 50 - b.Dy())
			opts.GeoM.Scale(2, 2)
			opts.GeoM.Translate(x, y)
			screen.DrawImage(fish, opts)
			opts.GeoM.Reset()
		}
	}*/
}

func (m *Magazine) Update() {
	m.pages[m.activeIndex].Update()
	m.buttonGraphic.Update()
}

func LoadMagazineIndexPage(eHub *tasks.EventHub, b image.Rectangle) (*ebitenui.UI, error) {
	bground, err := loaders.LoadImageAssetAsEbitenImage("uiSprites/magazineAlt")
	if err != nil {
		return nil, err
	}

	magNineSlice := eimage.NewNineSlice(bground, [3]int{9, bground.Bounds().Dx() - 18, 9}, [3]int{8, 9, 10})
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
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	return &ui, nil
}

func LoadFishPages(eHub *tasks.EventHub, fishImgMap map[string]*ebiten.Image) (*ebitenui.UI, error) {

	bground, err := loaders.LoadImageAssetAsEbitenImage("uiSprites/magazineAlt")
	if err != nil {
		return nil, err
	}
	magNineSlice := eimage.NewNineSlice(bground, [3]int{9, bground.Bounds().Dx() - 18, 9}, [3]int{8, 9, 10})

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

	fishDescriptions := LoadFishDescriptions()

	for key, fish := range fishImgMap {

		button, err := LoadStackSpriteSelectButton(key, fish, 16, eHub)
		if err != nil {
			return nil, err
		}

		container := LoadStackedButtonWithText(button, fishDescriptions[key], eHub)
		leftPage.AddChild(container)
	}

	rootContainer.AddChild(leftPage)
	rootContainer.AddChild(rightPage)

	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}

	return &ui, nil
}

func MagSubscriptions(magUi *Magazine, eHub *tasks.EventHub) {
	eHub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		println(ev.ButtonText, "button event received")
		switch ev.ButtonText {
		case "Fish":
			magUi.activeIndex = 1
		}
	})

	eHub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		switch ev.ButtonText {
		case "Kirbensis":
			block, err := InitStoreTextBlock(250, 50, "Kirbensis", eHub)
			if err != nil {
				return
			}
			magUi.pages[magUi.activeIndex].AddWindow(block)
		}
	})

}

func InitStoreTextBlock(x, y int, fishName string, eHub *tasks.EventHub) (*widget.Window, error) {

	windowContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(colornames.Red)),
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Padding(widget.Insets{
				Top: 20},
			),
			widget.RowLayoutOpts.Spacing(20),
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
		),
		),
	)

	infoBlock, err := NewTextBlock(eHub, StoreMenu)
	if err != nil {
		return nil, fmt.Errorf("erorr inititaing store info block: %q", err)
	}

	buyButton := LoadSubmitButton("Buy", eHub, 10)

	windowContainer.AddChild(infoBlock)
	windowContainer.AddChild(buyButton)

	window := widget.NewWindow(
		//Set the main contents of the window
		widget.WindowOpts.Contents(windowContainer),
		//Set the titlebar for the window (Optional)
		//Set the window above everything else and block input elsewhere
		widget.WindowOpts.Modal(),
		//Set how to close the window. CLICK_OUT will close the window when clicking anywhere
		//that is not a part of the window object
		widget.WindowOpts.CloseMode(widget.CLICK_OUT),
		//Indicates that the window is draggable. It must have a TitleBar for this to work
		widget.WindowOpts.Draggable(),
		//Set the window resizeable
		widget.WindowOpts.Resizeable(),
		//Set the minimum size the window can be
		widget.WindowOpts.MinSize(200, 100),
		//Set the maximum size a window can be
		widget.WindowOpts.MaxSize(300, 300),
		//Set the callback that triggers when a move is complete
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Moved")
		}),
		//Set the callback that triggers when a resize is complete
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Window Resized")
		}),
		widget.WindowOpts.Location(image.Rect(x, y, x+200, y+200)),
	)
	return window, nil
}
