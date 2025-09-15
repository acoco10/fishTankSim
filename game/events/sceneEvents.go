package events

import (
	"image"
)

type DayOver struct {
}

func (d DayOver) Type() string {
	return "DayOver"
}

type DayOverTransitionComplete struct {
}

func (d DayOverTransitionComplete) Type() string {
	return "DayOverTransitionComplete"
}

type NewDay struct {
	Day     int
	NTasks  int
	DayType string
}

type BedTime struct{}

func (b BedTime) Type() string { return "BedTime" }

func (n NewDay) Type() string {
	return "NewDay"
}

type FishTankLayout struct {
	image.Rectangle
}

func (f FishTankLayout) Type() string {
	return "FishTankLayout"
}

type LeavingFishScene struct {
}

func (l LeavingFishScene) Type() string {
	return "LeavingFishScene"
}

type Zoom struct{}

func (z Zoom) Type() string {
	return "Zoom"
}

type UnZoom struct{}

func (z UnZoom) Type() string {
	return "UnZoom"
}

type LightEvent struct {
	Day bool
}

func (l LightEvent) Type() string {
	return "LightEvent"
}
