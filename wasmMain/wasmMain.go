package main

import (
	"encoding/json"
	"fishTankWebGame/game"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"syscall/js"
)

func main() {
	loadFunc := js.Global().Get("loadSaveData")
	promise := loadFunc.Invoke()

	// Wait for it
	result, err := game.AwaitPromise(promise)
	if err != nil {
		fmt.Println("Error loading data:", err)
		return
	}

	bytes := []byte(js.Global().Get("JSON").Call("stringify", result).String())
	var data map[string]interface{}
	json.Unmarshal(bytes, &data)

	fmt.Println("Loaded in Go:", data)

	g := game.NewGame()
	err = ebiten.RunGame(g)
	if err != nil {
		log.Fatal(err)
	}
}
