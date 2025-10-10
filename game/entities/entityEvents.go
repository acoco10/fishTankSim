package entities

import (
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/acoco10/fishTankWebGame/game/util"
)

type CreatureReachedPoint struct {
	PointTypeReached util.InterestPoint
	PointID          uint32
	CreatureID       uint32
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

type WriteToWhiteBoard struct {
	Msg               string
	PreferredPosition uint8
	NoErase           bool
	Later             bool
	EventDriven       tasks.Event
	EventToPublish    []tasks.Event
	EraseAfterFlag    bool
	Insets            [2]float64 //this is set from preferred position once we get to whiteboard,
	// not super intuitive but allows to reuse the data structures
}

func (w WriteToWhiteBoard) Type() string {
	return "WriteToWhiteBoard"
}

type DisableWhiteBoard struct {
	UnLockEvent tasks.Event
	Condition   func(event tasks.Event) bool
}

func (dwb DisableWhiteBoard) Type() string {
	return "DisableWhiteBoard"
}

type WhiteBoardErased struct {
	When uint8
}

func (wbc WhiteBoardErased) Type() string {
	return "WhiteBoard"
}

type EraseRequest struct {
	time    uint8
	onClick bool
}

func (er EraseRequest) Type() string {
	return "EraseRequest"
}

type RequestZoom struct {
	Reason ZoomState
}

func (z RequestZoom) Type() string {
	return "RequestZoom"
}

type DebrisCaptured struct {
	Id uint32
}

func (z DebrisCaptured) Type() string {
	return "DebrisCaptured"
}
