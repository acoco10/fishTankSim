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
	"golang.org/x/image/colornames"
	"image/color"
	"log"
	"strings"
)

// import references for loading from consts in entities
const (
	//fish
	goldfish  = string(entities.GoldFish)
	mollyFish = string(entities.MollyFish)
	kirbensis = string(entities.Kirbensis)
	angelFish = string(entities.AngelFish)
	guppyFish = string(entities.Guppy)

	//items
	PHBoost    = string(entities.PhBoost)
	PHReduce   = string(entities.PhReduce)
	Plants1    = "plantPack1"
	Fertilizer = string(entities.Fertilizer)

	//props
	castle    = string(entities.Castle)
	zenFriend = string(entities.ZenFriend)
	zenBridge = string(entities.ZenBridge)
	logProp   = string(entities.Log)
	hotRock   = string(entities.HotRock)
	coolRock  = string(entities.CoolRock)

	//source of truth for ui names of fish visible to user
	goldFishTitle   = "Goldfish"
	guppyFishTitle  = "Guppy"
	angleFishTitle  = "Angel Fish"
	kirbensisTile   = "Kirbensis"
	mollyFishTitle  = "Molly Fish"
	plants1Title    = "Plant Pack One"
	phBoostTitle    = "pH Boost"
	phReduceTitle   = "pH Reduce"
	fertilizerTitle = "Fertilizer"

	zenFriendTitle = "Zen Friend"
	zenBridgeTitle = "Zen Bridge"
	logTitle       = "Log"
	castleTitle    = "Castle"

	buyPrefix        = "Buy:"
	plantPrefix      = "Plant:"
	itemPrefix       = "Item:"
	fishPrefix       = "Fish:"
	decorationPrefix = "Decoration:"
	buyButtonTitle   = "Buy"

	//page names
	Index           = "Index"
	FishStore       = "Fish Store"
	plantStore      = "Plant Store"
	itemStore       = "Item Store"
	decorationStore = "Decoration Store"

	//tags
	noNextPage = "noNextPage"
	noLastPage = "noLastPage"
)

type PurchaseType int8

const (
	Decoration PurchaseType = iota
	Plant
	Item
	Fish
)

const (
	magazineWidth  = 329
	magazineHeight = 446
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
		activeIndex: Index,
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
// keep this wrapper consistent, and we can add to create store page function
func (m *Magazine) AddStorePage(items []StoreItem, hub *tasks.EventHub, pageName string) {
	var tag string
	if pageName == itemStore {
		tag = noNextPage
	}
	page := createStorePage(pageName, items, hub, tag)
	m.pages[pageName] = page
	m.pageOrder = append(m.pageOrder, pageName)
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
		createTextPage(contents[0], leftPage, magazineWidth)
	}
	if len(contents) > 1 {
		createTextPage(contents[1], rightPage, magazineWidth)
	}

	button := LoadSubmitButton(LastPage, hub, "")
	leftPage.AddChild(button)

	button2 := LoadSubmitButton(NextPage, hub, "")
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

func createTextPage(content TextContent, container *widget.Container, width int) *widget.Text {
	textContainer := widget.NewContainer(widget.ContainerOpts.WidgetOpts(
		widget.WidgetOpts.LayoutData(widget.RowLayoutData{Position: widget.RowLayoutPositionCenter})),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(10),
			),
		),
	)

	var contentText *widget.Text
	titleface := registry.FontMap["RockSalt"]
	if content.Title != "" {

		titleText := widget.NewText(
			widget.TextOpts.Text(content.Title, &titleface, color.RGBA{R: 60, G: 160, B: 200, A: 255}),
			widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionStart),
			widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Position: widget.RowLayoutPositionCenter})))

		textContainer.AddChild(titleText)
	}
	textFace := registry.FontMap["RockSalt_18"]
	if content.Content != "" {
		contentText = widget.NewText(
			widget.TextOpts.Text(content.Content, &textFace, colornames.Darkblue),
			widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionStart),
			widget.TextOpts.Padding(&widget.Insets{Top: 10, Left: 50, Right: 25, Bottom: 30}),
			widget.TextOpts.MaxWidth(float64(width)),
		)
		textContainer.AddChild(contentText)
	}

	container.AddChild(textContainer)

	return contentText
}

// createStorePage creates a page with purchasable items
func createStorePage(title string, items []StoreItem, hub *tasks.EventHub, tag string) *widget.Container {
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

	// Add items to right page
	for _, item := range rightItems {
		itemContainer := createStoreItemWidget(item, hub)
		rightPage.AddChild(itemContainer)

	}

	if tag != noLastPage {
		lastButton := LoadSubmitButton(LastPage, hub, "")
		leftPage.AddChild(lastButton)
	}

	if tag != noNextPage {
		nextButton := LoadSubmitButton(NextPage, hub, "")
		rightPage.AddChild(nextButton)
	}

	return root
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
	face := registry.FontMap["RockSalt_12"]
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
		eventName = buyPrefix + item.Name
	}
	buyButton := LoadSubmitButton(buyButtonTitle, hub, eventName)

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
		widget.TextOpts.Text(Index, &face, color.RGBA{R: 60, G: 160, B: 200, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionEnd),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Position: widget.RowLayoutPositionCenter})))

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
			case "Item":
				buyType = Item
			case "Plant":
				buyType = Plant
			case "Decoration":
				buyType = Decoration
			case "Fish":
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
		if ev.ButtonText == NextPage {
			if mag.pageIndex < len(mag.pages)-1 {
				mag.pageIndex++
				mag.activeIndex = mag.pageOrder[mag.pageIndex]
			}
		}
		if ev.ButtonText == LastPage {
			if mag.pageIndex > 0 {
				mag.pageIndex--
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
		kirbensis:  2.0,
		mollyFish:  1.0,
		goldfish:   1.0,
		guppyFish:  1.0,
		angelFish:  5,
		PHBoost:    0.25,
		PHReduce:   0.25,
		Plants1:    1.0,
		Fertilizer: .50,
		zenFriend:  12,
		zenBridge:  5,
		castle:     25,
		logProp:    5,
		hotRock:    .75,
		coolRock:   .75,
	}

	mag := CreateSimpleMagazine(hub)
	mag.CostMap = costMap

	// Add index page

	// Load fish images
	fishImages := LoadFishSprites()

	// Create store items
	fishDescriptions := LoadFishDescriptions()
	fishNames := []string{goldfish, mollyFish, kirbensis, guppyFish, angelFish}
	fishTitle := []string{goldFishTitle, mollyFishTitle, kirbensisTile, guppyFishTitle, angleFishTitle}
	fish := make([]StoreItem, len(fishNames))

	itemImgs := loadStoreImgs()

	for i := 0; i < len(fishNames); i++ {
		fish[i] = StoreItem{
			Name:        fishTitle[i],
			Description: fishDescriptions[fishNames[i]],
			Image:       fishImages[fishNames[i]],
			Price:       costMap[fishNames[i]],
			EventName:   buyPrefix + fishPrefix + fishNames[i],
		}
	}

	decorationDescriptions := LoadDecorationDescriptions()
	decorationNames := []string{castle, zenBridge, zenFriend, logProp}
	decorationTitles := []string{castleTitle, zenBridgeTitle, zenFriendTitle, logTitle}
	decorations := make([]StoreItem, len(decorationNames))

	for i := 0; i < len(decorationNames); i++ {
		decorations[i] = StoreItem{
			Name:        decorationTitles[i],
			Description: decorationDescriptions[decorationNames[i]],
			Image:       itemImgs[decorationNames[i]],
			Price:       costMap[decorationNames[i]],
			EventName:   buyPrefix + fishPrefix + decorationNames[i],
		}
	}

	mag.AddStorePage(fish, hub, FishStore)

	plantItems := []StoreItem{
		{
			Name: plants1Title,
			Description: "A pack of 3 that may contain ferns, grass or leafy plants with a chance at a rare version of each.\n\n" +
				"Plants boost your ecosystem",
			Price:     costMap[Plants1],
			Image:     itemImgs[Plants1],
			EventName: buyPrefix + plantPrefix + Plants1,
		},
		{
			Name:        fertilizerTitle,
			Description: "Double chance at a rare on your next plant",
			Price:       costMap[Fertilizer],
			Image:       itemImgs[Fertilizer],
			EventName:   buyPrefix + itemPrefix + Fertilizer,
		},
	}

	mag.AddStorePage(plantItems, hub, plantStore)

	mag.AddStorePage(decorations, hub, decorationStore)

	mag.activeIndex = Index
	// Add accessories store page
	accessoryItems := []StoreItem{
		{
			Name:        phBoostTitle,
			Description: "Raises tank pH level.",
			Image:       itemImgs["phaidb"],
			Price:       costMap[PHBoost],
			EventName:   buyPrefix + itemPrefix + PHBoost,
		},
		{
			Name:        phReduceTitle,
			Description: "Lowers tank pH level.",
			Image:       itemImgs["phaidr"],
			Price:       costMap[PHReduce],
			EventName:   buyPrefix + itemPrefix + PHReduce,
		},
		{
			Name:        "Cool Rock",
			Description: "Cools down your tank.",
			Image:       itemImgs["coolRock"],
			Price:       0.75,
			EventName:   buyPrefix + decorationPrefix + coolRock,
		},
		{
			Name:        "Hot Rock",
			Description: "Heats up your tank.",
			Image:       itemImgs["hotRock"],
			Price:       0.75,
			EventName:   buyPrefix + decorationPrefix + hotRock,
		},
	}

	mag.AddStorePage(accessoryItems, hub, itemStore)

	pageOrder := []string{FishStore, decorationStore, plantStore, itemStore}

	mag.AddIndexPage(pageOrder, hub, Index)
	mag.activeIndex = Index

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

func LoadFishSprites() map[string]*ebiten.Image {
	fishList := []entities.FishList{entities.GoldFish, entities.MollyFish, entities.Kirbensis, entities.Guppy, entities.AngelFish}
	imgs := make(map[string]*ebiten.Image)
	for _, fish := range fishList {
		fishSprite, err := entities.LoadFishAnimations(fish, 2)
		if err != nil {
			log.Fatal(err)
		}
		imgs[string(fish)] = fishSprite["swimming"].GetFirstFrameAsStaticImage()
	}
	return imgs
}

func LoadDecorationDescriptions() map[string]string {
	return map[string]string{
		zenBridge: "Add some calm and tranquility to your tank.\n\n -Increases PH by 1.0, \n\n-Decreases fish stress by 0.25\n\n+zen set",
		zenFriend: "Add some calm and tranquility to your tank.\n\n-Decreases PH by 2.0, \n\n-Decreases fish stress by 1n\n\n+zen set",
		logProp:   "A cheap nature booster.\n\n -Increase PH by 2.0, \n\n+Eco System",
		castle:    "A decoration fit for a king.\n\n -Decreases PH by 2.0, \n\n+Cave structure",
	}
}

func LoadFishDescriptions() map[string]string {
	// Keep your existing implementation
	descriptionMap := map[string]string{

		goldfish: "Goldfish are one of the earliest " +
			"fish to be kept as pets. They are easy to care for and flexible.",
		mollyFish: "Molly fish are a friendly, active fish. " +
			"They will swim right up to the glass" +
			"when they are ready to eat.",
		kirbensis: "An exotically patterned fish that prefer cave-like structure.",
		guppyFish: "Guppies are Hardy fish that come in a variety of vibrant colors.",
		angelFish: "Angel fish need warm temperatures and acidic temperatures. They are famous for their unique shape and beautiful stripe patterns",
	}

	return descriptionMap
}

func loadStoreImgs() map[string]*ebiten.Image {
	return util.LoadDirectoryImages("images/storeAssets")
}
