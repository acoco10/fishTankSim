package ui

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	eimage "github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"log"
	"strings"
)

const (
	Basics                           = "Basics"
	SpeciesGuide                     = "Species Guide"
	EcoSystemGuide                   = "Eco System Guide"
	UnlockTextPreFix                 = "FishUnlocked:"
	JournalPageImageWidth            = 310
	JournalPAgeImageHeight           = 465
	TextWidthModifierThatLooksDecent = 200
	NextPage                         = "Next Page >"
	LastPage                         = "< Last Page"
	BaseImage                        = "BaseImage"
	HighlightImage                   = "HighlightImage"
	ImageRef                         = "ImageReference"

	rowLayOutSpacing = 50
)

type JournalPages struct {
	root         *widget.Container
	leftPage     *widget.Container
	rightPage    *widget.Container
	indexButtons map[string]*widget.Button
}
type SpeciesPages struct {
	JournalPages
	speciesButtons [2]*widget.Container
	order          [2]string
}

type Journal struct {
	pages                 map[string]*JournalPages
	speciesPages          []*SpeciesPages
	pageIndex             int
	activeIndex           string
	species               map[string]TextContent
	pageOrder             []string
	unlockables           map[string]bool
	unlockButtons         map[string]*widget.Button
	unLockOrder           []string
	indexButtons          []*widget.Button
	buttonUpdaters        []uint32
	submitButtonImage     *widget.ButtonImage
	indexButtonImage      *widget.ButtonImage
	submitButtonHighlight *widget.ButtonImage
	indexButtonHighlight  *widget.ButtonImage
}

// StoreItem represents an item that can be purchased

// TextContent represents content for a text page

func (m *Journal) ActiveWindow() *widget.Container {
	return m.pages[m.activeIndex].root
}

func CreateNewJournal(hub *tasks.EventHub) *Journal {
	journal := &Journal{
		pages:       make(map[string]*JournalPages),
		activeIndex: "info",
	}

	return journal
}

// AddTextPage adds a simple text page to the magazine
func (m *Journal) AddTextPage(page *JournalPages, pageName string) {
	m.pages[pageName] = page
	m.pageOrder = append(m.pageOrder, pageName)

}

// AddIndexPage adds a navigation page with buttons to other pages
func createJournalPageBase(rightTitle string, leftTitle string, hub *tasks.EventHub) (*widget.Container, *widget.Container, *widget.Container, map[string]*widget.Button) {
	magNineSlice, flippedMagNineSlice := LoadJournalNineSlice()

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

	rightPageIndex, butMap := MakeJournalIndex(hub)

	customRightPage := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(JournalPageImageWidth*2, JournalPAgeImageHeight*2)),
		widget.ContainerOpts.BackgroundImage(flippedMagNineSlice),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(0),
				widget.RowLayoutOpts.Padding(&widget.Insets{Top: 30, Left: 0, Right: 0, Bottom: 40}),
			),
		),
	)

	customLeftPage := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(magNineSlice),
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(JournalPageImageWidth*2, JournalPAgeImageHeight*2)),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(0),
				widget.RowLayoutOpts.Padding(&widget.Insets{Top: 40, Left: 0, Right: 0, Bottom: 40}),
			),
		),
	)

	rtc := MakeContentContainer()
	ltc := MakeContentContainer()

	rightTitleText := makeTitle(rightTitle)
	leftTitleText := makeTitle(leftTitle)

	customRightPage.AddChild(rightPageIndex)
	rightPage := makePageContainer(flippedMagNineSlice)
	rightPage.AddChild(rightTitleText)
	customRightPage.AddChild(rightPage)
	rightPage.AddChild(rtc)

	leftPage := makePageContainer(magNineSlice)
	leftPage.AddChild(leftTitleText)
	customLeftPage.AddChild(leftPage)
	leftPage.AddChild(ltc)

	rootContainer.AddChild(customLeftPage, customRightPage)

	lastpg := LoadSubmitButton(LastPage, hub, "")
	leftPage.AddChild(lastpg)

	button2 := LoadSubmitButton(NextPage, hub, "")
	rightPage.AddChild(button2)
	return rootContainer, ltc, rtc, butMap
}

func makePageContainer(image *eimage.NineSlice) *widget.Container {
	Page := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(JournalPageImageWidth, JournalPAgeImageHeight)),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(0),
				widget.RowLayoutOpts.Padding(&widget.Insets{Top: 20, Left: 10, Right: 10}),
			),
		),
	)
	return Page
}

func makeTitle(title string) *widget.Text {
	face := registry.FontMap["RockSalt"]

	Title := widget.NewText(
		widget.TextOpts.Text(title, &face, color.RGBA{R: 60, G: 160, B: 200, A: 255}),
		widget.TextOpts.Position(widget.TextPositionCenter, widget.TextPositionStart),
		widget.TextOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{Position: widget.RowLayoutPositionCenter})))
	return Title
}

// createTextPage creates a page with title, text content, and optional image
func createTextJournalPages(leftTitle string, rightTitle string, contents []TextContent, hub *tasks.EventHub) (*widget.Container, JournalPages) {

	root, leftPage, rightPage, butMap := createJournalPageBase(rightTitle, leftTitle, hub)
	pagesStruct := JournalPages{root: root, rightPage: rightPage, leftPage: leftPage, indexButtons: butMap}

	if len(contents) != 0 {
		if contents[0].Image != nil {
			createImagePage(leftPage, contents)
		} else {
			createTextPage(contents[0], leftPage, JournalPageImageWidth*2)
		}
		if len(contents) > 1 {
			createTextPage(contents[1], rightPage, JournalPageImageWidth*2)
		}
	}

	return root, pagesStruct

}

func MakeJournalIndex(hub *tasks.EventHub) (*widget.Container, map[string]*widget.Button) {

	speciesButtonCont, speciesButton := LoadOutlineTextColoreBg(colornames.Darkgoldenrod, "Species Guide", hub, string(SpeciesGuide))

	basicsButtonCont, BasicButton := LoadOutlineTextColoreBg(colornames.Darkcyan, "Basics", hub, string(Basics))
	ecoButtonCont, ecoButton := LoadOutlineTextColoreBg(colornames.Cornflowerblue, "Eco System", hub, string(EcoSystemGuide))

	indexContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(100, 10),
			widget.WidgetOpts.LayoutData(
				widget.RowLayoutData{Position: widget.RowLayoutPositionCenter, MaxHeight: 10})),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
				widget.RowLayoutOpts.Spacing(20),
				widget.RowLayoutOpts.Padding(&widget.Insets{Top: -26, Left: 10, Right: 10}),
			)))

	indexContainer.AddChild(basicsButtonCont, speciesButtonCont, ecoButtonCont)

	indexButMap := map[string]*widget.Button{
		Basics:         BasicButton,
		SpeciesGuide:   speciesButton,
		EcoSystemGuide: ecoButton,
	}

	return indexContainer, indexButMap
}

func MakeContentContainer() *widget.Container {
	textContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(widget.WidgetOpts.MinSize(JournalPageImageWidth, JournalPAgeImageHeight+TextWidthModifierThatLooksDecent)),
		//widget.ContainerOpts.BackgroundImage(eimage.NewNineSliceColor(colornames.Darkblue)),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{Position: widget.RowLayoutPositionCenter, MaxHeight: JournalPAgeImageHeight + TextWidthModifierThatLooksDecent}),
		))
	return textContainer
}

func InitHiglightOnButton(button *widget.Button, btnImg *widget.ButtonImage, highlightImage *widget.ButtonImage) uint32 {
	ent := &entities.Entity{}
	ent.UpdateFunc = HighLightButtonUpdater
	ent.Parameters = make(map[string]any)
	ent.Parameters["button"] = button
	ent.Parameters["buttonImage"] = btnImg
	ent.Parameters["highlightImage"] = highlightImage
	ent.Parameters[entities.Counter] = 0
	return entities.RegisterEntity(ent)
}

func defaultJournalSubs(journal *Journal, hub *tasks.EventHub) {
	hub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		if ev.ButtonText == LastPage {
			if journal.pageIndex > 0 {
				journal.pageIndex--
				journal.activeIndex = journal.pageOrder[journal.pageIndex]
			}
		}
		if ev.ButtonText == NextPage {
			if journal.pageIndex < len(journal.pageOrder)-1 {
				journal.pageIndex++
				journal.activeIndex = journal.pageOrder[journal.pageIndex]
			}
		}
	})

}

func LoadJournalNineSlice() (*eimage.NineSlice, *eimage.NineSlice) {
	bground, err := util.LoadImageAssetAsEbitenImage("uiSprites/journalLeft")
	if err != nil {
		log.Fatal("Error loading journal right page:", err)
	}

	flipImg, err := util.LoadImageAssetAsEbitenImage("uiSprites/journalRight")
	if err != nil {
		log.Fatal("Error loading journal right page:", err)
	}

	magNineSlice := eimage.NewNineSlice(
		bground, [3]int{32, bground.Bounds().Dx() - 64, 32},
		[3]int{32, bground.Bounds().Dy() - 64, 32})

	flipNineSlice := eimage.NewNineSlice(
		flipImg, [3]int{32, bground.Bounds().Dx() - 64, 32},
		[3]int{32, bground.Bounds().Dy() - 64, 32})

	return magNineSlice, flipNineSlice
}

func CreateJournal(hub *tasks.EventHub) *Journal {
	journal := CreateNewJournal(hub)

	controls :=
		"Left Click - Pick up item \n\nSpace - Zoom, Unzoom \n\nLeft Click (Zoomed) - Select Fish\n\n" +
			"Left Click(with selected item) - Use\n\nHand Icon - White(default) Green(Held item can be used) \n\nESC or E - Return item to shelf\n\n" +
			"ESC- close window"

	dailyCare :=
		"To properly care for your fish you will need to monitor their environment and tailor it to their individual and species specific needs." +
			"\n\nZoom in to get a closer look at your fish and click on them to see their mood."

	controlsPage := TextContent{Content: controls}
	fishCare101 := TextContent{Content: dailyCare}
	contents := []TextContent{controlsPage, fishCare101}

	phText :=
		"pH is a measure of a substances hydrogen ion content\n\n" +
			" - Different species of fish have different pH preferences\n\n" +
			" - Monitored by using a litmus strip on tank water\n\n" +
			"-  The strip will change color after being submerged for a few seconds, acidic substances cause it to the litmus paper to " +
			"turn red and basic substances cause it to turn green\n\n" +
			" - Compare with the color with legend image to determine the value as closely as you can\n\n" +
			" - pH balancers will increase or decrease your tanks pH immediately but only last a few day \n\n" +
			" - long term, decorations made of different substances will permanently increase and decrease the pH of " +
			"your tank but the tank takes a few days to reach equilibrium\n\n" +
			"Acidic = Lower pH\n\n" +
			"Basic = less acidic, higher pH\n\n"

	temperature :=
		" - Fish have different ideal temperatures ranges based on species  type and individual preferences\n\n" +
			" - Until you have a Tank Heater, temperatures will fluctuate" +
			"- Hot and Cold Rocks are a cheap temporary solution but not precise and your tank will still be influenced by temperature fluctuations in your room."

	ph := TextContent{Content: phText}
	temp := TextContent{Content: temperature}
	content2 := []TextContent{ph, temp}

	mfishNotes := "October 14th\n" +
		"- Its a balmy 75 degree fall day, my new Molly Fish Cassandra seems quite happy.\n\n" +
		"December 21st\n\n" +
		"- The water in my new apartment is so basic, I had to use lots of treatments to raise it to a level she likes."
	mollyFishSpeciesGuide := TextContent{Title: "Molly Fish Notes", Content: mfishNotes}

	goldFishNotes := "September 01\n\n" +
		"- I got my first fish this week, Thomas! A delightful Goldfish." +
		"He is taking some time to adjust to his new home but seems health so far\n\n" +
		"- I've heard Gold fish are quite hardy and can thrive in a variety of environments\n\n" +
		"September 23\n\n" +
		"-Thomas has settled in nicely to his new home, he seems fine with our normal water even though its rather acidic\n\n" +
		"January 21st\n\n" +
		"-Our water heater went out and my room grew quite cold\n\n" +
		"-Poor Thomas got sick and died soon after"

	goldFishGuide := TextContent{Title: "Goldfish Notes", Content: goldFishNotes}

	_, basicsPage := createTextJournalPages("Controls", "Fish Care 101", contents, hub)
	_, ecoPage := createTextJournalPages(EcoSystemGuide+":PH", EcoSystemGuide+":Temperature", content2, hub)

	journal.AddTextPage(&basicsPage, Basics)
	journal.AddTextPage(&ecoPage, EcoSystemGuide)

	defaultJournalSubs(journal, hub)
	setupJournalEventSubs(journal, hub)
	journal.activeIndex = Basics

	mButtonCont, mButton := LoadOutlineTextButtonSubmitBgGetButton("Molly Fish Notes", hub, UnlockTextPreFix+string(entities.MollyFish))

	gfButtonCont, gfButton := LoadOutlineTextButtonSubmitBgGetButton("Gold Fish Notes", hub, UnlockTextPreFix+string(entities.GoldFish))

	journal.species = make(map[string]TextContent)
	journal.species[mollyFish] = mollyFishSpeciesGuide
	journal.species[goldfish] = goldFishGuide

	journal.unlockables = make(map[string]bool)
	journal.unlockables[mollyFish] = false
	journal.unlockables[goldfish] = false

	journal.unlockButtons = make(map[string]*widget.Button)
	journal.unlockButtons[mollyFish] = mButton
	journal.unlockButtons[goldfish] = gfButton

	speciesOnPage1 := [2]string{goldfish, mollyFish}
	speciesButtonPageOne := [2]*widget.Container{gfButtonCont, mButtonCont}

	unlockPage := createSpeciesPages("", speciesOnPage1, speciesButtonPageOne, hub)

	journal.pages[SpeciesGuide] = &unlockPage.JournalPages
	journal.speciesPages = append(journal.speciesPages, &unlockPage)
	journal.unLockOrder = append(journal.unLockOrder, speciesOnPage1[0], speciesOnPage1[1])
	journal.pageOrder = append(journal.pageOrder, SpeciesGuide)

	submitButtonImg := loadSubmitButtonImage()
	submitButtonHighlight := loadHighlightSubmitButtonImage()
	indexButtonImage := LoadColoredImageButton(colornames.Goldenrod)
	indexButtonHighlight := LoadColoredImageButton(colornames.Lightgoldenrodyellow)

	journal.submitButtonImage = submitButtonImg
	journal.submitButtonHighlight = submitButtonHighlight

	journal.indexButtonImage = indexButtonImage
	journal.indexButtonHighlight = indexButtonHighlight

	return journal
}

func createSpeciesPages(title string, species [2]string, buttons [2]*widget.Container, hub *tasks.EventHub) SpeciesPages {
	root, leftPage, rightPage, butMap := createJournalPageBase(title, title, hub)

	jp := JournalPages{root: root, leftPage: leftPage, rightPage: rightPage, indexButtons: butMap}
	sp := SpeciesPages{JournalPages: jp, speciesButtons: buttons}
	leftLayoutContainer := makeButContainer()

	rightLayoutContainer := makeButContainer()

	for i, button := range buttons {
		sp.order[i] = species[i]
		if i < 1 {
			leftLayoutContainer.AddChild(button)
		} else {
			rightLayoutContainer.AddChild(button)
		}
		for _, child := range button.Children() {
			child.GetWidget().Disabled = true //setting for unlockable buttons, not generic at all!!
		}
	}

	leftPage.AddChild(leftLayoutContainer)
	rightPage.AddChild(rightLayoutContainer)

	return sp
}

func makeButContainer() *widget.Container {
	leftLayoutContainer := widget.NewContainer(widget.ContainerOpts.WidgetOpts(
		widget.WidgetOpts.LayoutData(widget.RowLayoutData{Position: widget.RowLayoutPositionCenter}), widget.WidgetOpts.MinSize(JournalPageImageWidth*2+100, JournalPAgeImageHeight*2)),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Padding(&widget.Insets{Left: 50, Right: 50, Top: 50}),
				widget.RowLayoutOpts.Direction(widget.DirectionVertical),
				widget.RowLayoutOpts.Spacing(rowLayOutSpacing),
			),
		),
	)
	return leftLayoutContainer
}

func setupJournalEventSubs(journal *Journal, hub *tasks.EventHub) {
	hub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		_, pageExists := journal.pages[ev.ButtonText]
		if pageExists {
			journal.activeIndex = ev.ButtonText
			for i, pgName := range journal.pageOrder {
				if pgName == ev.ButtonText {
					journal.pageIndex = i
				}
			}

		} else {
			println("mag page:", ev.ButtonText, "does not exist")
		}
	})

	hub.Subscribe(events.ButtonClickedEvent{}, func(e tasks.Event) {
		ev := e.(events.ButtonClickedEvent)
		if strings.HasPrefix(ev.ButtonText, UnlockTextPreFix) {
			species := strings.TrimPrefix(ev.ButtonText, UnlockTextPreFix)
			var pageWithSpecies *SpeciesPages
			specFound := false
			var q int
			var p int
			for i, spec := range journal.unLockOrder {
				if spec == species {
					pageWithSpecies = journal.speciesPages[p]
					specFound = true
					q = i
					break
				}
				if p > 0 && i%2 == 0 {
					p++
				}
			}

			if !specFound {
				return
			}
			var pageToMutate *widget.Container
			if q%2 == 0 {
				pageToMutate = pageWithSpecies.leftPage
			} else {
				pageToMutate = pageWithSpecies.rightPage
			}

			pageToMutate.RemoveChildren()
			text := createTextPage(journal.species[species], pageToMutate, JournalPageImageWidth*2)
			clr := color.RGBA{5, 5, 5, 20}
			text.SetColor(clr)
			ent := &entities.Entity{}
			ent.Parameters = make(map[string]any)
			ent.Parameters["text"] = text
			ent.Parameters["color"] = clr
			ent.UpdateFunc = TextOpacityUpdater
			entities.RegisterEntity(ent)

			rect := pageToMutate.GetWidget().Rect
			makeUiParticles(rect)
			for _, id := range journal.buttonUpdaters {
				entities.RemoveEntity(id)
			}
		}
	})

	hub.Subscribe(events.NewSpecies{}, func(e tasks.Event) {
		ev := e.(events.NewSpecies)
		_, exists := journal.unlockButtons[ev.Species]

		if !exists {
			return
		}

		journal.unlockButtons[ev.Species].GetWidget().Disabled = false
		id := InitHiglightOnButton(journal.unlockButtons[ev.Species], journal.submitButtonImage, journal.submitButtonHighlight)
		journal.buttonUpdaters = append(journal.buttonUpdaters, id)
		journal.unlockables[ev.Species] = true

		for _, page := range journal.pages {
			//need to highlight species button on every page, state is not centralized
			id2 := InitHiglightOnButton(page.indexButtons[SpeciesGuide], journal.indexButtonImage, journal.indexButtonHighlight)
			journal.buttonUpdaters = append(journal.buttonUpdaters, id2)
		}
	})
}

func makeUiParticles(rect image.Rectangle) {
	pEnt := entities.NewGenericParticleSystem(
		float64(rect.Min.X+220),
		float64(rect.Min.Y+500),
		image.Rect(0, 0, ScreenWidth*4, ScreenHeight*4), 8)

	pEnt.PConfig = &entities.ParticleConfig{
		XVariance:         70,
		YVariance:         20,
		XVelocityVariance: 15,
		YVelocityVariance: 200,
		MaxLife:           0,
		BaseYVelocity:     -1000,
		Scale:             4.0,
		AlphaDecay:        0.5}
	pEnt.SpawnRate = 250
	pEnt.MaxParticles = 1000
	pEnt.EndAfter = 1
	ent := &entities.Entity{ParticleSystem: pEnt, Z: 1, Sprite: pEnt.Sprite}
	ent.SetOverUIEffect()
	entities.RegisterEntity(ent)

	pEnt2 := entities.NewGenericParticleSystem(
		float64(rect.Max.X-220),
		float64(rect.Min.Y+500),
		image.Rect(0, 0, ScreenWidth*4, ScreenHeight*4), 8)
	pEnt2.PConfig = &entities.ParticleConfig{
		XVariance:         70,
		YVariance:         20,
		XVelocityVariance: 15,
		YVelocityVariance: 200,
		MaxLife:           0,
		BaseYVelocity:     -1000,
		Scale:             4.0}
	pEnt2.SpawnRate = 250
	pEnt2.MaxParticles = 1000
	pEnt2.EndAfter = 1
	ent2 := &entities.Entity{ParticleSystem: pEnt2, Z: 1, Sprite: pEnt2.Sprite}
	ent2.SetOverUIEffect()
	entities.RegisterEntity(ent2)

}

func makeIconMap() []TextContent {
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

	return []TextContent{phContent, tempContent, thumbsUp, thumbsNeutral, thumbsDown}
}

func HighLightButtonUpdater(ent *entities.Entity) {
	counter := entities.GetParam[int](ent, entities.Counter)
	counter++
	if counter == 30 {
		button := entities.GetParam[*widget.Button](ent, "button")
		buttonimg := entities.GetParam[*widget.ButtonImage](ent, "buttonImage")
		button.SetImage(buttonimg)
	}

	if counter == 60 {
		button := entities.GetParam[*widget.Button](ent, "button")
		highlightimg := entities.GetParam[*widget.ButtonImage](ent, "highlightImage")
		button.SetImage(highlightimg)
		counter = 0
	}

	ent.Parameters[entities.Counter] = counter
}

func TextOpacityUpdater(ent *entities.Entity) {
	clr := entities.GetParam[color.RGBA](ent, "color")
	txt := entities.GetParam[*widget.Text](ent, "text")
	if clr.A < 255 {
		clr.A += 1
	} else {
		entities.RemoveEntity(ent.Id)
	}
	txt.SetColor(clr)

	ent.Parameters["color"] = clr
}
