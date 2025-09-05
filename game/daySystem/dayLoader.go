package daySystem

import (
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/tasks"
)

var taskRelatedEvents []tasks.Event

func LoadDaysTasks(log *sceneManagement.GameLog) {

	LoadDefualtTasks(log)

}
