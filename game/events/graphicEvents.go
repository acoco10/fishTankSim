package events

type ClickMeGraphicEvent struct {
	X, Y        float64
	SpriteWidth float64
}

type FadeInTextEvent struct {
	X, Y float64
	Text string
}

type TurnOffGraphic struct {
	X, Y float64
}

type GraphicCreated struct {
	GraphicType string
	GraphicId   int
}
