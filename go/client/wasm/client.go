package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/oprf/go/client/core"
	"github.com/oprf/go/common"
)

// pseudonymize execute the PseudonimizeRequest and returns the pseudonymized data bytes.
// It creates the client request, generate the random public information, call the server for
// an evaluation and finalize the protocol.
func pseudonymize(request *PseudonimizeRequest) ([][]byte, error) {
	// Set up the client with the mode and suite
	client := core.NewClient(serverURL)
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
