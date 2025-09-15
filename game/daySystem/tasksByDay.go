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

	/*	taskCondition12 := func(e tasks.Event) bool {
		ev, ok := e.(events.UnFocusEvent)
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
	}*/

	gameLog.TaskManager.NewTask(events.PHGuess{}, "1. Take a ph reading of your tank", taskCondition11)
	//gameTask.UIEffect = uiEvent

	taskCondition2 := func(e tasks.Event) bool {
		ev, ok := e.(events.ButtonClickedEvent)
		return ok && ev.ButtonText == "Confirm for prop select"
	}

	gameLog.TaskManager.NewTask(events.ButtonClickedEvent{}, "2. Pick your first tank decoration", taskCondition2)

	taskCondition4, text := FeedAllFishTask(3)
	gameLog.TaskManager.NewTask(entities.AllFishFed{}, text, taskCondition4)

	gameLog.DayType = sceneManagement.Free
}

func LoadDay2Tasks(gameLog *sceneManagement.GameLog) {
	println("loading day 2 tasks")

	taskCondition1 := func(e tasks.Event) bool {
		ev, ok := e.(events.ButtonClickedEvent)
		return ok && ev.ButtonText == "Go do your Chores?: Yes"
	}

	gameLog.TaskManager.NewTask(events.ButtonClickedEvent{}, "1. Do your chores", taskCondition1)

	taskCondition2 := func(e tasks.Event) bool {
		_, ok := e.(events.MoneyAdded)
		return ok
	}

	gameLog.TaskManager.NewTask(events.MoneyAdded{}, "2. Stash your allowance.", taskCondition2)

	taskCondition3 := func(e tasks.Event) bool {
		log.Printf("Day 2 purchase task condition met")
		ev, ok := e.(events.PurchaseSuccessful)
		return ok && entities.FishList(ev.Purchase) != ""
	}

	gameLog.TaskManager.NewTask(events.PurchaseSuccessful{}, "3. Buy a new fish.", taskCondition3)

	taskCondition4, text := FeedAllFishTask(4)
	gameLog.TaskManager.NewTask(entities.AllFishFed{}, text, taskCondition4)
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

	gameLog.TaskManager.NewTask(events.UISpriteAction{}, "1. Go to Camp", taskCondition1)

	taskCondition, text := FeedAllFishTask(2)
	gameLog.TaskManager.NewTask(entities.AllFishFed{}, text, taskCondition)

	gameLog.DayType = sceneManagement.Camp
}

func FeedAllFishTask(taskn int) (condition func(e tasks.Event) bool, text string) {

	taskCondition := func(e tasks.Event) bool {
		_, ok := e.(entities.AllFishFed)
		return ok
	}

	taskText := fmt.Sprintf("%d. Feed your fish ", taskn)

	return taskCondition, taskText

}

func LoadDefualtTasks(gameLog *sceneManagement.GameLog) {
	//uiEvent := events.UISpriteAction{UiSpriteData: "phreader", UiSpriteAction: "highlight"}
	taskCondition11 := func(e tasks.Event) bool {
		_, ok := e.(events.PHGuess)
		return ok
	}

	/*	taskCondition12 := func(e tasks.Event) bool {
		ev, ok := e.(events.UnFocusEvent)
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
	}*/

	gameLog.TaskManager.NewTask(events.PHGuess{}, "1. Take a ph reading", taskCondition11)
	//gameTask.UIEffect = uiEvent

	taskCondition4, text := FeedAllFishTask(2)
	gameLog.TaskManager.NewTask(entities.AllFishFed{}, text, taskCondition4)

	gameLog.DayType = sceneManagement.Free
}
