package util

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
)

func UpdateTimers(timers map[any]*entities.Timer) {
	for _, timer := range timers {
		timer.Update()
	}
}
