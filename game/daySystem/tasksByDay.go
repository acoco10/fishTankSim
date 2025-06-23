package daySystem

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/geometry"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/tasks"
)

func LoadDay1Tasks(gameLog *sceneManagement.GameLog) {
	taskCondition := func(e tasks.Event) bool {
		ev, ok := e.(entities.CreatureReachedPoint)
		if ev.Point == nil {
			return false
		}
		return ok && ev.Point.PType == geometry.Food
	}

	gameTask := tasks.NewTask(entities.CreatureReachedPoint{}, "1. Feed your fish", taskCondition)
	gameTask.Subscribe(gameLog.GlobalEventHub)

	taskCondition2 := func(e tasks.Event) bool {
		ev, ok := e.(entities.SendData)
		return ok && ev.DataFor == "statsMenu"
	}

	gameTask2 := tasks.NewTask(entities.SendData{}, "2. Click your fish", taskCondition2)
	gameTask2.Subscribe(gameLog.GlobalEventHub)

	taskCondition3 := func(e tasks.Event) bool {
		_, ok := e.(entities.AllFishFed)
		return ok
	}

	gameTask3 := tasks.NewTask(entities.AllFishFed{}, "3. Feed them until they're full", taskCondition3)
	gameTask3.Subscribe(gameLog.GlobalEventHub)

	gameLog.Tasks = append(gameLog.Tasks, gameTask, gameTask2, gameTask3)
}

func LoadDay2Tasks(gameLog *sceneManagement.GameLog) {
	println("loading day 2 tasks")
	taskCondition := func(e tasks.Event) bool {
		_, ok := e.(events.MoneyAdded)
		return ok
	}

	gameTask := tasks.NewTask(events.MoneyAdded{}, "1. Collect Your Allowance", taskCondition)
	gameTask.Subscribe(gameLog.GlobalEventHub)

	taskCondition2 := func(e tasks.Event) bool {
		_, ok := e.(entities.AllFishFed)
		return ok
	}

	gameTask2 := tasks.NewTask(entities.AllFishFed{}, "2. Feed all your fish.", taskCondition2)
	gameTask2.Subscribe(gameLog.GlobalEventHub)

	gameLog.Tasks = []*tasks.Task{}
	gameLog.Tasks = append(gameLog.Tasks, gameTask, gameTask2)

	gameTask3 := tasks.NewTask(entities.AllFishFed{}, "3. Buy a new fish.", taskCondition2)
	gameTask2.Subscribe(gameLog.GlobalEventHub)

	gameLog.Tasks = []*tasks.Task{}
	gameLog.Tasks = append(gameLog.Tasks, gameTask, gameTask2, gameTask3)
}

func LoadDay3Tasks(gameLog *sceneManagement.GameLog) {
	println("loading day 2 tasks")
	taskCondition := func(e tasks.Event) bool {
		_, ok := e.(events.MoneyAdded)
		return ok
	}

	gameTask := tasks.NewTask(events.MoneyAdded{}, "1. Feed all your fish", taskCondition)
	gameTask.Subscribe(gameLog.GlobalEventHub)

	taskCondition2 := func(e tasks.Event) bool {
		_, ok := e.(entities.AllFishFed)
		return ok
	}

	gameTask2 := tasks.NewTask(entities.AllFishFed{}, "2. Buy a decoration for your fish tank", taskCondition2)
	gameTask2.CountRequired = 2
	gameTask2.Subscribe(gameLog.GlobalEventHub)

	taskCondition3 := func(e tasks.Event) bool {
		ev, ok := e.(entities.SendData)
		return ok && ev.DataFor == "statsMenu"
	}

	gameTask3 := tasks.NewTask(entities.SendData{}, "3. Do you chores", taskCondition3)
	gameTask3.Subscribe(gameLog.GlobalEventHub)

	gameLog.Tasks = []*tasks.Task{}
	gameLog.Tasks = append(gameLog.Tasks, gameTask, gameTask2, gameTask3)
}
