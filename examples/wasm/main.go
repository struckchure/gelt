//go:build wasm
// +build wasm

package main

import (
	"fmt"
	"syscall/js"
)

// GOOS=js GOARCH=wasm go build -o main.wasm

func main() {
	fmt.Println("WASM Loaded!")

	js.Global().Call("alert", "Hello from Go WebAssembly!")
}
