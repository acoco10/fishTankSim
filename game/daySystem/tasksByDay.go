package daySystem

import (
	"fmt"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/events"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"log"
)

func LoadDay1Tasks(gameLog *sceneManagement.GameLog) {

	taskCondition := func(e tasks.Event) bool {
		ev, ok := e.(events.UISpriteAction)
		return ok && ev.UiSpriteAction == "ph reading"
	}

	uiEvent := events.UISpriteAction{UiSprite: "phreader", UiSpriteAction: "highlight"}

	gameTask := tasks.NewTask(events.UISpriteAction{}, "1. Take a ph reading of your tank", taskCondition)
	gameTask.UIEffect = uiEvent

	taskCondition2 := func(e tasks.Event) bool {
		ev, ok := e.(entities.SendData)
		return ok && ev.DataFor == "statsMenu"
	}

	gameTask2 := tasks.NewTask(entities.SendData{}, "2. Pick your first tank decoration", taskCondition2)

	taskCondition3 := func(e tasks.Event) bool {
		_, ok := e.(entities.AllFishFed)
		return ok
	}

	gameTask3 := tasks.NewTask(entities.AllFishFed{}, "3. Feed your fish!", taskCondition3)

	gameLog.Tasks = append(gameLog.Tasks, gameTask, gameTask2, gameTask3)
	gameLog.DayType = sceneManagement.Camp
}

func LoadDay2Tasks(gameLog *sceneManagement.GameLog) {
	println("loading day 2 tasks")

	taskCondition1 := func(e tasks.Event) bool {
		ev, ok := e.(events.ButtonClickedEvent)
		return ok && ev.ButtonText == "Go do your Chores?: Yes"
	}

	gameTask := tasks.NewTask(events.MoneyAdded{}, "1. Do your chores", taskCondition1)

	taskCondition2 := func(e tasks.Event) bool {
		_, ok := e.(events.MoneyAdded)
		return ok
	}

	gameTask2 := tasks.NewTask(events.NewPurchase{}, "2. Stash your allowance.", taskCondition2)

	taskCondition3 := func(e tasks.Event) bool {
		log.Printf("Day 2 purchase task condition met")
		_, ok := e.(events.NewPurchase)
		return ok
	}

	gameTask3 := tasks.NewTask(events.NewPurchase{}, "3. Buy a new fish.", taskCondition3)

	gameTask4 := FeedAllFishTask(4)

	gameLog.Tasks = []*tasks.Task{}
	gameLog.Tasks = append(gameLog.Tasks, gameTask, gameTask2, gameTask3, gameTask4)
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

	taskCondition3 := func(e tasks.Event) bool {
		ev, ok := e.(entities.SendData)
		return ok && ev.DataFor == "statsMenu"
	}

	gameTask3 := tasks.NewTask(entities.SendData{}, "3. Do you chores", taskCondition3)

	gameLog.Tasks = []*tasks.Task{}
	gameLog.Tasks = append(gameLog.Tasks, gameTask, gameTask2, gameTask3)
	gameLog.DayType = sceneManagement.Free
}

func FeedAllFishTask(taskn int) *tasks.Task {

	taskCondition := func(e tasks.Event) bool {
		_, ok := e.(entities.AllFishFed)
		return ok
	}

	text := fmt.Sprintf("%d. Feed all your fish ", taskn)

	task := tasks.NewTask(entities.AllFishFed{}, text, taskCondition)
	return task

}
