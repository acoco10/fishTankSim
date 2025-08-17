package soundFX

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	resource "github.com/quasilyte/ebitengine-resource"
)

var audioContext = audio.NewContext(48000)
var loadedSounds *resource.Loader

const (
	BestAdventureEver resource.AudioID = iota
	CardBoard         resource.AudioID = iota
	Coins1            resource.AudioID = iota
	Crash             resource.AudioID = iota
	DayTimeJazz       resource.AudioID = iota
	FailedStart       resource.AudioID = iota
	IndieCafe         resource.AudioID = iota
	Kaching           resource.AudioID = iota
	MoneyCounter      resource.AudioID = iota
	MowerRunning      resource.AudioID = iota
	PickUpOne         resource.AudioID = iota
	PlopSound         resource.AudioID = iota
	PouringFood       resource.AudioID = iota
	SelectSound       resource.AudioID = iota
	SuccessMusic      resource.AudioID = iota
	TropicalHouse     resource.AudioID = iota
	WaterBubbles      resource.AudioID = iota
	SelectSound2      resource.AudioID = iota
	WhiteBoardMarker1 resource.AudioID = iota
	WhiteBoardMarker2 resource.AudioID = iota
)
