package entities

import (
	"encoding/json"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/acoco10/fishTankWebGame/game/events"
	"log"
)

type SavedFish struct {
	Name      string  `json:"Name"`
	Size      int     `json:"Size"`
	Progress  float32 `json:"Progress"`
	NextLevel float32 `json:"NextLevel"`
	FishType  string  `json:"FishType"`
	MaxSpeed  float32 `json:"MaxSpeed"`
}

type SaveGameState struct {
	Fish  []SavedFish    `json:"State"`
	Tasks []*events.Task `json:"tasks"`
}

func (gs *SaveGameState) ToJSON() string {
	b, _ := json.Marshal(gs)
	return string(b)
}

func (gs *SaveGameState) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), gs)
}

func LoadSaveJson(fileName string) string {
	contents, err := assets.DataDir.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	var saveState = SaveGameState{}

	err = json.Unmarshal(contents, &saveState)

	if err != nil {
		log.Fatal(err)
	}

	return saveState.ToJSON()
}
