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

type SendData struct {
	DataFor string
	Data    string
}

type RequestData struct {
	DataType   string
	RequestFor any
}

type UISpriteAction struct {
	UiSprite       string
	UiSpriteAction string
}

type ButtonEvent struct {
	ButtonText string
	EType      string
}
