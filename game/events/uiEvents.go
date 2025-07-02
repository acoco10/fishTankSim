package events

type ButtonEvent struct {
	ButtonText string
	EType      string
}

type ButtonClickedEvent struct {
	ButtonText string
}

type UISpriteAction struct {
	UiSprite       string
	UiSpriteAction string
}

type CloseWindow struct {
}
