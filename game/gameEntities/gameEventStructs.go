package gameEntities

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

type PointGenerated struct {
	Point *Point
}

type SpriteHovered struct {
	sprite *Sprite
}
