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
	temporary     bool
}
type Environment struct {
	Temperature     int
	NaturalPHLevel  float64
	modifiedPHlevel float64
	NatureBoosters  int
	modifiers       []TankModifier
}

func InitEnvironment(e *Environment) {
	e.Temperature = rand.Intn(10) + 67
	e.NaturalPHLevel = rand.Float64()*3 + 6
}

func (env *Environment) Subscribe(hub *tasks.EventHub) {
	hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		env.Temperature = rand.Intn(10) + 67

		for _, mod := range env.modifiers {
			switch mod.EffectedParam {
			case "PH":
				if env.NaturalPHLevel != env.NaturalPHLevel+mod.Amount {
				}
				util.Lerp64(env.NaturalPHLevel, env.NaturalPHLevel+mod.Amount, mod.Amount)
			}
		}

	})
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
		mod := TankModifier{Name: "Castle"}
		return mod, true
	case "Log":
		mod := TankModifier{Name: "Log", EffectedParam: "PH", Amount: 2}
		return mod, true
	}
	return TankModifier{}, false
}
