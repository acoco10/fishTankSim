package ui

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
	"strings"
)

type PurchaseType int8

const (
	Decoration PurchaseType = iota
	Plant
	Item
	Fish
)

type Magazine struct {
	pages       map[string]*widget.Container
	pageIndex   int
	activeIndex string
	pageOrder   []string
	CostMap     map[string]float64
}

// StoreItem represents an item that can be purchased
type StoreItem struct {
	Name        string
	Description string
	Image       *ebiten.Image
	Price       float64
	EventName   string // What event to publish when purchased
}

// TextContent represents content for a text page
type TextContent struct {
	Title   string
	Content string
	Image   *ebiten.Image // Optional image
}

func (m *Magazine) ActiveWindow() *widget.Container {
	return m.pages[m.activeIndex]
}

// CreateSimpleMagazine creates a magazine with just the pages you specify
func CreateSimpleMagazine(hub *tasks.EventHub) *Magazine {
	mag := &Magazine{
		pages:       make(map[string]*widget.Container),
		activeIndex: "info",
	}

	// Set up event subscriptions
	return mag
}

// AddTextPage adds a simple text page to the magazine
func (m *Magazine) AddTextPage(page *widget.Container, pageName string) error {
	m.pages[pageName] = page
	m.pageOrder = append(m.pageOrder, pageName)
	return nil
}

// AddStorePage adds a store page with purchasable items
func (m *Magazine) AddStorePage(items []StoreItem, hub *tasks.EventHub, pageName string) error {
	page, err := createStorePage(pageName, items, hub)
	if err != nil {
		return err
	}

	m.pages[pageName] = page
	m.pageOrder = append(m.pageOrder, pageName)
	return nil
}

// AddIndexPage adds a navigation page with buttons to other pages
func (m *Magazine) AddIndexPage(buttons []string, hub *tasks.EventHub, pageName string) {
	page := createIndexPage(buttons, hub)
	m.pageOrder = append([]string{pageName}, m.pageOrder...)
	m.pages[pageName] = page
}

// Helper function to create the basic magazine page structure
func createMagazinePageBase() (*widget.Container, *widget.Container, *widget.Container) {
	magNineSlice, flippedMagNineSlice, err := LoadMagNineSlice()
	if err != nil {
		return nil, nil, nil
	}

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(ScreenHeight*2, ScreenWidth*2)),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Spacing(0, 0),
				widget.GridLayoutOpts.Padding(&widget.Insets{}),
				widget.GridLayoutOpts.Stretch([]bool{true, true}, []bool{true, true}),
			),
		),
	)

	leftPage := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(800, 900)),
		widget.ContainerOpts.BackgroundImage(magNineSlice),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(50),
				widget.RowLayoutOpts.Padding(&widget.Insets{Top: 50, Left: 10, Right: 10}),
			),
		),
	)

	rightPage := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(800, 900)),
		widget.ContainerOpts.BackgroundImage(flippedMagNineSlice),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(50),
				widget.RowLayoutOpts.Padding(&widget.Insets{Top: 50, Left: 10, Right: 10}),
			),
		),
	)

	rootContainer.AddChild(leftPage, rightPage)
	return rootContainer, leftPage, rightPage
}

// createTextPage creates a page with title, text content, and optional image
func createTextPages(contents []TextContent, hub *tasks.EventHub) *widget.Container {

	root, leftPage, rightPage := createMagazinePageBase()
	if contents[0].Image != nil {
		createImagePage(leftPage, contents)
	} else {
		createTextPage(contents[0], leftPage)
	}
	if len(contents) > 1 {
		createTextPage(contents[1], rightPage)
	}

	button := LoadSubmitButton("Back", hub, "")
	leftPage.AddChild(button)

	button2 := LoadSubmitButton("Forward", hub, "")
	rightPage.AddChild(button2)

	return root

}

func createImagePage(container *widget.Container, contents []TextContent) {
	textContainer := widget.NewContainer(widget.ContainerOpts.WidgetOpts(
		widget.WidgetOpts.LayoutData(widget.RowLayoutData{Position: widget.RowLayoutPositionCenter})),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(10),
			),
		),
	)

	face, _ := registry.FontMap["RockSalt"]
	if contents[0].Title != "" {
		titleText := widget.NewText(
			widget.TextOpts.Text(contents[0].Title, &face, color.RGBA{R: 60, G: 160, B: 200, A: 255}),
			widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
			widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Position: widget.RowLayoutPositionCenter})))
		textContainer.AddChild(titleText)
	}

	for _, content := range contents {
		if content.Content != "" {
			contentText := widget.NewText(
				widget.TextOpts.Text(content.Content, &face, color.Black),
				widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
				widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.GridLayoutData{HorizontalPosition: widget.GridLayoutPositionStart})),
			)
			params := make(map[string]any)
			params["imageScale"] = 2
			params["minWidth"] = 64
			params["minHeight"] = 64

			gc := LoadGraphic(content.Image, params, contentText)
			textContainer.AddChild(gc)
		}
	}

	container.AddChild(textContainer)

}

func createTextPage(content TextContent, container *widget.Container) {
	textContainer := widget.NewContainer(widget.ContainerOpts.WidgetOpts(
		widget.WidgetOpts.LayoutData(widget.RowLayoutData{Position: widget.RowLayoutPositionCenter})),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(10),
			),
		),
	)

	face, _ := registry.FontMap["RockSalt"]
	if content.Title != "" {

		titleText := widget.NewText(
			widget.TextOpts.Text(content.Title, &face, color.RGBA{R: 60, G: 160, B: 200, A: 255}),
			widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
			widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Position: widget.RowLayoutPositionCenter})))

		textContainer.AddChild(titleText)
	}

	if content.Content != "" {
		contentText := widget.NewText(
			widget.TextOpts.Text(content.Content, &face, color.Black),
			widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
			widget.TextOpts.Padding(&widget.Insets{Top: 10, Left: 25, Right: 10, Bottom: 30}),
			widget.TextOpts.Padding(&widget.Insets{Top: 10, Left: 25, Right: 10, Bottom: 30}),
			widget.TextOpts.MaxWidth(700),
		)
		textContainer.AddChild(contentText)
	}

	container.AddChild(textContainer)
}

// createStorePage creates a page with purchasable items
func createStorePage(title string, items []StoreItem, hub *tasks.EventHub) (*widget.Container, error) {
	root, leftPage, rightPage := createMagazinePageBase()

	face := registry.FontMap["RockSalt"]
	titleText := widget.NewText(
		widget.TextOpts.Text(title, &face, color.RGBA{R: 60, G: 160, B: 200, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Position: widget.RowLayoutPositionCenter})))
	leftPage.AddChild(titleText)

	fillerText := widget.NewText(
		widget.TextOpts.Text("", &face, color.RGBA{R: 60, G: 160, B: 200, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionCenter),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Position: widget.RowLayoutPositionCenter})))
	rightPage.AddChild(fillerText)

	// Distribute items between left and right pages
	leftItems := items[:len(items)/2]
	rightItems := items[len(items)/2:]

	// Add items to left page
	for _, item := range leftItems {
		itemContainer := createStoreItemWidget(item, hub)
		leftPage.AddChild(itemContainer)
	}

	button := LoadSubmitButton("Back", hub, "")
	leftPage.AddChild(button)

	// Add items to right page
	for _, item := range rightItems {
		itemContainer := createStoreItemWidget(item, hub)
		rightPage.AddChild(itemContainer)

	}

	button2 := LoadSubmitButton("Forward", hub, "")
	rightPage.AddChild(button2)

	return root, nil
}

// createStoreItemWidget creates a widget for a single store item
func createStoreItemWidget(item StoreItem, hub *tasks.EventHub) *widget.Container {
	container := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
				widget.RowLayoutOpts.Spacing(10),
			),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.RowLayoutData{Position: widget.RowLayoutPositionCenter}),
		),
	)

	// Create sprite button if image provided

	params := make(map[string]any)
	params["imageScale"] = 2
	params["minWidth"] = 150
	params["minHeight"] = 150
	if item.Image != nil {
		println("making image button for item:", item.Name)
		button, _ := LoadStackSpriteSelectButton(
			item.Name,
			item.Image,
			hub,
			params,
		)
		container.AddChild(button)
	}

	// Create text container with description and price
	textContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(5),
			),
		),
	)

	// Description
	face, _ := util.LoadFont(12, "nk57")
	descText := widget.NewText(
		widget.TextOpts.Text(item.Description, &face, color.Black),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
		widget.TextOpts.MaxWidth(300),
	)

	// Price
	priceText := widget.NewText(
		widget.TextOpts.Text(fmt.Sprintf("Price: $%.2f", item.Price), &face, color.RGBA{R: 60, G: 160, B: 50, A: 255}),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
	)

	// Buy button
	eventName := item.EventName
	if eventName == "" {
		eventName = "Buy: " + item.Name
	}
	buyButton := LoadSubmitButton("Buy", hub, eventName)

	textContainer.AddChild(descText, priceText, buyButton)
	container.AddChild(textContainer)

	return container
}

// createIndexPage creates a navigation page with buttons
func createIndexPage(buttonTexts []string, hub *tasks.EventHub) *widget.Container {
	root, leftPage, _ := createMagazinePageBase()

	// Add title
	face, _ := registry.FontMap["RockSalt"]
	titleText := widget.NewText(
		widget.TextOpts.Text("Index", &face, color.RGBA{R: 60, G: 160, B: 200, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionEnd),
	)
	leftPage.AddChild(titleText)

	// Add buttons
	buttonContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.RowLayoutData{Position: widget.RowLayoutPositionCenter}),
		),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(10),
			),
		),
	)

	for _, buttonText := range buttonTexts {
		ttl := strings.ToTitle(buttonText)
		button := LoadSubmitButton(ttl, hub, buttonText)
		buttonContainer.AddChild(button)
	}

	leftPage.AddChild(buttonContainer)

	return root
}

// setupMagazineEvents handles all magazine-related events
func setupMagazineEvents(mag *Magazine, hub *tasks.EventHub) {
	hub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		_, pageExists := mag.pages[ev.ButtonText]
		if pageExists {
			mag.activeIndex = ev.ButtonText
			for i, pgName := range mag.pageOrder {
				if pgName == ev.ButtonText {
					mag.pageIndex = i
				}
			}
		} else {
			println("mag page:", ev.ButtonText, "does not exist")
		}
		// Handle purchase events
		if strings.HasPrefix(ev.ButtonText, "Buy:") {
			f := func(c rune) bool {
				return c == ':'
			}
			fields := strings.FieldsFunc(ev.ButtonText, f)
			var buyType PurchaseType
			switch strings.Trim(fields[1], " ") {
			case "item":
				buyType = Item
			case "plant":
				buyType = Plant
			case "decoration":
				buyType = Decoration
			case "fish":
				buyType = Fish
			}

			if cost, exists := mag.CostMap[fields[2]]; exists {
				purchaseEvent := events.BuyAttempt{
					Cost:     cost,
					Item:     fields[2],
					ItemType: uint8(buyType),
				}
				hub.Publish(purchaseEvent)
			}
		}
	})
}

func defaultMagSubs(mag *Magazine, hub *tasks.EventHub) {
	hub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		if ev.ButtonText == "Back" {
			if mag.pageIndex > 0 {
				mag.pageIndex--
				mag.activeIndex = mag.pageOrder[mag.pageIndex]
			}
		}
		if ev.ButtonText == "Forward" {
			if mag.pageIndex < len(mag.pageOrder)-1 {
				mag.pageIndex++
				mag.activeIndex = mag.pageOrder[mag.pageIndex]
			}

		}
	})
}

func LoadCollectionPage() {

}

// Example usage functions:

func CreateFishMagazine(hub *tasks.EventHub) (*Magazine, error) {
	// Setup cost map
	costMap := map[string]float64{
		"kirbensis":  2.0,
		"mollyFish":  1.0,
		"goldFish":   1.0,
		"guppy":      1.0,
		"angelFish":  5,
		"phBoost":    0.25,
		"phReduce":   0.25,
		"plantPack1": 1.0,
		"fertilizer": .50,
		"zenFriend":  12,
		"zenBridge":  5,
		"castle":     25,
		"log":        5,
		"hotRock":    .75,
		"coolRock":   .75,
	}

	/*	styleCostMap := map[string]float64{
		"gravelecolor" : 5,
		"discoball"
	}*/

	mag := CreateSimpleMagazine(hub)
	mag.CostMap = costMap

	// Add index page

	// Load fish images
	fishImages, err := LoadFishSprites()
	if err != nil {
		return nil, err
	}

	// Create store items
	fishDescriptions, fishNames := LoadFishDescriptions()
	fish := make([]StoreItem, 5)

	itemImgs := loadStoreImgs()

	for i := 0; i < 5; i++ {
		fish[i] = StoreItem{
			Name:        fishNames[i],
			Description: fishDescriptions[i],
			Image:       fishImages[i],
			Price:       costMap[fishNames[i]],
			EventName:   "Buy: fish:" + fishNames[i],
		}
	}

	err = mag.AddStorePage(fish, hub, "fish")
	if err != nil {
		return nil, err
	}

	plantItems := []StoreItem{
		{
			Name:        "Plant Pack 1",
			Description: "A pack of 3 that may contain ferns, grass or leafy plants with a chance at a rare version of each",
			Price:       costMap["plantPack1"],
			Image:       itemImgs["plantPack1"],
			EventName:   "Buy: plant:plantPack1",
		},
		{
			Name:        "Fertilizer",
			Description: "Double chance at a rare on your next plant",
			Price:       costMap["fertilizer"],
			Image:       itemImgs["fertilizer"],
			EventName:   "Buy: item:fertilizer",
		},
	}

	err = mag.AddStorePage(plantItems, hub, "Plants")
	if err != nil {
		return nil, err
	}

	decorations := []StoreItem{
		{
			Name:        "Zen Friend",
			Description: "Add some calm and tranquility to your tank.\n\n-Decreases PH by 2.0, \n\n-Decreases fish stress by 1n\n\n-zen set",
			Price:       costMap["zenFriend"],
			Image:       itemImgs["zenFriend"],
			EventName:   "Buy: decoration:zenFriend",
		},
		{
			Name:        "Zen Bridge",
			Description: "Add some calm and tranquility to your tank.\n\n -Increases PH by 1.0, \n\n-Decreases fish stress by 0.25\n\n-zen set",
			Price:       costMap["zenBridge"],
			Image:       itemImgs["zenBridge"],
			EventName:   "Buy: decoration:zenBridge",
		},
		{
			Name:        "Log",
			Description: "A cheap nature booster.\n\n -Decrease PH by 2.0, \n\n-Increase environment by 1.0",
			Price:       costMap["log"],
			Image:       itemImgs["log"],
			EventName:   "Buy: decoration:log",
		},
		{
			Name:        "Castle",
			Description: "A decoration fit for a king.\n\n -Increases PH by 2.0, \n\n-Cave structure",
			Price:       costMap["castle"],
			Image:       itemImgs["castle"],
			EventName:   "Buy: decoration:castle",
		},
	}

	err = mag.AddStorePage(decorations, hub, "Decorations")

	mag.activeIndex = "info"
	// Add accessories store page
	accessoryItems := []StoreItem{
		{
			Name:        "PH Boost",
			Description: "Raises tank PH levels",
			Image:       itemImgs["phaidb"],
			Price:       0.25,
			EventName:   "Buy: item:phBoost",
		},
		{
			Name:        "PH Reduce",
			Description: "Lowers tank PH levels",
			Image:       itemImgs["phaidr"],
			Price:       0.25,
			EventName:   "Buy: item:phReduce",
		},
		{
			Name:        "Cool Rock",
			Description: "Cools down your tank",
			Image:       itemImgs["coolRock"],
			Price:       0.75,
			EventName:   "Buy: decoration:coolRock",
		},
		{
			Name:        "Hot Rock",
			Description: "Heats up your tank",
			Image:       itemImgs["hotRock"],
			Price:       0.75,
			EventName:   "Buy: decoration:hotRock",
		},
	}

	mag.AddStorePage(accessoryItems, hub, "items")

	var keys []string
	for key, _ := range mag.pages {
		keys = append(keys, key)
	}
	mag.AddIndexPage(keys, hub, "info")

	setupMagazineEvents(mag, hub)
	defaultMagSubs(mag, hub)

	return mag, nil
}

// Keep your existing helper functions
func LoadMagNineSlice() (*eimage.NineSlice, *eimage.NineSlice, error) {
	bground, err := util.LoadImageAssetAsEbitenImage("uiSprites/magazineLeft")
	if err != nil {
		return nil, nil, err
	}

	flipImg, err := util.LoadImageAssetAsEbitenImage("uiSprites/magazineRight")
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

func LoadFishSprites() ([12]*ebiten.Image, error) {
	fishList := []entities.FishList{entities.Guppy, entities.Kirbensis, entities.GoldFish, entities.MollyFish, entities.AngelFish}
	imgs := [12]*ebiten.Image{}
	for i, fish := range fishList {
		fishSprite, err := entities.LoadFishAnimations(fish, 2)
		if err != nil {
			log.Fatal(err)
		}
		imgs[i] = fishSprite["swimming"].GetFirstFrameAsStaticImage()
	}

	return imgs, nil
}

func LoadFishDescriptions() ([5]string, [5]string) {
	// Keep your existing implementation
	descriptionMap := [5]string{
		"Guppies are Hardy fish that come in a variety of vibrant colors.",
		"An exotically patterned fish that prefer cave-like structure.",
		"Goldfish are one of the earliest " +
			"fish to be kept as pets. They are easy to care for and flexible.",
		"Molly fish are a friendly, active fish. " +
			"They will swim right up to the glass" +
			"when they are ready to eat.",
		"Angel fish need warm temperatures and acidic temperatures. They are famous for their unique shape and beautiful stripe patterns",
	}

	names := [5]string{"Guppy", "Kirbensis", "GoldFish", "MollyFish", "AngelFish"}
	return descriptionMap, names
}

func loadStoreImgs() map[string]*ebiten.Image {
	return util.LoadDirectoryImages("images/storeAssets")
}
