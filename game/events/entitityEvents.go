package events

type FishLevelUp struct {
	FishID uint32
}

func (f FishLevelUp) Type() string {
	return "FishLevelUp"
}

type UnFocus struct {
	EntID uint32
}

func (u UnFocus) Type() string {
	return "UnFocus"
}

type Focus struct {
	EntID uint32
}

func (f Focus) Type() string {
	return "Focus"
}

type WriteToWhiteBoard struct {
	Msg               string
	PreferredPosition string
}

func (w WriteToWhiteBoard) Type() string {
	return "WriteToWhiteBoard"
}

type NewProp struct {
	PropId uint32
	Name   string
}

func (w NewProp) Type() string {
	return "NewProp"
}
