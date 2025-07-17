package main

import (
	"encoding/json"
	"github.com/acoco10/fishTankWebGame/game/entities"
	"github.com/acoco10/fishTankWebGame/game/loader"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/scenes"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

type GameState struct {
	Username string        `json:"username"`
	State    []interface{} `json:"state"`
}

func main() {

	var data GameState

	/*stateData, err := assets.DataDir.ReadFile("data/saveWithTasks.json")
	if err != nil {
		fmt.Errorf("cant read test save file from embed dir %t", err)
	}


	err = json.Unmarshal(stateData, &data)
	if err != nil {
		log.Fatal(err)
	}*/

	loader.LoadShaderRegistry()
	loader.LoadFontRegistry()
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
	//p := profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook)
	//defer p.Stop()
	g := scenes.NewGame(gameLog, scenes.NewUser)
	err = ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}
