//go:build js && wasm
// +build js,wasm

package main

import (
	"log"
	"syscall/js"
)

// WASM sources
// https://golangbot.com/webassembly-using-go/
// https://ian-says.com/articles/golang-in-the-browser-with-web-assembly/
// https://github.com/tinygo-org/tinygo/tree/master/src/examples/wasm
// https://about.sourcegraph.com/go/gophercon-2019-get-going-with-webassembly/

const serverURL = "/api"

func main() {
	log.Println("Go Web Assembly")

	js.Global().Set("pseudonymize", wrapper())

	select {}
}
