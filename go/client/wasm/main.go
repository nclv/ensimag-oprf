//go:build js && wasm

package main

import (
	"log"
	"syscall/js"
)

// TODO: change the server URL
const serverURL = "http://localhost:1323"

func main() {
	log.Println("Go Web Assembly")

	js.Global().Set("pseudonymize", wrapper())

	select {}
}
