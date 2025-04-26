package soundFX

import (
	"bytes"
	"fishTankWebGame/assets"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/audio"
	resource "github.com/quasilyte/ebitengine-resource"
	"io"
)

const (
	JazzE resource.AudioID = iota
	SunSetVibe
	WaterBubbles
)

var audioContext = audio.NewContext(44100)

var SoundData = map[string][]byte{}

// key is file path, not just name
func LoadSounds() (*resource.Loader, error) {

	var rLoader *resource.Loader

	soundDir, err := assets.SoundDir.ReadDir("soundFx")

	if err != nil {
		return rLoader, fmt.Errorf("error reading sound files: %w", err)
	}

	audioRegMap := map[resource.AudioID]resource.AudioInfo{}

	for i, dir := range soundDir {
		name := dir.Name()
		endIndex := len(name) - 4
		sName := name[:endIndex]
		println(i, "Loading sound:", sName)
		song, err := assets.SoundDir.ReadFile(name)

		if err != nil {
			return rLoader, fmt.Errorf("error reading sound file: %w", err)
		}

		SoundData[name] = song
		audioRegMap[resource.AudioID(i)] = resource.AudioInfo{Path: name, Volume: 0.2}

	}

	l := resource.NewLoader(audioContext)

	l.OpenAssetFunc = func(path string) io.ReadCloser {
		return io.NopCloser(bytes.NewReader(SoundData[path]))
	}
	return l, nil

}
