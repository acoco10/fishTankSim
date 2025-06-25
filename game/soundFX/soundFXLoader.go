package soundFX

import (
	"bytes"
	"fmt"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/hajimehoshi/ebiten/v2/audio"
	resource "github.com/quasilyte/ebitengine-resource"
	"io"
	"sort"
)

const (
	BestAdventureEver resource.AudioID = iota
	PickUpOne         resource.AudioID = iota
	PlopSound         resource.AudioID = iota
	PouringFood       resource.AudioID = iota
	SelectSound       resource.AudioID = iota
	SuccessMusic      resource.AudioID = iota
	WaterBubbles      resource.AudioID = iota
	SelectSound2      resource.AudioID = iota
	WhiteBoardMarker1 resource.AudioID = iota
	WhiteBoardMarker2 resource.AudioID = iota
)

var audioContext = audio.NewContext(44100)

var SoundData = map[string][]byte{}

// key is file path, not just name
func LoadSounds() (*resource.Loader, error) {
	var rLoader *resource.Loader

	volumeMap := map[resource.AudioID]float64{
		BestAdventureEver: -0.5,
		PickUpOne:         0.0,
		PlopSound:         -0.3,
		PouringFood:       0.0,
		SelectSound:       -0.5,
		SuccessMusic:      -0.6,
		WaterBubbles:      0.0,
		SelectSound2:      0.0,
		WhiteBoardMarker1: 0.5,
		WhiteBoardMarker2: 0.5,
	}

	soundDir, err := assets.SoundDir.ReadDir("soundFx")

	if err != nil {
		return rLoader, fmt.Errorf("error reading sound files: %w", err)
	}

	sort.Slice(soundDir, func(i, j int) bool {

		return soundDir[i].Name()[0] < soundDir[j].Name()[0]
	})

	audioRegMap := map[resource.AudioID]resource.AudioInfo{}

	for i, dir := range soundDir {
		name := dir.Name()
		endIndex := len(name) - 4
		sName := name[:endIndex]
		println(i, "Loading sound:", sName)

		song, err := assets.SoundDir.ReadFile("soundFx/" + name)
		if err != nil {
			return rLoader, fmt.Errorf("error reading sound file: %w", err)
		}

		SoundData[name] = song
		println("saving audio id:", resource.AudioID(i))
		vol := 0.0

		if volumeMap[resource.AudioID(i)] != 0.0 {
			vol = volumeMap[resource.AudioID(i)]
		}

		println(resource.AudioID(i), "volume =", vol)
		audioRegMap[resource.AudioID(i)] = resource.AudioInfo{Path: name, Volume: vol}

	}

	l := resource.NewLoader(audioContext)
	l.AudioRegistry.Assign(audioRegMap)
	l.OpenAssetFunc = func(path string) io.ReadCloser {
		return io.NopCloser(bytes.NewReader(SoundData[path]))
	}
	return l, nil

}
