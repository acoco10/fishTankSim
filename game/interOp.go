//go:build js && wasm
// +build js,wasm

package game

import (
	"fmt"
	"syscall/js"
)

var gameStarted = false

func SaveToBackend(data []byte) {
	js.Global().Call("saveGame", string(data))
}

func LoadSaveDataFromJS(this js.Value, args []js.Value) interface{} {

	if len(args) < 1 {
		fmt.Println("No data received")
		return nil
	}

	state := args[0].Get("state").String()
	fmt.Println("WASM received state:", state)
	return nil

}

func AwaitPromise(p js.Value) (js.Value, error) {
	ch := make(chan struct{})
	var result js.Value
	var err error

	thenFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		result = args[0]
		close(ch)
		return nil
	})

	catchFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		err = fmt.Errorf("promise rejected: %v", args[0])
		close(ch)
		return nil
	})

	p.Call("then", thenFunc).Call("catch", catchFunc)
	<-ch
	return result, err
}
