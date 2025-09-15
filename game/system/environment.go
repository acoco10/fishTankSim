package system

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"log"
	"math/rand"
)

type modifierType = int8

const (
	Temporary modifierType = iota
	Permanent
)

type TankModifier struct {
	Name          string
	EffectedParam string
	Amount        float64
	Duration      int
	DaysActive    int
	ModType       modifierType
	FallOff       float64
}
type Environment struct {
	Temperature         int
	ModifiedTemperature int
	NaturalPHLevel      float64
	ModifiedPHLevel     float64
	NatureBoosters      int
	modifiers           []TankModifier
	usedTempModifier    []TankModifier
}

func InitEnvironment(e *Environment) {
	e.Temperature = rand.Intn(10) + 65
	e.NaturalPHLevel = rand.Float64()*5 + 3
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

		for i, mod := range env.modifiers {
			switch mod.EffectedParam {
			case "PH":
				graphics.NewFadeInTextGraphicCentered("PH temporarily increased", 180)
				mod.DaysActive++
				env.ModifiedPHLevel += mod.Amount
			case "Temperature":
				graphics.NewFadeInTextGraphicCentered("temperature temporarily increased", 180)
				mod.DaysActive++
				env.ModifiedTemperature += int(mod.Amount)
			}
			if mod.ModType == Temporary {
				mod.Amount = max(mod.Amount-mod.FallOff, 0)
				env.usedTempModifier = append(env.usedTempModifier, mod)
				if mod.DaysActive >= mod.Duration {
					if len(env.modifiers) > 1 {
						env.modifiers = append(env.modifiers[0:i], env.modifiers[i+1:]...)
					} else {
						env.modifiers = []TankModifier{}
					}
				}
			}
		}
	})

	hub.Subscribe(events.ItemUsed{}, func(e tasks.Event) {
		ev := e.(events.ItemUsed)
		if ev.Name == "phBoost" {
			log.Printf("Modifier added to tank environment")
			env.AppendModifier(TankModifier{Name: ev.Name, EffectedParam: "PH", Amount: .25, Duration: 1})
		}
		if ev.Name == "phReduce" {
			log.Printf("Modifier added to tank environment")
			env.AppendModifier(TankModifier{Name: ev.Name, EffectedParam: "PH", Amount: -.25, Duration: 1})
		}
	})

}

func (env *Environment) AppendModifier(modifier TankModifier) {
	env.modifiers = append(env.modifiers, modifier)
}

func (env *Environment) AddTankModifier(name string) {
	mod, exists := MakeModifierFromPropName(name)
	if exists {
		env.modifiers = append(env.modifiers, mod)
	}
}

func MakeModifierFromPropName(name string) (TankModifier, bool) {
	switch name {
	case "castle":
		mod := TankModifier{Name: "castle", EffectedParam: "PH", Amount: -2}
		return mod, true
	case "log":
		mod := TankModifier{Name: "log", EffectedParam: "PH", Amount: 2}
		return mod, true
	case "hotRock":
		mod := TankModifier{Name: "hotRock", EffectedParam: "Temperature", Amount: 5, Duration: 3}
		return mod, true
	case "coolRock":
		mod := TankModifier{Name: "coolRock", EffectedParam: "Temperature", Amount: -5, Duration: 3}
		return mod, true
	}

	return TankModifier{}, false
}
