package entities

import (
	"encoding/json"
	"github.com/acoco10/fishTankWebGame/game/util"
	"log"
	"os"
)

// EntityTemplate defines the declarative configuration for an entity
type EntityTemplate struct {
	// Flags
	Flags entFlags

	// Z-depth
	Z     int
	BaseZ int

	// Sprite settings
	UnFocusable                    bool
	NoShaderOnFocus                bool
	AbleToBeUnfocusedAutomatically bool
	Draw                           bool
	DontDraw                       bool

	// Timers (duration in seconds)
	Timers map[string]float64

	// Images to load
	ImageAssets map[uint8]string // key -> asset path

	// State machine factory
	StateMachineFactory func() *StateMachine
}

// EntityTemplates - exported so it can be hot-reloaded
var EntityTemplates = buildUiDatTemplates()

func buildUiDatTemplates() map[string]EntityTemplate {
	return map[string]EntityTemplate{
		string(GrandpasJournal): {
			Flags:           likeWindow | altImage | clickTransition | revert,
			NoShaderOnFocus: true,
			Timers: map[string]float64{
				"freeze": 0.5,
			},
			ImageAssets: map[uint8]string{
				Alternate: "uiSprites/magazineAndJournal",
			},
			StateMachineFactory: func() *StateMachine {
				sm := InitStateMachine(nil, ClickedUpdater, nil, nil)
				sm.AppendState(AltImageHovered, PublishPickedUpEvent)
				return sm
			},
		},

		string(FishFood): {
			Flags: updater | outline | autoTransition1,
			Timers: map[string]float64{
				GenFood: 0.1,
			},
			StateMachineFactory: func() *StateMachine {
				sm := InitStateMachine(nil, ClickedUpdater, nil, nil)
				sm.AppendState(ActivationRectUpdaterFishFood, nil)
				return sm
			},
		},

		string(Phreader): {
			Flags:                          swirl | updater | outline | autoTransition1 | altImage,
			AbleToBeUnfocusedAutomatically: true,
			Timers: map[string]float64{
				DoAtTime: 1.75,
			},
			Z:     1,
			BaseZ: 1,
			StateMachineFactory: func() *StateMachine {
				sm := InitStateMachine(nil, ClickedUpdater, nil, nil)
				sm.AppendState(ActivationRectUpdaterPhReader, nil)
				sm.AppendState(UsedInActivationRect, nil)
				return sm
			},
		},

		string(PiggyBank): {
			Flags:                          lowLight,
			AbleToBeUnfocusedAutomatically: true,
			StateMachineFactory: func() *StateMachine {
				return InitPiggyBankStateMachine()
			},
		},

		string(Thermometer): {
			Flags: autoTransition1 | opacityFade | altImage,
			Z:     11,
			BaseZ: 11,
			Timers: map[string]float64{
				Reset: 2.5,
			},
			StateMachineFactory: func() *StateMachine {
				sm := InitStateMachine(nil, ClickedUpdater, AddTempGuage, nil)
				sm.AppendState(FadeOutOnNotHovered, nil)
				sm.AppendState(TurnOffEveryThingOnSpriteAnimationComplete, nil)
				return sm
			},
		},

		string(Magazine): {
			Flags:       revert,
			UnFocusable: true,
			Draw:        false,
			DontDraw:    true,
			Timers: map[string]float64{
				"freeze": 0.5,
			},
			StateMachineFactory: func() *StateMachine {
				return InitStateMachine(DisabledState, AltImageHovered, nil, PublishPickedUpEvent)
			},
		},

		string(LightSwitch): {
			Flags:       noOffset,
			UnFocusable: true,
			Z:           5,
			BaseZ:       5,
		},

		string(WhiteBoard): {
			UnFocusable: true,
			// Note: WhiteBoard needs special initialization with wbSprite
			// Handle this in the switch case
		},

		string(Pillow): {
			Draw:        false,
			UnFocusable: true,
		},
	}
}

// ReloadTemplates - call this to hot-reload template changes
func ReloadTemplates() {
	EntityTemplates = buildUiDatTemplates()
	log.Println("Entity templates reloaded!")
}

// ApplyTemplate applies a template to an entity
func (ent *Entity) ApplyTemplate(template EntityTemplate) error {
	// Apply flags
	ent.flags = template.Flags

	// Apply Z-depth
	if template.Z != 0 {
		ent.Z = template.Z
	}
	if template.BaseZ != 0 {
		ent.UiData.BaseZ = template.BaseZ
	}

	// Apply sprite settings
	if template.UnFocusable {
		ent.Sprite.UnFocusable = true
	}
	if template.NoShaderOnFocus {
		ent.Sprite.NoShaderOnFocus = true
	}
	if template.AbleToBeUnfocusedAutomatically {
		ent.Sprite.AbleToBeUnfocusedAutomatically = true
	}

	if template.DontDraw {
		ent.UiData.DontDraw = true
	}

	// Apply timers
	for key, duration := range template.Timers {
		ent.UiData.Timers[key] = util.NewTimer(duration)
	}

	// Load image assets
	for key, assetPath := range template.ImageAssets {
		img, err := util.LoadImageAssetAsEbitenImage(assetPath)
		if err != nil {
			return err
		}
		ent.Parameters.Images[key] = img
	}

	// Apply state machine
	if template.StateMachineFactory != nil {
		ent.StateMachine = template.StateMachineFactory()
	}

	return nil
}

type EntityConfigJSON struct {
	Flags                          []string           `json:"flags"`
	Z                              int                `json:"z"`
	BaseZ                          int                `json:"baseZ"`
	UnFocusable                    bool               `json:"unFocusable"`
	NoShaderOnFocus                bool               `json:"noShaderOnFocus"`
	AbleToBeUnfocusedAutomatically bool               `json:"ableToBeUnfocusedAutomatically"`
	Draw                           *bool              `json:"draw"` // Pointer to distinguish false from unset
	DontDraw                       bool               `json:"dontDraw"`
	Timers                         map[string]float64 `json:"timers"`
	ImageAssets                    map[string]string  `json:"imageAssets"`
}

// Flag name to bit mapping
var flagNameMap = map[string]entFlags{
	"noOffset":        noOffset,
	"updater":         updater,
	"outline":         outline,
	"autoTransition1": autoTransition1,
	"swirl":           swirl,
	"altImage":        altImage,
	"lowLight":        lowLight,
	"opacityFade":     opacityFade,
	"revert":          revert,
	"likeWindow":      likeWindow,
	"clickTransition": clickTransition,
}

// ParseFlags converts string flags to bitflags
func ParseFlags(flagNames []string) entFlags {
	var flags entFlags
	for _, name := range flagNames {
		if flag, ok := flagNameMap[name]; ok {
			flags |= flag
		} else {
			log.Printf("Warning: unknown flag name: %s", name)
		}
	}
	return flags
}

func loadTemplatesFromJSON(path string) (map[string]EntityTemplate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// JSON uses string keys, we need to map to EntityType
	var jsonConfig map[string]EntityConfigJSON
	if err := json.Unmarshal(data, &jsonConfig); err != nil {
		return nil, err
	}

	templates := make(map[string]EntityTemplate)

	for key, config := range jsonConfig {

		template := EntityTemplate{
			Flags:                          ParseFlags(config.Flags),
			Z:                              config.Z,
			BaseZ:                          config.BaseZ,
			UnFocusable:                    config.UnFocusable,
			NoShaderOnFocus:                config.NoShaderOnFocus,
			AbleToBeUnfocusedAutomatically: config.AbleToBeUnfocusedAutomatically,
			DontDraw:                       config.DontDraw,
			Timers:                         config.Timers,
		}

		// Handle Draw field (default is true, only set if explicitly false in JSON)
		if config.Draw != nil {
			template.Draw = *config.Draw
		} else {
			template.Draw = true
		}

		templates[key] = template
	}

	return templates, nil
}
