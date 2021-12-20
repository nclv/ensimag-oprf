package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"syscall/js"
)

// wrapper returns a javascript promise that pseudonymize an array of JSON input.
// https://withblue.ink/2020/10/03/go-webassembly-http-requests-and-promises.html
func wrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		jsonInput := args[0].String()
		log.Printf("input : %s\n", jsonInput)

		// Handler for the Promise
		// We need to return a Promise because HTTP requests are blocking in Go
		// All HTTP requests should be wrapped in a goroutine
		// Anonymous because should access jsonInput
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			reject := args[1]

			// Parse the JSON request
			var pseudonimizeRequest PseudonimizeRequest
			if err := json.Unmarshal([]byte(jsonInput), &pseudonimizeRequest); err != nil {
				log.Println("JSON unmarshalling error :", err, jsonInput)
				rejectPromise(reject, err)

				return nil
			}

			// Validate the OPRF mode
			if err := pseudonimizeRequest.ValidateMode(); err != nil {
				log.Println(err)
				rejectPromise(reject, err)

				return nil
			}
			// Validate the encryption suite
			if err := pseudonimizeRequest.ValidateSuite(); err != nil {
				log.Println(err)
				rejectPromise(reject, err)

				return nil
			}

			log.Println(pseudonimizeRequest)

			go func() {
				outputs, err := pseudonymize(&pseudonimizeRequest)
				if err != nil {
					log.Println("pseudonymization error")
					rejectPromise(reject, err)

					return
				}

				log.Println(outputs)

				// Encode the [][]byte outputs to []string
				encodedOutputs := make([]interface{}, len(outputs))
				for index, output := range outputs {
					encodedOutputs[index] = base64.StdEncoding.EncodeToString(output)
				}

				log.Println(encodedOutputs)

				// map[string]interface{} is parsed by js.ValueOf and put into a javascript Object
				data := map[string]interface{}{"pseudonymized_data": encodedOutputs}
				objectConstructor := js.Global().Get("Object")
				dataJS := objectConstructor.New(data)

				// Resolve the Promise by sending the object
				resolve.Invoke(dataJS)
			}()

			// The handler of a Promise doesn't return any value
			return nil
		})

		// Create and return the Promise object
		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
