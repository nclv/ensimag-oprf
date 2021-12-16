// +build js,wasm

package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"syscall/js"

	"github.com/cloudflare/circl/oprf"
	"github.com/oprf/go/client/core"
	"github.com/oprf/go/common"
)

type WrappedPseudonimizeRequest struct {
	Data []json.RawMessage `json:"data"`
	Mode oprf.Mode `json:"mode"`
	Suite oprf.SuiteID `json:"suite"`
}

type PseudonimizeRequest struct {
	Data [][]byte `json:"data"`
	Mode oprf.Mode `json:"mode"`
	Suite oprf.SuiteID `json:"suite"`
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

func pseudonymize(request *PseudonimizeRequest) ([][]byte, error) {
	// Setup
	client := core.NewClient("http://localhost:1323")
	client.SetupOPRFClient(request.Suite, request.Mode)

	// Request of pseudonymization
	clientRequest := client.CreateRequest(request.Data)

	// The public information
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

		return nil, errors.New("couldn't generate random information")
	}

	info := hex.EncodeToString(token)
	log.Println("Public information : ", info)

	evaluationRequest := common.NewEvaluationRequest(
		request.Suite, request.Mode, info, clientRequest.BlindedElements(),
	)
	evaluation := client.EvaluateRequest(evaluationRequest)

	for _, element := range evaluation.Elements {
		log.Println("Evaluation : ", base64.StdEncoding.EncodeToString(element))
	}

	outputs := client.Finalize(clientRequest, evaluation, info)

	for _, output := range outputs {
		log.Println("Output : ", base64.StdEncoding.EncodeToString(output))
	}

	return outputs, nil
}

func wrapper() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "invalid number of arguments passed"
		}

		jsonInput := args[0].String()
		log.Printf("json input %s\n", jsonInput)

		var pseudonimizeRequest PseudonimizeRequest
		if err := json.Unmarshal([]byte(jsonInput), &pseudonimizeRequest); err != nil {
			log.Println("JSON unmarshalling error :", err, jsonInput)

			return err
		}

		log.Println(pseudonimizeRequest)

		outputs, err := pseudonymize(&pseudonimizeRequest)
		if err != nil {
			log.Println("pseudonymization error :", err)

			return err
		}

		log.Println(outputs)

		encodedOutputs := make([]string, len(outputs))
		for index, output := range outputs {
			encodedOutputs[index] = base64.StdEncoding.EncodeToString(output)
		}

		log.Println(encodedOutputs)

		return encodedOutputs
	})
}

func main() {
	log.Println("Go Web Assembly")
	js.Global().Set("pseudonymize", wrapper())
	select {}
}
