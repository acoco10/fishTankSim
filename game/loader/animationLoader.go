package loader

import (
	"encoding/json"
	"github.com/acoco10/QuickDrawAdventure/animations"
	"github.com/acoco10/QuickDrawAdventure/spriteSheet"
	"github.com/acoco10/fishTankWebGame/assets"
	"log"
)

type AsepriteFrame struct {
	Frame struct {
		X int `json:"x"`
		Y int `json:"y"`
		W int `json:"w"`
		H int `json:"h"`
	} `json:"frame"`
	Duration int `json:"duration"` // ms
}

type AsepriteData struct {
	Frames []AsepriteFrame `json:"frames"`
}

func LoadAnimation(path string) (*animations.Animation, *spritesheet.SpriteSheet, error) {

	var data AsepriteData

	fileData, err := assets.AnimationDataDir.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	err = json.Unmarshal(fileData, &data)
	if err != nil {
		return nil, nil, err
	}

	frameCount := len(data.Frames)
	if frameCount == 0 {
		panic("no frames found")
	}

	// Just pick first duration

	log.Printf("Width result from animation loader:%d", data.Frames[0].Frame.W)
	log.Printf("Height result from animation loader:%d", data.Frames[0].Frame.H)
	log.Printf("Duration result from animation loader:%d", data.Frames[0].Duration)

	sps := spritesheet.NewSpritesheet(len(data.Frames), 1, data.Frames[0].Frame.W, data.Frames[0].Frame.H)
	ani := animations.NewAnimation(0, len(data.Frames)-1, 1, float32(1000/data.Frames[0].Duration))

	return ani, sps, nil
}
