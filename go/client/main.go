package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"log"

	"github.com/cloudflare/circl/oprf"
)

var (
	mode  oprf.Mode
	suite oprf.SuiteID
)

func init() {
	modeFlag := flag.Int("mode", int(oprf.BaseMode), "mode")
	suiteFlag := flag.Int("suite", int(oprf.OPRFP256), "cipher suite")

	flag.Parse()

	mode = oprf.Mode(*modeFlag)
	suite = oprf.SuiteID(*suiteFlag)
}

func main() {
	log.Println(mode, suite)

	// Setup
	client := NewClient("http://localhost:1323")
	client.SetupOPRFClient(suite, mode)

	// Request of pseudonymization
	clientRequest := client.CreateRequest([][]byte{{0x00}, []byte("Hello")})

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

		return
	}

	info := hex.EncodeToString(token)
	log.Println("Public information : ", info)

	evaluationRequest := &EvaluationRequest{
		Suite:           suite,
		Mode:            mode,
		Info:            info,
		BlindedElements: clientRequest.BlindedElements(),
	}
	evaluation := client.EvaluateRequest(evaluationRequest)

	for _, element := range evaluation.Elements {
		log.Println("Evaluation : ", base64.StdEncoding.EncodeToString(element))
	}

	outputs := client.Finalize(clientRequest, evaluation, info)

	for _, output := range outputs {
		log.Println("Output : ", base64.StdEncoding.EncodeToString(output))
	}
}
