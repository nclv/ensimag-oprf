//go:build js && wasm
// +build js,wasm

package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"syscall/js"

	"github.com/cloudflare/circl/oprf"
	"github.com/oprf/go/client/core"
	"github.com/oprf/go/common"
)

type WrappedPseudonimizeRequest struct {
	Data  []json.RawMessage `json:"data"`
	Mode  oprf.Mode         `json:"mode"`
	Suite oprf.SuiteID      `json:"suite"`
}

type PseudonimizeRequest struct {
	Data  [][]byte     `json:"data"`
	Mode  oprf.Mode    `json:"mode"`
	Suite oprf.SuiteID `json:"suite"`
}

func (p *PseudonimizeRequest) Validate() (bool, error) {
	switch p.Mode {
	case oprf.BaseMode, oprf.VerifiableMode:

	}

	return true, nil
}

func (p *PseudonimizeRequest) UnmarshalJSON(data []byte) error {
	var wp WrappedPseudonimizeRequest
	if err := json.Unmarshal(data, &wp); err != nil {
		return err
	}

	p.Data = make([][]byte, len(wp.Data))

	var input string
	for index, stringInput := range wp.Data {
		if err := json.Unmarshal(stringInput, &input); err != nil {
			return err
		}
		p.Data[index] = []byte(input)
	}

	p.Suite = wp.Suite
	p.Mode = wp.Mode

	return nil
}

// pseudonymize execute the PseudonimizeRequest and returns the pseudonymized data bytes.
// It creates the client request, generate the random public information, call the server for
// an evaluation and finalize the protocol.
func pseudonymize(request *PseudonimizeRequest) ([][]byte, error) {
	// Set up the client with the mode and suite
	client := core.NewClient("http://localhost:1323")
	if err := client.SetupOPRFClient(request.Suite, request.Mode); err != nil {
		return [][]byte{}, fmt.Errorf("couldn't setup OPRF client : %w", err)
	}

	// Request of pseudonymization
	clientRequest := client.CreateRequest(request.Data)

	// The public information (client SECRET)
	// It is RECOMMENDED that this metadata be constructed with some type of higher-level
	// domain separation to avoid cross protocol attacks or related issues.
	// For example, protocols using this construction might ensure that the metadata uses
	// a unique, prefix-free encoding.
	// Any system which has multiple POPRF applications should distinguish client inputs to
	// ensure the POPRF results are separate.
	// info := "7465737420696e666f"
	// Generate a random information for each request for non-deterministic results
	token := make([]byte, 256)
	if _, err := rand.Read(token); err != nil {
		log.Println(err)

		return nil, fmt.Errorf("couldn't generate random information : %w", err)
	}

	info := hex.EncodeToString(token)
	// DO NOT SHARE THE PUBLIC INFORMATION
	// log.Println("Public information : ", info)

	evaluationRequest := common.NewEvaluationRequest(
		request.Suite, request.Mode, info, clientRequest.BlindedElements(),
	)
	evaluation := client.EvaluateRequest(evaluationRequest)

	//for _, element := range evaluation.Elements {
	//	log.Println("Evaluation : ", base64.StdEncoding.EncodeToString(element))
	//}

	outputs := client.Finalize(clientRequest, evaluation, info)

	//for _, output := range outputs {
	//	log.Println("Output : ", base64.StdEncoding.EncodeToString(output))
	//}

	return outputs, nil
}

func wrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		jsonInput := args[0].String()
		log.Printf("input : %s\n", jsonInput)

		// Handler for the Promise
		// We need to return a Promise because HTTP requests are blocking in Go
		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			resolve := args[0]
			reject := args[1]

			var pseudonimizeRequest PseudonimizeRequest
			if err := json.Unmarshal([]byte(jsonInput), &pseudonimizeRequest); err != nil {
				log.Println("JSON unmarshalling error :", err, jsonInput)

				// Handle errors: reject the Promise if we have an error
				errorConstructor := js.Global().Get("Error")
				errorObject := errorConstructor.New(err.Error())
				reject.Invoke(errorObject)

				return nil
			}

			log.Println(pseudonimizeRequest)

			// Run this code asynchronously
			go func() {

				outputs, err := pseudonymize(&pseudonimizeRequest)
				if err != nil {
					log.Println("pseudonymization error")
					// Handle errors: reject the Promise if we have an error
					errorConstructor := js.Global().Get("Error")
					errorObject := errorConstructor.New(err.Error())
					reject.Invoke(errorObject)

					return
				}

				log.Println(outputs)

				encodedOutputs := make([]interface{}, len(outputs))
				for index, output := range outputs {
					encodedOutputs[index] = base64.StdEncoding.EncodeToString(output)
				}

				log.Println(encodedOutputs)

				data := map[string]interface{}{"encoded outputs": encodedOutputs}
				objectConstructor := js.Global().Get("Object")
				dataJS := objectConstructor.New(data)

				// Resolve the Promise
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

func main() {
	log.Println("Go Web Assembly")
	// https://withblue.ink/2020/10/03/go-webassembly-http-requests-and-promises.html

	js.Global().Set("pseudonymize", wrapper())

	select {}
}
