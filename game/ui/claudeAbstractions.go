//go:build old

package ui

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

// ButtonConfig holds common button configuration
type ButtonConfig struct {
	Text           string
	EventText      string // If empty, uses Text
	FontSize       float64
	FontFamily     string
	MinWidth       int
	MinHeight      int
	Padding        widget.Insets
	Colors         ButtonColors
	Tooltip        []string
	DisableOnClick bool
	Image          *ebiten.Image // For sprite buttons
}

type ButtonColors struct {
	Idle     color.Color
	Hover    color.Color
	Pressed  color.Color
	Disabled color.Color
}

// DefaultButtonConfig returns sensible defaults
func DefaultButtonConfig(text string) ButtonConfig {
	return ButtonConfig{
		Text:           text,
		EventText:      "",
		FontSize:       16,
		FontFamily:     "nk57",
		MinWidth:       100,
		MinHeight:      40,
		Padding:        widget.Insets{Left: 10, Right: 10, Top: 5, Bottom: 5},
		Colors:         DefaultButtonColors(),
		DisableOnClick: true,
	}
}

func DefaultButtonColors() ButtonColors {
	return ButtonColors{
		Idle:     color.NRGBA{0, 0, 0, 255},
		Hover:    color.NRGBA{255, 255, 0, 255},
		Pressed:  color.NRGBA{255, 255, 0, 255},
		Disabled: color.NRGBA{128, 128, 128, 255},
	}
}

// SimpleButton creates a basic button with minimal config
func SimpleButton(text string, hub *tasks.EventHub) *widget.Button {
	return ConfigurableButton(DefaultButtonConfig(text), hub)
}

// QuickButton for rapid prototyping - just text and event
func QuickButton(text, eventText string, hub *tasks.EventHub) *widget.Button {
	config := DefaultButtonConfig(text)
	config.EventText = eventText
	return ConfigurableButton(config, hub)
}

// ConfigurableButton creates a button from a config struct
func ConfigurableButton(config ButtonConfig, hub *tasks.EventHub) *widget.Button {
	// Load button image (you might want to make this configurable too)
	buttonImage, err := loadSpriteSelectButtonImage(config.Text)
	if err != nil {
		// Fallback to a default button image or handle error
		buttonImage = loadDefaultButtonImage()
	}

	face, err := util.LoadFont(config.FontSize, config.FontFamily)
	if err != nil {
		// Handle font loading error
		face, _ = util.LoadFont(16, "nk57") // fallback
	}

	eventText := config.EventText
	if eventText == "" {
		eventText = config.Text
	}

	var button *widget.Button
	button = widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text(config.Text, &face, &widget.ButtonTextColor{
			Idle:     config.Colors.Idle,
			Hover:    config.Colors.Hover,
			Pressed:  config.Colors.Pressed,
			Disabled: config.Colors.Disabled,
		}),
		widget.ButtonOpts.TextPadding(&config.Padding),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.MinSize(config.MinWidth, config.MinHeight),
		),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if !button.GetWidget().Disabled {
				if config.DisableOnClick {
					button.GetWidget().Disabled = true
				}
				ev := events.ButtonClickedEvent{ButtonText: eventText}
				hub.Publish(ev)
			}
		}),
	)

	// Add tooltip if provided
	if len(config.Tooltip) > 0 {
		tooltipContainer := MakeToolTipContainer(config.Tooltip)
		button.GetWidget().ToolTip = widget.NewToolTip(
			widget.ToolTipOpts.Content(tooltipContainer),
		)
	}

	return button
}

// MenuBuilder for creating test menus quickly
type MenuBuilder struct {
	container *widget.Container
	hub       *tasks.EventHub
	layout    string // "vertical", "horizontal", "grid"
	spacing   int
}

func NewMenuBuilder(hub *tasks.EventHub) *MenuBuilder {
	return &MenuBuilder{
		hub:     hub,
		layout:  "vertical",
		spacing: 10,
	}
}

func (mb *MenuBuilder) SetLayout(layout string, spacing int) *MenuBuilder {
	mb.layout = layout
	mb.spacing = spacing
	return mb
}

func (mb *MenuBuilder) AddButton(text string) *MenuBuilder {
	return mb.AddButtonWithEvent(text, text)
}

func (mb *MenuBuilder) AddButtonWithEvent(text, eventText string) *MenuBuilder {
	if mb.container == nil {
		mb.initContainer()
	}

	button := QuickButton(text, eventText, mb.hub)
	mb.container.AddChild(button)
	return mb
}

func (mb *MenuBuilder) AddCustomButton(config ButtonConfig) *MenuBuilder {
	if mb.container == nil {
		mb.initContainer()
	}

	button := ConfigurableButton(config, mb.hub)
	mb.container.AddChild(button)
	return mb
}

func (mb *MenuBuilder) AddSpriteButton(text string, sprite *ebiten.Image) *MenuBuilder {
	if mb.container == nil {
		mb.initContainer()
	}

	// You'd need to implement LoadStackSpriteButton or similar
	container, _ := LoadStackSpriteSelectButtonWithToolTip(text, sprite, 16, mb.hub, nil)
	mb.container.AddChild(container)
	return mb
}

func (mb *MenuBuilder) Build() *widget.Container {
	if mb.container == nil {
		mb.initContainer()
	}
	return mb.container
}

func (mb *MenuBuilder) initContainer() {

	switch mb.layout {
	case "horizontal":
		layout := widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionHorizontal),
			widget.RowLayoutOpts.Spacing(mb.spacing),
		)
		mb.container = widget.NewContainer(
			widget.ContainerOpts.Layout(layout),
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				}),
			),
		)
	case "grid":
		layout := widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(3),
			widget.GridLayoutOpts.Spacing(mb.spacing, mb.spacing),
		)
		mb.container = widget.NewContainer(
			widget.ContainerOpts.Layout(layout),
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				}),
			),
		)
	default: // vertical
		layout := widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(mb.spacing),
		)
		mb.container = widget.NewContainer(
			widget.ContainerOpts.Layout(layout),
			widget.ContainerOpts.WidgetOpts(
				widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionCenter,
				}),
			),
		)
	}

}

// Quick test menu creation functions
func QuickTestMenu(hub *tasks.EventHub, buttons ...string) *widget.Container {
	mb := NewMenuBuilder(hub)
	for _, buttonText := range buttons {
		mb.AddButton(buttonText)
	}
	return mb.Build()
}

func QuickHorizontalMenu(hub *tasks.EventHub, buttons ...string) *widget.Container {
	mb := NewMenuBuilder(hub).SetLayout("horizontal", 20)
	for _, buttonText := range buttons {
		mb.AddButton(buttonText)
	}
	return mb.Build()
}

// Button style presets
func SubmitButtonConfig(text string) ButtonConfig {
	config := DefaultButtonConfig(text)
	config.Colors = ButtonColors{
		Idle:     color.NRGBA{255, 255, 0, 255},
		Hover:    color.NRGBA{0, 0, 100, 255},
		Pressed:  color.NRGBA{0, 0, 100, 255},
		Disabled: color.NRGBA{128, 128, 128, 255},
	}
	config.FontFamily = "reglisseClean"
	config.FontSize = 22
	return config
}

func BackButtonConfig() ButtonConfig {
	config := DefaultButtonConfig("Back")
	config.EventText = "back"
	config.MinWidth = 40
	config.MinHeight = 68
	config.Padding = widget.Insets{} // No padding for back button
	return config
}

// Utility function you'll need to implement
func loadDefaultButtonImage() *widget.ButtonImage {
	// Return a simple default button image
	// You might want to create a basic colored rectangle or load a default asset
	return loadSubmitButtonImage() // fallback to existing function
}

// Example usage functions to show how clean this becomes:

func CreateFeatureTestMenu(hub *tasks.EventHub) *widget.Container {
	return QuickTestMenu(hub, "Test Feature A", "Test Feature B", "Settings", "Back")
}

func CreateAdvancedTestMenu(hub *tasks.EventHub) *widget.Container {
	return NewMenuBuilder(hub).
		SetLayout("vertical", 15).
		AddButton("Simple Button").
		AddButtonWithEvent("Custom Event", "special_event").
		AddCustomButton(SubmitButtonConfig("Submit")).
		AddCustomButton(BackButtonConfig()).
		Build()
}

func CreateGridMenu(hub *tasks.EventHub) *widget.Container {
	return NewMenuBuilder(hub).
		SetLayout("grid", 10).
		AddButton("Option 1").
		AddButton("Option 2").
		AddButton("Option 3").
		AddButton("Option 4").
		AddButton("Option 5").
		AddButton("Option 6").
		Build()
}
