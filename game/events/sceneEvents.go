package events

import (
	"image"
)

type DayOver struct {
}

type DayOverTransitionComplete struct {
}

type NewDay struct {
	NTasks int
	Type   string
}

type FishTankLayout struct {
	image.Rectangle
}

type IntroFlow struct {
	Event string
}
