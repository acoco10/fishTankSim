package main

import (
	"encoding/json"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/registry"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/scenes"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	_ "net/http/pprof"
)

type GameState struct {
	Username string        `json:"username"`
	State    []interface{} `json:"state"`
}

func main() {

	/*	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()*/

	// Your game code here

	var data GameState

	/*stateData, err := assets.DataDir.ReadFile("data/saveWithTasks.json")
	if err != nil {
		fmt.Errorf("cant read test save file from embed dir %t", err)
	}


	err = json.Unmarshal(stateData, &data)
	if err != nil {
		log.Fatal(err)
	}*/

	registry.LoadFontRegistry()
	var state entities.SaveGameState

	b, err := json.Marshal(data.State)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &state.Fish)
	if err != nil {
		log.Fatal(err)
	}

	gameLog := sceneManagement.NewGameLog(state, "")

	g := scenes.NewGame(gameLog, scenes.NewUser)
	err = ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}
