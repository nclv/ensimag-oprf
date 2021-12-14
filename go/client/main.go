package main

import (
	"encoding/base64"
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

	info := "7465737420696e666f"
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
