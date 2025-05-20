//go:build !js || !wasm
// +build !js !wasm

package game

func SaveToBackend(data []byte) {
	// No-op: Only used in WASM build
}

func WasmStartUp() {
	// No-op: Only used in WASM build
}
