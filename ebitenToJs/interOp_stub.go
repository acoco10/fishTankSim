//go:build !js || !wasm
// +build !js !wasm

package ebitenToJs

func SaveToBackend(data string) {
	// No-op: Only used in WASM build
}

func WasmStartUp() {
	// No-op: Only used in WASM build
}
