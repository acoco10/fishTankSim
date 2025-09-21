package system

import (
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/graphics"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
	"math"
	"math/rand"
)

type modifierType = int8

const (
	Temporary modifierType = iota
	Permanent
)

const (
	pH       = "PH"
	Temp     = "temperature"
	phBoost  = "phBoost"
	phReduce = "phReduce"
	hotRock  = "hotRock"
	coolRock = "coolRock"

	significantly = "significantly"
	somewhat      = "somewhat"
	slightly      = "slightly"
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
	GraphicManager      *graphics.GraphicManager
	modConfigs          map[string]TempEnvModifierConfig
}

type TempEnvModifierConfig struct {
	averageAmount float64
	stdDevAmount  float64

	minimumDuration int
	maximumDuration int

	parameterEffected string
}

func SetupModConfigs() map[string]TempEnvModifierConfig {
	configMap := make(map[string]TempEnvModifierConfig)

	configMap[phBoost] = TempEnvModifierConfig{averageAmount: 0.7, stdDevAmount: 0.1, maximumDuration: 3, minimumDuration: 1, parameterEffected: pH}
	configMap[phReduce] = TempEnvModifierConfig{averageAmount: -0.7, stdDevAmount: 0.1, maximumDuration: 3, minimumDuration: 1, parameterEffected: pH}
	configMap[hotRock] = TempEnvModifierConfig{averageAmount: 5, stdDevAmount: 2, maximumDuration: 5, minimumDuration: 2, parameterEffected: Temp}
	configMap[coolRock] = TempEnvModifierConfig{averageAmount: -5, stdDevAmount: 2, maximumDuration: 5, minimumDuration: 2, parameterEffected: Temp}

	return configMap
}

func InitEnvironment(e *Environment) {
	e.Temperature = rand.Intn(20) + 55
	e.NaturalPHLevel = rand.Float64()*4 + 5
	e.ModifiedPHLevel = e.NaturalPHLevel
	e.modConfigs = SetupModConfigs()
	gm := graphics.GraphicManager{}
	gm.Timers = make(map[string]*util.Timer)
	gm.Timers[graphics.Trigger] = util.NewTimer(1.5)
	gm.Timers[graphics.Trigger].TimerUpdater = graphics.TriggerTimerUpdater
	e.GraphicManager = &gm

}

func (env *Environment) Subscribe(hub *tasks.EventHub) {
	hub.Subscribe(events.NewDay{}, func(e tasks.Event) {
		env.Temperature = rand.Intn(30) + 55
		filteredToRemoveOldMods := env.modifiers[:0]

		if len(env.modifiers) > 0 {
			env.GraphicManager.Timers[graphics.Trigger].TurnOn()
		}
		for _, mod := range env.modifiers {
			env.handleMod(mod)

			if mod.DaysActive < mod.Duration {
				filteredToRemoveOldMods = append(filteredToRemoveOldMods, mod)
			}
		}

		env.modifiers = filteredToRemoveOldMods
	})

	hub.Subscribe(events.ItemUsed{}, func(e tasks.Event) {
		ev := e.(events.ItemUsed)
		_, exists := env.modConfigs[ev.Name]
		if !exists {
			return
		}
		env.AppendModifier(env.MakeModifierFromConfig(ev.Name))
	})
}

func (env *Environment) handleMod(mod TankModifier) {
	switch mod.EffectedParam {
	case pH:
		env.ModifiedPHLevel += mod.Amount
	case Temp:
		env.ModifiedTemperature += int(mod.Amount)

	}
	msg := makeModMsg(mod)
	env.GraphicManager.Strings = append(env.GraphicManager.Strings, msg)
	mod.DaysActive++
	if mod.ModType == Temporary {
		mod.Amount = max(mod.Amount-mod.FallOff, 0)
		env.usedTempModifier = append(env.usedTempModifier, mod)
	}
}

func makeModMsg(mod TankModifier) string {
	modifierMsg := "increased"
	if mod.Amount < 0 {
		modifierMsg = "decreased"
	}
	var param = mod.EffectedParam
	amountAdj := getModifierAdjective(mod.Amount, param)
	msg := formModSentence(param, modifierMsg, amountAdj)
	return msg
}

func getModifierAdjective(amount float64, param string) string {
	var amountAdj string
	amt := math.Abs(amount)
	switch param {
	case pH:
		if amt > 0.5 {
			amountAdj = significantly
		} else if amount > 0.25 {
			amountAdj = somewhat
		} else {
			amountAdj = slightly
		}
	case Temp:
		if amt > 5 {
			amountAdj = significantly
		} else if amt > 3 {
			amountAdj = somewhat
		} else {
			amountAdj = slightly
		}
	}
	return amountAdj
}

func (env *Environment) MakeModifierFromConfig(name string) TankModifier {
	mod := env.modConfigs[name]
	duration := rand.Intn(mod.maximumDuration-mod.minimumDuration) + mod.minimumDuration
	amount := rand.NormFloat64()*mod.stdDevAmount + mod.averageAmount
	return TankModifier{Name: name, EffectedParam: mod.parameterEffected, Amount: amount, Duration: duration}
}

func formModSentence(param string, modifier string, adjective string) string {
	return util.UpperCaseFirstLetter(param) + " " + modifier + " " + adjective + "!"
}

func (env *Environment) AppendModifier(modifier TankModifier) {
	env.modifiers = append(env.modifiers, modifier)
}

func (env *Environment) AddTankModifier(name string) {
	mod, exists := env.MakeModifierFromPropName(name)
	if exists {
		env.modifiers = append(env.modifiers, mod)
	}
}

func (env *Environment) MakeModifierFromPropName(name string) (TankModifier, bool) {
	switch name {
	case "castle":
		mod := TankModifier{Name: "castle", EffectedParam: "PH", Amount: -2}
		return mod, true
	case "log":
		mod := TankModifier{Name: "log", EffectedParam: "PH", Amount: 2}
		return mod, true
	case hotRock, coolRock:
		env.AppendModifier(env.MakeModifierFromConfig(name))

	}

	return TankModifier{}, false
}
