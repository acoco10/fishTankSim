package events

type ButtonEvent struct {
	ButtonText string
	EType      string
}

func (b ButtonEvent) Type() string {
	return "ButtonEvent"
}

type ButtonClickedEvent struct {
	ButtonText string
}

func (b ButtonClickedEvent) Type() string {
	return "ButtonClickedEvent"
}

type UISpriteAction struct {
	UiSprite       string
	UiSpriteAction string
}

func (U UISpriteAction) Type() string {
	return "UISpriteAction"
}

type CloseWindow struct {
	OverRide bool
}

func (c CloseWindow) Type() string {
	return "CloseWindow"
}

type WindowClosed struct{}

func (w WindowClosed) Type() string {
	return "WindowClosed"
}

type AddWindow struct {
	CursorOccupied bool
}

func (a AddWindow) Type() string {
	return "AddWindow"
}

type WindowOpened struct{}

func (w WindowOpened) Type() string {
	return "WindowOpened"
}

type PHGuess struct {
	Guess float64
}

func (P PHGuess) Type() string {
	return "PHGuess"
}
