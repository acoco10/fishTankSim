package soundFX

import (
	"bytes"
	"github.com/acoco10/fishTankWebGame/assets"
	resource "github.com/quasilyte/ebitengine-resource"
	"io"
	"log"
	"sort"
)

var loadedOGGs *resource.Loader

var oggData = map[string][]byte{}

// ReadSeekCloser wraps bytes.Reader to implement the full interface
type ReadSeekCloser struct {
	*bytes.Reader
}

func (rsc *ReadSeekCloser) Close() error {
	return nil
}

func LoadOggs() {
	volumeMap := map[resource.AudioID]float64{
		BestAdventureEver: 0.5,
		CardBoard:         0.7,
		Coins1:            0.2,
		Crash:             0.1,
		DayTimeJazz:       -0.5,
		FailedStart:       0.9,
		IndieCafe:         0.0,
		Kaching:           0.3,
		MoneyCounter:      0.5,
		MowerRunning:      -0.2,
		PickUpOne:         0.0,
		PlopSound:         -0.3,
		PouringFood:       0.0,
		SelectSound:       -0.5,
		SuccessMusic:      -0.8,
		TropicalHouse:     -0.5,
		WaterBubbles:      0.0,
		SelectSound2:      0.0,
		WhiteBoardMarker1: 0.5,
		WhiteBoardMarker2: 0.5,
	}

	soundDir, err := assets.OGG.ReadDir("soundFx/oggs")
	if err != nil {
		log.Fatal("Error reading sound file directory")
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

		song, err2 := assets.OGG.ReadFile("soundFx/oggs/" + name)
		if err2 != nil {
			log.Fatal("Error reading sound files into song data", name)
		}

		oggData[name] = song // Fixed: now storing in oggData
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

	// Fixed: Use oggData and proper ReadSeekCloser
	l.OpenAssetFunc = func(path string) io.ReadCloser {
		data, exists := oggData[path]
		if !exists {
			log.Printf("Sound file not found: %s", path)
			return nil
		}
		if len(data) == 0 {
			log.Printf("Empty sound file: %s", path)
			return nil
		}
		return &ReadSeekCloser{bytes.NewReader(data)}
	}

	loadedOGGs = l
}
