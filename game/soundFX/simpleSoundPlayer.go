package soundFX

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	resource "github.com/quasilyte/ebitengine-resource"
	"log"
)

type SoundPlayer struct {
	*resource.Loader
	*audio.Player
	queue []*audio.Player
}

func NewSoundPlayer() (*SoundPlayer, error) {
	l, err := LoadSounds()
	if err != nil {
		return &SoundPlayer{}, err
	}

	s := SoundPlayer{Loader: l, Player: &audio.Player{}}
	return &s, nil
}

func (s *SoundPlayer) Update() {
	playNext := false
	if len(s.queue) > 0 {
		if !s.queue[0].IsPlaying() {
			s.queue = s.queue[1:]
			playNext = true
		}
	}
	if playNext && len(s.queue) > 0 {
		s.queue[0].Play()
	}
}

func (s *SoundPlayer) AddToQueue(id resource.AudioID) {
	sfx := s.Loader.LoadWAV(id).Player
	err := sfx.Rewind()
	if err != nil {
		log.Printf("%q Rewind: %s", id, err)
	}
	s.queue = append(s.queue, sfx)
}

func (s *SoundPlayer) Play(id resource.AudioID) {
	sfx := s.Loader.LoadWAV(id).Player
	err := sfx.Rewind()
	if err != nil {
		log.Printf("%q Rewind: %s", id, err)
	}

	s.Player = sfx
	sfx.Play()
}

func (s *SoundPlayer) Pause() {
	s.Player.Pause()
}
