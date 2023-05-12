package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/cloudflare/circl/oprf"
	"github.com/ensimag-oprf/go/client/core"
)

// PseudonymizeResponse contains the pseudonymized data output and the public information if requested.
type PseudonymizeResponse struct {
	Info    string
	Outputs [][]byte
}

// pseudonymize execute the PseudonimizeRequest and returns the PseudonymizedResponse.
// It creates the client request, generate the random static information, call the server for
// an evaluation and finalize the protocol.
func pseudonymize(request *PseudonimizeRequest) (*PseudonymizeResponse, error) {
	suite, err := oprf.GetSuite(request.Suite)
	if err != nil {
		return nil, fmt.Errorf("the suite is not suppported : %w", err)
	}

	// Set up the client with the mode and suite
	client := core.NewClient(serverURL, suite, request.Mode)

	// Request of pseudonymization
	finalizeData, oprfEvaluationRequest, err := client.Blind(request.Data)
	if err != nil {
		return nil, fmt.Errorf("couldn't create the request : %w", err)
	}

	// The static information (client SECRET)
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

	blindedElements, err := core.SerializeElements(oprfEvaluationRequest.Elements)
	if err != nil {
		log.Fatal(err)
	}

	// Evaluate the request
	evaluationRequest := core.NewEvaluationRequest(
		suite, request.Mode, info, blindedElements,
	)

	evaluationResponse, err := client.EvaluateRequest(evaluationRequest)
	if err != nil {
		return nil, fmt.Errorf("couldn't evaluate the request : %w", err)
	}

	// Deserialize the public key
	// publicKey, err := core.DeserializePublicKey(request.Suite, evaluationResponse.SerializedPublicKey)
	// if err != nil {
	// 	return nil, fmt.Errorf("coudn't deserialize the public key : %w", err)
	// }
	// // Set the public key on the client
	// if err := client.SetOPRFClientPublicKey(publicKey); err != nil {
	// 	return nil, fmt.Errorf("coudn't update the client's public key : %w", err)
	// }

	// Finalize the protocol
	outputs, err := client.Finalize(finalizeData, evaluationResponse.Evaluation, info)
	if err != nil {
		return nil, fmt.Errorf("couldn't finalize the request : %w", err)
	}

	// Build the response
	response := &PseudonymizeResponse{Outputs: outputs}
	if request.ReturnInfo {
		response.Info = info
	}

	return response, nil
}
