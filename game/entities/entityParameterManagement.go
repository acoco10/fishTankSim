package entities

import (
	"github.com/acoco10/fishTankWebGame/game/tasks"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"image/color"
)

type intParam uint8

const (
	//integers
	OtherX = iota
	OtherY
	Power
	Tries
	IndexCounter
	IndexCounterMax
	stopFrame
	AltSpriteHeight
	AltSpriteWidth
	lastZ
	UnsubId
	UnsubId2
)

const (
	//strings
	Position = iota
	Tag
	Text
	SpriteWText
	PublishString
)

const (
	//linkedEntities
	linkedCreature = iota
	LinkedGraphic1
	linkedGraphic2
)

const (
	//Events
	EventAtTime = iota
	PointEvent
)

const (
	// floats
	TextOffset = iota
	XEffectOffSet
	YEffectOffset
)

const (
	//rectangles
	TrashCan = iota
	Net
	TankBounds
	NetOpening
	CapturedBoundaries
)

const (
	//Buttonimages
	ButtonHighlight = iota
	ButtonImage
)

const (
	//buttons
	Button = iota
)

const (
	//Images
	Primary = iota
	Alternate
)

const (
	//stringLists
	Events = iota
)

type strParams [20]string

type EntityParameters struct {
	Floats           [20]float64
	Ints             [20]int
	EntIds           [10]uint32
	Flags            entFlags
	Strings          strParams
	Rectangles       [20]image.Rectangle
	Events           [20]tasks.Event
	UiButtonImages   [3]*widget.ButtonImage
	UiButtonWidget   *widget.Button
	Colors           [2]color.RGBA
	UiText           *widget.Text
	Images           [10]*ebiten.Image
	StringLists      [10][]string
	StringListsIndex [10]int
}
