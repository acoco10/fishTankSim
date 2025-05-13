package main

import (
	"encoding/json"
	"fmt"
	"github.com/acoco10/fishTankWebGame/game"
	"github.com/acoco10/fishTankWebGame/game/gameEntities"
	"github.com/acoco10/fishTankWebGame/game/sceneManagement"
	"github.com/acoco10/fishTankWebGame/game/soundFX"

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

	gameLog := sceneManagement.GameLog{}
	gameLog.Save = &state

	eHub := gameEntities.NewEventHub()
	gameLog.GlobalEventHub = eHub

	songP, err := soundFX.NewSongPlayer()
	if err != nil {
		log.Fatal(err)
	}

	gameLog.SongPlayer = songP

	soundP, err := soundFX.NewSongPlayer()
	if err != nil {
		log.Fatal(err)
	}

	gameLog.SoundPlayer = soundP

	var g ebiten.Game

	if len(gameLog.Save.Fish) > 0 {
		g = game.NewGame(&gameLog, game.ExistingUser)
	} else {
		g = game.NewGame(&gameLog, game.NewUser)
	}
	err = ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}

}
