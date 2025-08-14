package system

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"math/rand"
)

type TankModifier struct {
	Name          string
	EffectedParam string
	Amount        float64
	Duration      int
	DaysActive    int
}
type Environment struct {
	Temperature        int
	NaturalPHLevel     float64
	ModifiedPHLevel    float64
	NatureBoosters     int
	modifiers          []TankModifier
	temporaryModifiers []TankModifier
	usedTempModifier   []TankModifier
}

func InitEnvironment(e *Environment) {
	e.Temperature = rand.Intn(10) + 65
	e.NaturalPHLevel = rand.Float64()*3 + 6
	e.ModifiedPHLevel = e.NaturalPHLevel
}

func (env *Environment) Subscribe(hub *tasks.EventHub) {
	hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		env.Temperature = rand.Intn(10) + 67

		for i, mod := range env.usedTempModifier {
			switch mod.EffectedParam {
			case "PH":
				env.ModifiedPHLevel -= mod.Amount
				if len(env.usedTempModifier) > 1 {
					env.usedTempModifier = append(env.usedTempModifier[0:i], env.usedTempModifier[i+1:]...)
				} else {
					env.usedTempModifier = []TankModifier{}
				}
			}
		}

		for _, mod := range env.modifiers {
			switch mod.EffectedParam {
			case "PH":
				if env.ModifiedPHLevel != env.ModifiedPHLevel+mod.Amount {
					mod.DaysActive++
					change := util.Lerp64(env.ModifiedPHLevel, env.ModifiedPHLevel+mod.Amount, float64(mod.DaysActive)*0.25)
					env.ModifiedPHLevel = change
				}
			}
		}

		for i, mod := range env.temporaryModifiers {
			switch mod.EffectedParam {
			case "PH":
				mod.DaysActive++
				env.ModifiedPHLevel += mod.Amount
				env.usedTempModifier = append(env.usedTempModifier, mod)
				if mod.DaysActive >= mod.Duration {
					if len(env.temporaryModifiers) > 1 {
						env.temporaryModifiers = append(env.temporaryModifiers[0:i], env.temporaryModifiers[i+1:]...)
					} else {
						env.temporaryModifiers = []TankModifier{}
					}
				}
			}
		}

	})

	hub.Subscribe(events.UISpriteAction{}, func(e tasks.Event) {
		ev := e.(events.UISpriteAction)
		if ev.UiSprite == "PHModifier" {
			if ev.UiSpriteAction == "ph+" {
				env.AppendTempModifier(TankModifier{Name: ev.UiSpriteAction, EffectedParam: "PH", Amount: .2, Duration: 1})
			}
			if ev.UiSpriteAction == "ph-" {
				env.AppendTempModifier(TankModifier{Name: ev.UiSpriteAction, EffectedParam: "PH", Amount: -.2, Duration: 1})
			}
		}
	})

}

func (env *Environment) AppendTempModifier(modifier TankModifier) {
	env.temporaryModifiers = append(env.temporaryModifiers, modifier)
}

func (env *Environment) AddTankModifier(name string) {
	mod, exists := MakeModifierFromPropName(name)
	if exists {
		env.modifiers = append(env.modifiers, mod)
	}
}

func MakeModifierFromPropName(name string) (TankModifier, bool) {
	switch name {
	case "Castle":
		mod := TankModifier{Name: "Castle", EffectedParam: "PH", Amount: -2}
		return mod, true
	case "Log":
		mod := TankModifier{Name: "Log", EffectedParam: "PH", Amount: 2}
		return mod, true
	}
	return TankModifier{}, false
}
