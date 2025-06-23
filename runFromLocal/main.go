package main

import (
	"encoding/json"
	"fmt"
	"github.com/acoco10/fishTankWebGame/assets"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/scenes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/pkg/profile"
	"log"
)

type GameState struct {
	Username string        `json:"username"`
	State    []interface{} `json:"state"`
}

func main() {

	stateData, err := assets.DataDir.ReadFile("data/saveWithTasks.json")
	if err != nil {
		fmt.Errorf("cant read test save file from embed dir %t", err)
	}

	var data GameState
	json.Unmarshal(stateData, &data)

	var state entities.SaveGameState

	b, err := json.Marshal(data.State)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &state.Fish)
	if err != nil {
		log.Fatal(err)
	}

	gameLog := sceneManagement.NewGameLog(state)
	p := profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook)
	defer p.Stop()
	g := scenes.NewGame(gameLog, scenes.NewUser)
	err = ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}
