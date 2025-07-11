package daySystem

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/geometry"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"log"
)

func LoadDay1Tasks(gameLog *sceneManagement.GameLog) {

	taskCondition := func(e tasks.Event) bool {
		ev, ok := e.(entities.CreatureReachedPoint)
		if ev.Point == nil {
			return false
		}
		return ok && ev.Point.PType == geometry.Food
	}

	uiEvent := events.UISpriteAction{UiSprite: "fishFood", UiSpriteAction: "highlight"}

	gameTask := tasks.NewTask(entities.CreatureReachedPoint{}, "1. Feed your fish", taskCondition)
	gameTask.UIEffect = uiEvent

	taskCondition2 := func(e tasks.Event) bool {
		ev, ok := e.(entities.SendData)
		return ok && ev.DataFor == "statsMenu"
	}

	gameTask2 := tasks.NewTask(entities.SendData{}, "2. Click your fish", taskCondition2)

	taskCondition3 := func(e tasks.Event) bool {
		_, ok := e.(entities.AllFishFed)
		return ok
	}

	gameTask3 := tasks.NewTask(entities.AllFishFed{}, "3. Feed them until they're full", taskCondition3)

	gameLog.Tasks = append(gameLog.Tasks, gameTask, gameTask2, gameTask3)
	gameLog.DayType = sceneManagement.Free
}

func LoadDay2Tasks(gameLog *sceneManagement.GameLog) {
	println("loading day 2 tasks")
	taskCondition := func(e tasks.Event) bool {
		_, ok := e.(events.MoneyAdded)
		return ok
	}

	gameTask := tasks.NewTask(events.MoneyAdded{}, "1. Do Your Chores", taskCondition)

	gameTask2 := FeedAllFishTask()

	taskCondition3 := func(e tasks.Event) bool {
		log.Printf("Day 2 purchase task condition met")
		_, ok := e.(events.NewPurchase)
		return ok
	}

	gameTask3 := tasks.NewTask(events.NewPurchase{}, "3. Buy a new fish.", taskCondition3)

	gameLog.Tasks = []*tasks.Task{}
	gameLog.Tasks = append(gameLog.Tasks, gameTask, gameTask2, gameTask3)
	gameLog.DayType = sceneManagement.Chores
}

func LoadDay3Tasks(gameLog *sceneManagement.GameLog) {
	println("loading day 2 tasks")
	taskCondition := func(e tasks.Event) bool {
		_, ok := e.(events.MoneyAdded)
		return ok
	}

	gameTask := tasks.NewTask(events.MoneyAdded{}, "1. Feed all your fish", taskCondition)

	taskCondition2 := func(e tasks.Event) bool {
		_, ok := e.(entities.AllFishFed)
		return ok
	}

	gameTask2 := tasks.NewTask(entities.AllFishFed{}, "2. Buy a decoration for your fish tank", taskCondition2)
	gameTask2.CountRequired = 2

	taskCondition3 := func(e tasks.Event) bool {
		ev, ok := e.(entities.SendData)
		return ok && ev.DataFor == "statsMenu"
	}

	gameTask3 := tasks.NewTask(entities.SendData{}, "3. Do you chores", taskCondition3)

	gameLog.Tasks = []*tasks.Task{}
	gameLog.Tasks = append(gameLog.Tasks, gameTask, gameTask2, gameTask3)
	gameLog.DayType = sceneManagement.Free
}

func FeedAllFishTask() *tasks.Task {

	taskCondition := func(e tasks.Event) bool {
		_, ok := e.(entities.AllFishFed)
		return ok
	}

	task := tasks.NewTask(entities.AllFishFed{}, "2. Feed all your fish.", taskCondition)
	return task

}
