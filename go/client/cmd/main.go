package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"log"

	"github.com/cloudflare/circl/oprf"

	"github.com/ensimag-oprf/go/client/core"
)

const serverURL = "http://localhost:1323/api"

var (
	data  []string
	mode  oprf.Mode
	suite oprf.SuiteID
)

func commandLine() {
	modeFlag := flag.Int("mode", int(oprf.BaseMode), "mode")
	suiteFlag := flag.Int("suite", int(oprf.OPRFP256), "cipher suite")

	flag.Parse()
	data = flag.Args()

	mode = oprf.Mode(*modeFlag)
	suite = oprf.SuiteID(*suiteFlag)
}

func main() {
	commandLine()

	log.Println(mode, suite)

	// Set up the OPRF client
	client := core.NewClient(serverURL)
	if err := client.SetupOPRFClient(suite, mode); err != nil {
		return
	}

	// Convert the string input to bytes
	dataBytes := make([][]byte, len(data))
	for index, input := range data {
		dataBytes[index] = []byte(input)
	}
	// Set default inputs if no data is provided
	if len(dataBytes) == 0 {
		dataBytes = append(dataBytes, [][]byte{{0x00}, {0xFF}}...)
	}

	// Request of pseudonymization
	clientRequest, err := client.CreateRequest(dataBytes)
	if err != nil {
		return
	}

	// The static information
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

	evaluationRequest := core.NewEvaluationRequest(
		suite, mode, info, clientRequest.BlindedElements(),
	)
	evaluationResponse, err := client.EvaluateRequest(evaluationRequest)
	if err != nil {
		return
	}

	// Deserialize the public key
	publicKey, err := core.DeserializePublicKey(suite, evaluationResponse.SerializedPublicKey)
	if err != nil {
		return
	}
	// Set the public key on the client
	if err := client.SetOPRFClientPublicKey(publicKey); err != nil {
		return
	}

	for _, element := range evaluationResponse.Evaluation.Elements {
		log.Println("Evaluation : ", base64.StdEncoding.EncodeToString(element))
	}

	// Finalize the OPRF protocol
	outputs, err := client.Finalize(clientRequest, evaluationResponse.Evaluation, info)
	if err != nil {
		return
	}

	for _, output := range outputs {
		log.Println("Output : ", base64.StdEncoding.EncodeToString(output))
	}
}
