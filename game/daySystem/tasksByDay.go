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

	//uiEvent := events.UISpriteAction{UiSpriteData: "phreader", UiSpriteAction: "highlight"}
	taskCondition11 := func(e tasks.Event) bool {
		_, ok := e.(events.PHGuess)
		return ok
	}

	taskCondition12 := func(e tasks.Event) bool {
		ev, ok := e.(events.UnFocus)
		if !ok {
			return false
		}
		ent, exists := entities.GetEntity(ev.EntID)
		if !exists {
			return false
		}
		if ent.UiData != nil {
			if ent.UiData.Label == "phreader" {
				return true
			}
		}
		return false
	}
	gameTask := tasks.NewTask(events.UnFocus{}, "1. Take a ph reading of your tank", taskCondition12, gameLog.GlobalEventHub)
	//gameTask.UIEffect = uiEvent
	subT1 := &tasks.SubTask{Completed: false, Condition: taskCondition11, EventType: events.PHGuess{}}
	gameTask.SubTasks = append(gameTask.SubTasks, subT1)

	taskCondition2 := func(e tasks.Event) bool {
		ev, ok := e.(events.ButtonClickedEvent)
		return ok && ev.ButtonText == "Confirm for prop select"
	}

	gameTask2 := tasks.NewTask(events.ButtonClickedEvent{}, "2. Pick your first tank decoration", taskCondition2, gameLog.GlobalEventHub)

	gameLog.Tasks = []*tasks.Task{gameTask, gameTask2, FeedAllFishTask(3, gameLog.GlobalEventHub)}

	gameLog.DayType = sceneManagement.Free
}

func LoadDay2Tasks(gameLog *sceneManagement.GameLog) {
	println("loading day 2 tasks")

	taskCondition1 := func(e tasks.Event) bool {
		ev, ok := e.(events.ButtonClickedEvent)
		return ok && ev.ButtonText == "Go do your Chores?: Yes"
	}

	gameTask := tasks.NewTask(events.ButtonClickedEvent{}, "1. Do your chores", taskCondition1, gameLog.GlobalEventHub)

	taskCondition2 := func(e tasks.Event) bool {
		_, ok := e.(events.MoneyAdded)
		return ok
	}

	gameTask2 := tasks.NewTask(events.MoneyAdded{}, "2. Stash your allowance.", taskCondition2, gameLog.GlobalEventHub)

	taskCondition3 := func(e tasks.Event) bool {
		log.Printf("Day 2 purchase task condition met")
		ev, ok := e.(events.PurchaseSuccessful)
		return ok && entities.FishList(ev.Purchase) != ""
	}

	gameTask3 := tasks.NewTask(events.PurchaseSuccessful{}, "3. Buy a new fish.", taskCondition3, gameLog.GlobalEventHub)

	gameTask4 := FeedAllFishTask(4, gameLog.GlobalEventHub)

	gameLog.Tasks = []*tasks.Task{gameTask, gameTask2, gameTask3, gameTask4}
	gameLog.DayType = sceneManagement.Chores
}

func LoadDay3Tasks(gameLog *sceneManagement.GameLog) {
	println("loading day 3 tasks")
	/*	taskCondition := func(e tasks.Event) bool {
			_, ok := e.(events.MoneyAdded)
			return ok
		}

		gameTask2 := tasks.NewTask(entities.AllFishFed{}, "2. Buy a decoration for your fish tank", taskCondition2)*/

	taskCondition1 := func(e tasks.Event) bool {
		ev, ok := e.(events.UISpriteAction)
		return ok && ev.UiSprite == "door"
	}

	gameTask1 := tasks.NewTask(events.UISpriteAction{}, "1. Go to Camp", taskCondition1, gameLog.GlobalEventHub)

	gameLog.Tasks = []*tasks.Task{}
	gameLog.Tasks = append(gameLog.Tasks, gameTask1, FeedAllFishTask(2, gameLog.GlobalEventHub))

	gameLog.DayType = sceneManagement.Camp
}

func FeedAllFishTask(taskn int, hub *tasks.EventHub) *tasks.Task {

	taskCondition := func(e tasks.Event) bool {
		_, ok := e.(entities.AllFishFed)
		return ok
	}

	text := fmt.Sprintf("%d. Feed all your fish ", taskn)

	task := tasks.NewTask(entities.AllFishFed{}, text, taskCondition, hub)
	task.Type = tasks.FishFed

	return task

}
