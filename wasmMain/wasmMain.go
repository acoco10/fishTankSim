package main

import (
	"encoding/json"
	"fishTankWebGame/game"
	"fishTankWebGame/game/gameEntities"
	"fmt"
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
		fmt.Println("Error loading data:", err)
		return
	}
	println("go result from promise:", result.String())

	bytes := []byte(js.Global().Get("JSON").Call("stringify", result).String())
	var data GameState
	json.Unmarshal(bytes, &data)
	var state gameEntities.SaveGameState

	b, err := json.Marshal(data.State)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, &state.Fish)
	if err != nil {
		log.Fatal(err)
	}

	for n, fish := range state.Fish {
		println("fish", n, fish.Size)
	}

	g := game.NewGame(state)
	err = ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}
