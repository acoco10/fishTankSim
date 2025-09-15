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
	"golang.org/x/image/colornames"
	"image/color"
	"log"
	"strings"
)

type Journal struct {
	pages       map[string]*widget.Container
	pageIndex   int
	activeIndex string
	pageOrder   []string
}

// StoreItem represents an item that can be purchased

// TextContent represents content for a text page

func (m *Journal) ActiveWindow() *widget.Container {
	return m.pages[m.activeIndex]
}

func CreateNewJournal(hub *tasks.EventHub) *Journal {
	journal := &Journal{
		pages:       make(map[string]*widget.Container),
		activeIndex: "info",
	}

	// Set up event subscriptions
	return journal
}

// AddTextPage adds a simple text page to the magazine
func (m *Journal) AddTextPage(page *widget.Container, pageName string) error {
	m.pages[pageName] = page
	m.pageOrder = append(m.pageOrder, pageName)
	return nil
}

// AddIndexPage adds a navigation page with buttons to other pages
func createJournalPageBase() (*widget.Container, *widget.Container, *widget.Container) {
	magNineSlice, flippedMagNineSlice, err := LoadJournalNineSlice()
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
func createTextJournalPages(contents []TextContent, hub *tasks.EventHub) *widget.Container {

	root, leftPage, rightPage := createJournalPageBase()

	speciesButton := LoadOutlineTextColoreBg(colornames.Darkgoldenrod, "Species Guide", hub, "")
	basicsButton := LoadOutlineTextColoreBg(colornames.Darkcyan, "basics", hub, "")
	indexContainer := widget.NewContainer(widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 10)),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
				widget.RowLayoutOpts.Spacing(20),
				widget.RowLayoutOpts.Padding(&widget.Insets{Top: -42, Left: 10, Right: 10}),
			),
		))

	indexContainer.AddChild(basicsButton)
	indexContainer.AddChild(speciesButton)
	rightPage.AddChild(indexContainer)

	if contents[0].Image != nil {
		createImagePage(leftPage, contents)
	} else {
		createTextPage(contents[0], leftPage)
	}
	if len(contents) > 1 {
		createTextPage(contents[1], rightPage)
	}

	button2 := LoadSubmitButton("Back", hub, "")
	leftPage.AddChild(button2)

	button3 := LoadSubmitButton("Forward", hub, "")
	rightPage.AddChild(button3)

	return root

}

func createJournalImagePage(container *widget.Container, contents []TextContent) {
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

func createJournalTextPage(content TextContent, container *widget.Container) {
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

// createIndexPage creates a navigation page with buttons
func createJournalIndexPage(buttonTexts []string, hub *tasks.EventHub) *widget.Container {
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

func defaultJournalSubs(mag *Magazine, hub *tasks.EventHub) {
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

func LoadJournalNineSlice() (*eimage.NineSlice, *eimage.NineSlice, error) {
	bground, err := util.LoadImageAssetAsEbitenImage("uiSprites/journalLeft")
	if err != nil {
		return nil, nil, err
	}

	flipImg, err := util.LoadImageAssetAsEbitenImage("uiSprites/journalRight")
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

func CreateJournal(hub *tasks.EventHub) *Journal {
	mag := CreateNewJournal(hub)

	controls := "Left Click - pick up item \n\nSpace - zoom, un-zoom \n\nLeft Click (zoomed) - select fish\n\n" +
		"Left click(with selected item) - use\n\nHand Icon - white(default) green(held item can be used) \n\nE - return item to shelf\n\n" +
		"Esc- close window"

	dailyCare := "To properly care for your fish you will need to monitor their environment and tailor it to their needs.\n\nWhile zoomed you can check on each fish's condition."

	controlsPage := TextContent{Content: controls, Title: "Controls"}
	fishCare101 := TextContent{Content: dailyCare, Title: "GoldFish Care 101"}
	contents := []TextContent{controlsPage, fishCare101}

	phText := " - Monitored by using a test strip on tank water\n\n" +
		" - use the legend on the back of the box to determine the value as closely as you can\n\n" +
		" - different fish have different PH preferences\n\n" +
		" - ph balancers can decrease or increase your tank PH for a day the day after they are used\n\n" +
		" - long term, decorations permanently increase and decrease the PH of your tank but not instantly"

	temperature := " - Fish have different ideal temperatures ranges based on species and individual differences\n\n" +
		"- Use hot and cold rocks to increase or decrease your tanks temperature temporarily or a water heater to maintain a constant temperature"

	ph := TextContent{Content: phText, Title: "PH"}
	temp := TextContent{Content: temperature, Title: "Temperature"}
	content2 := []TextContent{ph, temp}

	mfs, err := entities.GenMollyFishStats()

	if err != nil {
		log.Fatal("error loading fish stats for grandpa's journal:", err)
	}

	idealPH := fmt.Sprintf("Desired PH: %0.2f\n\n", mfs.IdealPH)
	idealTemp := fmt.Sprintf("Favorite Temp: %d", mfs.IdealTemperature)
	mollyFishSpeciesGuide := TextContent{Title: "Molly Fish Notes", Content: idealPH + idealTemp}

	content3 := []TextContent{mollyFishSpeciesGuide}

	icons, err := util.LoadImageAssetAsEbitenImage("uiSprites/fishFactorIcons")
	if err != nil {
		log.Fatal(err)
	}
	iconLabels := []string{"thumbsUp", "thumbsNeutral", "thumbsDown", "otherFish", "structures", "temperature", "ph"}
	imageMap, indMap := util.ChopUpIcons(icons, iconLabels, 32)

	phContent := TextContent{Title: "Icons", Content: "PH", Image: imageMap["ph"]}
	tempContent := TextContent{Content: "Temperature", Image: imageMap["temperature"]}
	thumbsUp := TextContent{Content: "Ok", Image: indMap["thumbsUp"]}
	thumbsNeutral := TextContent{Content: "Ok", Image: indMap["thumbsNeutral"]}
	thumbsDown := TextContent{Content: "Ok", Image: indMap["thumbsDown"]}

	pages := createTextJournalPages(contents, hub)
	pages2 := createTextJournalPages(content2, hub)
	pages3 := createTextJournalPages(content3, hub)
	pages4 := createTextPages([]TextContent{phContent, tempContent, thumbsUp, thumbsNeutral, thumbsDown}, hub)

	mag.AddTextPage(pages, "basics")
	mag.AddTextPage(pages2, "Tank Environment")
	mag.AddTextPage(pages3, "Species Guide")
	mag.AddTextPage(pages4, "icons")

	mag.activeIndex = "basics"

	return mag
}
