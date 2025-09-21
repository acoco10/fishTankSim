package events

type FishLevelUp struct {
	FishID uint32
}

func (f FishLevelUp) Type() string {
	return "FishLevelUp"
}

type UnFocusEvent struct {
	EntID uint32
}

func (u UnFocusEvent) Type() string {
	return "UnFocusEvent"
}

type FocusEvent struct {
	EntID uint32
}

func (f FocusEvent) Type() string {
	return "FocusEvent"
}

type NewProp struct {
	PropId uint32
	Name   string
}

func (w NewProp) Type() string {
	return "NewProp"
}

type ItemUsed struct {
	Name string
}

func (f ItemUsed) Type() string {
	return "ItemUsed"
}

type PlacementMode struct {
}

func (p PlacementMode) Type() string {
	return "PlacementMode"
}

type WritingToWhiteBoard struct {
	Msg string
}

func (wtw WritingToWhiteBoard) Type() string {
	return "WritingToWhiteBoard"
}
