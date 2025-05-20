package main

import (
	"encoding/json"
	"github.com/acoco10/fishTankWebGame/game"
	"github.com/acoco10/fishTankWebGame/game/gameEntities"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/scenes"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"syscall/js"
)

type GameState struct {
	Username string        `json:"username"`
	State    []interface{} `json:"state"`
}

func main() {
	loadFunc := js.Global().Get("loadSaveData")
	promise := loadFunc.Invoke()
	// Wait for it
	result, err := game.AwaitPromise(promise)
	if err != nil {
		log.Fatal(err)
	}

	bytes := []byte(js.Global().Get("JSON").Call("stringify", result).String())

	println("go result from promise:", bytes)
	var state entities.SaveGameState

	err = json.Unmarshal(bytes, &state)
	if err != nil {
		log.Fatal("error unmarshalling bytes into save game state struct:", err)
	}

	for n, fish := range state.Fish {
		println("fish", n, fish.Size)
	}

	gameLog := sceneManagement.NewGameLog(state)

	var g ebiten.Game

	if len(gameLog.Save.Fish) > 0 {
		g = scenes.NewGame(gameLog, scenes.ExistingUser)
	} else {
		g = scenes.NewGame(gameLog, scenes.NewUser)
	}

	err = ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}

}
