package events

import "image"

type DayOver struct {
}

type DayOverTransitionComplete struct {
}

type NewDay struct {
	NTasks int
}

type FishTankLayout struct {
	image.Rectangle
}
