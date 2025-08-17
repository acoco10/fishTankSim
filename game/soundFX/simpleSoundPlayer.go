package soundFX

import (
	"github.com/acoco10/fishTankWebGame/game/util"
	"github.com/hajimehoshi/ebiten/v2/audio"
	resource "github.com/quasilyte/ebitengine-resource"
	"log"
	"time"
)

type SoundPlayer struct {
	sounds      map[resource.AudioID]*audio.Player
	timers      map[resource.AudioID]*util.Timer
	updateFuncs map[resource.AudioID]func(id *audio.Player, targetVol float64, currentVol float64, time float64)
	counter     float64
}

func (s *SoundPlayer) LoadPlayer(playerType string) {
	s.sounds = make(map[resource.AudioID]*audio.Player)
	soundList := []resource.AudioID{
		CardBoard,
		Coins1,
		Crash,
		MoneyCounter,
		Kaching,
		MowerRunning,
		FailedStart,
		PickUpOne,
		PlopSound,
		PouringFood,
		SelectSound,
		SuccessMusic,
		WaterBubbles,
		SelectSound2,
		WhiteBoardMarker1,
		WhiteBoardMarker2,
	}

	musicList := []resource.AudioID{
		BestAdventureEver,
		DayTimeJazz,
		TropicalHouse,
		IndieCafe,
	}

	if playerType == "sound" {
		for _, sound := range soundList {
			s.sounds[sound] = loadedSounds.LoadWAV(sound).Player
			s.sounds[sound].SetVolume(loadedSounds.LoadAudio(sound).Volume)
			bufferDuration := 64 * time.Millisecond
			s.sounds[sound].SetBufferSize(bufferDuration)
			s.updateFuncs = make(map[resource.AudioID]func(id *audio.Player, targetVol float64, currentVol float64, time float64))
			err := s.sounds[sound].Rewind()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if playerType == "music" {
		for _, sound := range musicList {
			s.sounds[sound] = loadedSounds.LoadWAV(sound).Player
			s.sounds[sound].SetVolume(loadedSounds.LoadAudio(sound).Volume)
			s.updateFuncs = make(map[resource.AudioID]func(id *audio.Player, targetVol float64, currentVol float64, time float64))
			err := s.sounds[sound].Rewind()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if playerType == "ogg" {
		for _, sound := range soundList {
			s.sounds[sound] = loadedOGGs.LoadOGG(sound).Player
			s.sounds[sound].SetVolume(loadedOGGs.LoadAudio(sound).Volume)
			s.updateFuncs = make(map[resource.AudioID]func(id *audio.Player, targetVol float64, currentVol float64, time float64))
			err := s.sounds[sound].Rewind()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	if playerType == "oggMusic" {
		for _, sound := range musicList {
			s.sounds[sound] = loadedOGGs.LoadOGG(sound).Player
			s.sounds[sound].SetVolume(loadedOGGs.LoadAudio(sound).Volume)
			s.updateFuncs = make(map[resource.AudioID]func(id *audio.Player, targetVol float64, currentVol float64, time float64))
			err := s.sounds[sound].Rewind()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

func (s *SoundPlayer) Play(id resource.AudioID) {
	s.sounds[id].Play()
}

func (s *SoundPlayer) Pause() {
	for _, sound := range s.sounds {
		sound.Pause()
		err := sound.Rewind()
		if err != nil {
			log.Fatal("Rewinding sound after pausing caused error")
		}
	}
}

func (s *SoundPlayer) FadeIn(id resource.AudioID) {
	err := s.sounds[id].Rewind()
	if err != nil {
		log.Fatal(err)
	}
	s.sounds[id].SetVolume(0.0)
	s.updateFuncs[id] = fade
	s.sounds[id].Play()
}

func (s *SoundPlayer) Update() {
	for key, playing := range s.sounds {
		if s.updateFuncs[key] != nil {
			targetVol := loadedSounds.LoadAudio(key).Volume
			if playing.Volume() >= targetVol-0.01 {
				s.counter = 0
				continue
			}
			s.updateFuncs[key](playing, targetVol, playing.Volume(), s.counter)
			s.counter += 0.016
		}
		if s.updateFuncs[key] != nil && playing.Position() > 1*time.Minute {
			s.updateFuncs[key](playing, 0, playing.Volume(), s.counter)
			s.counter += 0.016
			if playing.Volume() <= 0.015 {
				s.counter = 0
				continue
			}
		}

		if !playing.IsPlaying() && playing.Position() > 0 {
			err := playing.Rewind()
			if err != nil {
				log.Fatal(err)
			}
		}

	}
}

func fade(player *audio.Player, targetVol float64, currenVol float64, time float64) {
	newVol := util.Lerp64(currenVol, targetVol, time*0.001)
	player.SetVolume(newVol)
}
