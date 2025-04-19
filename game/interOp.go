//go:build js && wasm
// +build js,wasm

package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"syscall/js"
)

var gameStarted = false

func SaveToBackend(data string) {
	js.Global().Call("saveGame", data)
}

var loadedGameState string

func receiveGameState(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		println("No game state passed!")
		return nil
	}

	loadedGameState = args[0].String()

	println("Loaded game state from JS:", loadedGameState)

	// You can now use `loadedGameState` inside your game init logic
	if !gameStarted {
		gameStarted = true
		go func() {
			// Initialize your game struct using loadedGameState
			game := NewGame(loadedGameState) // Replace with your actual function
			if err := ebiten.RunGame(game); err != nil {
				println("Failed to start game:", err.Error())
			}
		}()
	}

	return nil
}

func WasmStartUp() {
	js.Global().Set("loadGameState", js.FuncOf(receiveGameState))
	js.Global().Set("wasmReady", js.ValueOf(true))
	// Proceed to start your Ebiten game (defer until game state is set if needed)

}
