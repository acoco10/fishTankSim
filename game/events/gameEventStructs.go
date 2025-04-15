package events

import "image"

type PropertyUpdate struct {
	Property string
	Value    bool
}

type DialogueEvent struct {
	Characters []string
}

type ButtonClickedEvent struct {
	ButtonText string
}

type CreatureReachedPoint struct {
	Point image.Point
}

type MouseButtonPressed struct {
	Point image.Point
}
