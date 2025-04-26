package soundFX

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	resource "github.com/quasilyte/ebitengine-resource"
	"log"
)

type SongPlayer struct {
	*resource.Loader
	*audio.Player
}

func NewSongPlayer() (*SongPlayer, error) {
	l, err := LoadSounds()
	if err != nil {
		return &SongPlayer{}, err
	}

	s := SongPlayer{l, &audio.Player{}}
	return &s, nil
}

func (s SongPlayer) Play(id resource.AudioID) {
	sfx := s.Loader.LoadWAV(id).Player

	err := sfx.Rewind()
	if err != nil {
		log.Printf("%q Rewind: %s", id, err)
	}
	sfx.Play()
}
