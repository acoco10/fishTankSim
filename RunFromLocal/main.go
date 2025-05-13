package main

import (
	"encoding/json"
	"fmt"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/acoco10/fishTankWebGame/game"
	"github.com/acoco10/fishTankWebGame/game/gameEntities"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/soundFX"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

type GameState struct {
	Username string        `json:"username"`
	State    []interface{} `json:"state"`
}

func main() {

	stateData, err := assets.DataDir.ReadFile("data/testSave.json")
	if err != nil {
		fmt.Errorf("cant read test save file from embed dir %t", err)
	}

	var data GameState
	json.Unmarshal(stateData, &data)
	var state gameEntities.SaveGameState

	b, err := json.Marshal(data.State)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &state.Fish)
	if err != nil {
		log.Fatal(err)
	}

	gameLog := sceneManagement.GameLog{Save: &gameEntities.SaveGameState{Fish: []gameEntities.SavedFish{}}}
	gameLog.Save = &state
	eHub := gameEntities.NewEventHub()
	gameLog.GlobalEventHub = eHub

	songP, err := soundFX.NewSongPlayer()
	if err != nil {
		log.Fatal(err)
	}

	soundP, err := soundFX.NewSongPlayer()
	if err != nil {
		log.Fatal(err)
	}

	gameLog.SongPlayer = songP
	gameLog.SoundPlayer = soundP

	g := game.NewGame(&gameLog, game.ExistingUser)
	err = ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}

}
