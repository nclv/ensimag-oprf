package main

import "syscall/js"

// rejectPromise handle errors: reject the Promise if we have an error
func rejectPromise(reject js.Value, err error) {
	errorConstructor := js.Global().Get("Error")
	errorObject := errorConstructor.New(err.Error())
	reject.Invoke(errorObject)
}
