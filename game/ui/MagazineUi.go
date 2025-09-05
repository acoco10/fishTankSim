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

type Magazine struct {
	pages       map[string]*widget.Container
	activeIndex string
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
	return nil
}

// AddStorePage adds a store page with purchasable items
func (m *Magazine) AddStorePage(items []StoreItem, hub *tasks.EventHub, pageName string) error {
	page, err := createStorePage(items, hub)
	if err != nil {
		return err
	}

	m.pages[pageName] = page
	return nil
}

// AddIndexPage adds a navigation page with buttons to other pages
func (m *Magazine) AddIndexPage(buttons []string, hub *tasks.EventHub, pageName string) {
	page := createIndexPage(buttons, hub)
	m.pages[pageName] = page
}

// Helper function to create the basic magazine page structure
func createMagazinePageBase() (*widget.Container, *widget.Container, *widget.Container) {
	magNineSlice, flippedMagNineSlice, err := LoadMagNineSlice()
	if err != nil {
		return nil, nil, nil
	}

	rootContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(1200, 500)),
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
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(600, 500)),
		widget.ContainerOpts.BackgroundImage(magNineSlice),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(10),
				widget.RowLayoutOpts.Padding(&widget.Insets{Top: 30}),
			),
		),
	)

	rightPage := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(400, 500)),
		widget.ContainerOpts.BackgroundImage(flippedMagNineSlice),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(10),
				widget.RowLayoutOpts.Padding(&widget.Insets{Top: 30}),
			),
		),
	)

	rootContainer.AddChild(leftPage, rightPage)
	return rootContainer, leftPage, rightPage
}

// createTextPage creates a page with title, text content, and optional image
func createTextPages(contents []TextContent, hub *tasks.EventHub) *widget.Container {

	root, leftPage, rightPage := createMagazinePageBase()
	createTextPage(contents[0], leftPage)
	if len(contents) > 1 {
		createTextPage(contents[1], rightPage)
	}

	button := LoadSubmitButton("Back", hub, "")
	leftPage.AddChild(button)

	button2 := LoadSubmitButton("Forward", hub, "")
	rightPage.AddChild(button2)

	return root

}

func createTextPage(content TextContent, container *widget.Container) {

	textContainer := widget.NewContainer(widget.ContainerOpts.WidgetOpts(
		widget.WidgetOpts.LayoutData(widget.RowLayoutData{Position: widget.RowLayoutPositionStart})),
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
		)
		textContainer.AddChild(titleText)
	}

	if content.Content != "" {
		contentText := widget.NewText(
			widget.TextOpts.Text(content.Content, &face, color.Black),
			widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
			widget.TextOpts.Padding(&widget.Insets{Top: 10, Left: 25, Right: 10, Bottom: 30}),
			widget.TextOpts.Padding(&widget.Insets{Top: 10, Left: 25, Right: 10, Bottom: 30}),
			widget.TextOpts.MaxWidth(450),
		)
		textContainer.AddChild(contentText)
	}

	// ... your existing code ...

	// Force same layout behavior as store pages

	// Add some debug after a render cycle

	container.AddChild(textContainer)
}

// createStorePage creates a page with purchasable items
func createStorePage(items []StoreItem, hub *tasks.EventHub) (*widget.Container, error) {
	root, leftPage, rightPage := createMagazinePageBase()

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
		} else {
			println("mag page:", ev.ButtonText, "does not exist")
		}
		// Handle page navigation

		// Handle purchase events
		if strings.HasPrefix(ev.ButtonText, "Buy:") {
			itemName := strings.TrimSpace(ev.ButtonText[len("Buy:"):])
			itemName = util.LowCase(itemName)

			if cost, exists := mag.CostMap[itemName]; exists {
				purchaseEvent := events.BuyAttempt{
					Cost: cost,
					Item: itemName,
				}
				hub.Publish(purchaseEvent)
			}
		}

		// Handle direct item purchases (like ph+ and ph-)
		if cost, exists := mag.CostMap[ev.ButtonText]; exists {
			purchaseEvent := events.BuyAttempt{
				Item: ev.ButtonText,
				Cost: cost,
			}
			hub.Publish(purchaseEvent)
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
		"mollyfish":  1.0,
		"goldFish":   1.0,
		"guppy":      1.0,
		"ph+":        0.25,
		"ph-":        0.25,
		"plantpack1": 1.0,
		"fertilizer": .50,
	}

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
	fish := make([]StoreItem, 4)

	for i := 0; i < 4; i++ {
		fish[i] = StoreItem{
			Name:        fishNames[i],
			Description: fishDescriptions[i],
			Image:       fishImages[i],
			Price:       costMap[util.LowCase(fishNames[i])],
			EventName:   "Buy: " + fishNames[i],
		}
	}

	err = mag.AddStorePage(fish, hub, "fish")
	if err != nil {
		return nil, err
	}

	plantItems := []StoreItem{
		{
			Name:        "Plant Pack One",
			Description: "A pack of 3 that may contain ferns, grass or leafy plants with a chance at a rare version of each",
			Price:       costMap["plantpack1"],
			EventName:   "Buy: plantpack1",
		},
		{
			Name:        "Plant Fertilizer",
			Description: "Double chance at getting a rare on the next plant",
			Price:       costMap["fertilizer"],
			EventName:   "Buy: fertilizer",
		},
	}
	err = mag.AddStorePage(plantItems, hub, "plants")
	if err != nil {
		return nil, err
	}
	// Add fish store page

	// Add info page

	mag.activeIndex = "info"
	// Add accessories store page
	accessoryItems := []StoreItem{
		{
			Name:        "pH Increaser",
			Description: "Raises tank pH levels",
			Price:       0.25,
			EventName:   "ph+",
		},
		{
			Name:        "pH Decreaser",
			Description: "Lowers tank pH levels",
			Price:       0.25,
			EventName:   "ph-",
		},
	}
	mag.AddStorePage(accessoryItems, hub, "items")

	var keys []string
	for key, _ := range mag.pages {
		keys = append(keys, key)
	}
	mag.AddIndexPage(keys, hub, "info")

	setupMagazineEvents(mag, hub)

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
	fishList := []entities.FishList{entities.Guppy, entities.Kirbensis, entities.GoldFish, entities.MollyFish}
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

func LoadFishDescriptions() ([4]string, [4]string) {
	// Keep your existing implementation
	descriptionMap := [4]string{
		"Guppies are Hardy fish that come in\n a variety of vibrant colors.",
		"An exotically patterned fish\n that prefer cave-like structure.",
		"Goldfish are one of the earliest " +
			"fish to be kept as pets.\n They are easy to care for and flexible.",
		"Molly fish are a friendly, active fish.\n " +
			"\n They will swim right up to the glass\n" +
			"when they are ready to eat.",
	}

	names := [4]string{"Guppy", "Kirbensis", "GoldFish", "MollyFish"}
	return descriptionMap, names
}

func CreateJournal(hub *tasks.EventHub) *Magazine {
	mag := CreateSimpleMagazine(hub)

	controls := "Left Click - pick up item \n\nSpace - zoom, un-zoom \n\nLeft Click (zoomed) - select fish\n\n" +
		"Left click(with selected item) - use\n\nHand Icon - white(default) green(held item can be used) \n\nE - return item to shelf\n\n" +
		"Esc- close window"

	dailyCare := "To properly care for your fish you will need to monitor their environment and tailor it to their needs.\n\nWhile zoomed you can check on each fish's condition."

	controlsPage := TextContent{Content: controls, Title: "Controls"}
	fishCare101 := TextContent{Content: dailyCare, Title: "GoldFish Care 101"}
	contents := []TextContent{controlsPage, fishCare101}

	phText := "PH can be monitored by using a test strip on your tank water, you can use the legend on the back of the box to determine the value as closely as you can." +
		"You'll get a chance to guess it exactly each day and receive a bonus if you get it right! PH can be modified long term by the decorations and plants that are added to your tank" +
		"and in the short term by buy a ph booster or a ph nullifier (this will only last that day), different fish have different PH preferences"

	temperature := "It's important that your fish dont get too hot or too cold, until you can afford a tank heater you'll need to make sure any fish you buy can survive at room temperature"

	ph := TextContent{Content: phText, Title: "PH"}
	temp := TextContent{Content: temperature, Title: "Temperature"}

	content2 := []TextContent{ph, temp}

	pages := createTextPages(contents, hub)
	pages2 := createTextPages(content2, hub)
	mag.AddTextPage(pages, "basics")
	mag.AddTextPage(pages2, "Tank Environment")

	mag.activeIndex = "basics"

	return mag
}
