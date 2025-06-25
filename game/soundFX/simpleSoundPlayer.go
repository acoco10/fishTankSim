package soundFX

import (
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/hajimehoshi/ebiten/v2/audio"
	resource "github.com/quasilyte/ebitengine-resource"
	"log"
)

var l *resource.Loader

type SoundPlayer struct {
	*resource.Loader
	*audio.Player
	queue  []*audio.Player
	queue2 []*audio.Player
	timers map[resource.AudioID]*entities.Timer
}

func NewSoundPlayer() (*SoundPlayer, error) {
	if l == nil {
		loader, err := LoadSounds()
		if err != nil {
			return &SoundPlayer{}, err
		}
		l = loader
	}

	s := SoundPlayer{Loader: l, Player: &audio.Player{}}
	s.timers = make(map[resource.AudioID]*entities.Timer)
	s.timers[WhiteBoardMarker2] = entities.NewTimer(1.4)

	return &s, nil
}
func (s *SoundPlayer) Update() {
	s.queue = UpdateQueue(s.queue)
	s.queue2 = UpdateQueue(s.queue2)
}
func UpdateQueue(queue []*audio.Player) []*audio.Player {
	if len(queue) > 0 {
		if queue[0].IsPlaying() {
			return queue
		}
	}
	if len(queue) > 0 {
		queue = queue[1:]
	}
	if len(queue) > 0 {
		queue[0].Play()
	}
	return queue
}

func (s *SoundPlayer) AddToQueue(id resource.AudioID, queue int) {

	if queue == 1 {
		sfx := s.Loader.LoadWAV(id).Player
		sfx.SetVolume(s.Loader.LoadAudio(id).Volume)
		err := sfx.Rewind()
		if err != nil {
			log.Printf("%q Rewind: %s", id, err)
		}
		s.queue = append(s.queue, sfx)
		if len(s.queue) == 1 {
			sfx.Play()
		}
	}
	if queue == 2 {
		sfx := s.Loader.LoadWAV(id).Player
		sfx.SetVolume(s.Loader.LoadAudio(id).Volume)
		err := sfx.Rewind()
		if err != nil {
			log.Printf("%q Rewind: %s", id, err)
		}
		s.queue2 = append(s.queue2, sfx)
		if len(s.queue2) == 1 {
			sfx.Play()
		}
	}
}

func (s *SoundPlayer) Play(id resource.AudioID) {

	sfx := s.Loader.LoadWAV(id).Player
	err := sfx.Rewind()
	if err != nil {
		log.Printf("%q Rewind: %s", id, err)
	}
	s.Player = sfx
	println("volume for", id, s.Loader.LoadAudio(id).Volume)
	s.Player.SetVolume(s.Loader.LoadAudio(id).Volume)

	sfx.Play()
}

func (s *SoundPlayer) Pause() {
	s.Player.Pause()
}
