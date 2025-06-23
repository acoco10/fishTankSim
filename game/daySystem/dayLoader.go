package daySystem

import "github.com/acoco10/fishTankWebGame/game/sceneManagement"

func LoadDaysTasks(log *sceneManagement.GameLog) {

	funcMap := make(map[int]func(log *sceneManagement.GameLog))

	funcMap[0] = LoadDay1Tasks
	funcMap[1] = LoadDay2Tasks
	funcMap[2] = LoadDay3Tasks

	if log.Day <= len(funcMap) {
		funcMap[log.Day-1](log)
	} else {
		println("task function not created yet for or day index is off", log.Day)
	}

}
