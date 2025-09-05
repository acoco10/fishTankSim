package entities

type CreatureReachedPoint struct {
	PointID    uint32
	CreatureID uint32
}

func (c CreatureReachedPoint) Type() string {
	return "CreatureReachedPoint"
}

type PointGenerated struct {
	PointId uint32
	Source  string
}

func (p PointGenerated) Type() string {
	return "PointGenerated"
}

type AllFishFed struct {
}

func (a AllFishFed) Type() string {
	return "AllFishFed"
}

type SendData struct {
	DataFor string
	Data    string
}

func (s SendData) Type() string {
	return "SendData"
}

type RequestData struct {
	DataType   string
	RequestFor any
}

func (r RequestData) Type() string {
	return "RequestData"
}

type FishEvent struct {
	fish  *CreatureData
	event string
}

func (f FishEvent) Type() string {
	return "FishEvent"
}

type TurnOnBubbles struct {
}

func (f TurnOnBubbles) Type() string {
	return "TurnOnBubbles"
}

type TurnOffBubbles struct {
}

func (f TurnOffBubbles) Type() string {
	return "TurnOffBubbles"
}

type PlacementPicked struct {
	PlacementFor string
	X            float32
	Y            float32
	Z            int
}

func (p PlacementPicked) Type() string {
	return "PlacementPicked"
}
