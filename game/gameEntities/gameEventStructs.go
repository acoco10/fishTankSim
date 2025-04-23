package gameEntities

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
	Point    *Point
	Creature *Creature
}

type MouseButtonPressed struct {
	Point *Point
}

type SpriteHovered struct {
	sprite *Sprite
}
