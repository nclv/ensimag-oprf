package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/circl/oprf"

	"github.com/ensimag-oprf/go/client/core"
)

const serverURL = "http://localhost:1323/api"

var (
	modeFlag *uint
	suiteID  string
	help     bool
)

func commandLine() {
	modeFlag = flag.Uint("mode", uint(oprf.BaseMode), "mode")

	flag.StringVar(&suiteID, "suite", "P256-SHA256", "Cipher suite : P256-SHA256, P384-SHA384 or P521-SHA512")
	flag.BoolVar(&help, "help", false, "Show the usage")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	mode := oprf.Mode(*modeFlag)
	data := flag.Args()

	// Handle required flags
	if help {
		flag.Usage()
		os.Exit(0)
	}

	suite, err := oprf.GetSuite(suiteID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(mode, suite)

	// Set up the client
	client := core.NewClient(serverURL, suite, mode)

	// Convert the string input to bytes
	dataBytes := make([][]byte, len(data))
	for index, input := range data {
		dataBytes[index] = []byte(input)
	}
	// Set default inputs if no data is provided
	if len(dataBytes) == 0 {
		dataBytes = [][]byte{{0x00}, {0xFF}}
	}

	// Request of pseudonymization
	finalizeData, oprfEvaluationRequest, err := client.Blind(dataBytes)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	info := hex.EncodeToString(token)
	log.Println("Public information : ", info)

	blindedElements, err := core.SerializeElements(oprfEvaluationRequest.Elements)
	if err != nil {
		log.Fatal(err)
	}

	evaluationRequest := core.NewEvaluationRequest(
		suite, mode, info, blindedElements,
	)

	evaluationResponse, err := client.EvaluateRequest(evaluationRequest)
	if err != nil {
		log.Fatal(err)
	}

	// Deserialize the public key
	// publicKey, err := core.DeserializePublicKey(suite, evaluationResponse.SerializedPublicKey)
	// if err != nil {
	// 	return
	// }
	// // Set the public key on the client
	// if err := client.SetOPRFClientPublicKey(publicKey); err != nil {
	// 	return
	// }

	for _, element := range evaluationResponse.Evaluation.Elements {
		data, err := element.MarshalBinary()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Evaluation : ", base64.StdEncoding.EncodeToString(data))
	}

	// Finalize the OPRF protocol
	outputs, err := client.Finalize(finalizeData, evaluationResponse.Evaluation, info)
	if err != nil {
		log.Fatal(err)
	}

	for _, output := range outputs {
		log.Println("Output : ", base64.StdEncoding.EncodeToString(output))
	}
}
