package main

import (
	"crypto/rand"
	"log"

	"github.com/oprf/go/utils"

	"github.com/cloudflare/circl/oprf"
)

func GetServer(suite oprf.SuiteID, mode oprf.Mode, privateKey *oprf.PrivateKey) *oprf.Server {
	var (
		server *oprf.Server
		err    error
	)

	if mode == oprf.BaseMode {
		server, err = oprf.NewServer(suite, privateKey)
	} else if mode == oprf.VerifiableMode {
		server, err = oprf.NewVerifiableServer(suite, privateKey)
	}

	if err != nil {
		log.Println(err)
	}

	return server
}

func GetClient(suite oprf.SuiteID, mode oprf.Mode, pkS *oprf.PublicKey) *oprf.Client {
	var (
		client *oprf.Client
		err    error
	)

	if mode == oprf.BaseMode {
		client, err = oprf.NewClient(suite)
	} else if mode == oprf.VerifiableMode {
		client, err = oprf.NewVerifiableClient(suite, pkS)
	}

	if err != nil {
		log.Println(err)
	}

	return client
}

// ServerSideOPRF performs a full OPRF protocol at server-side.
// There is no need for a client.
func ServerSideOPRF(input, info []byte, server *oprf.Server) {
	output, err := server.FullEvaluate(input, info)
	if err != nil {
		log.Println(err)
	}

	utils.PrintByteArray(output)
}

// ClientServerOPRF performs a full OPRF protocol with a client and a server.
func ClientServerOPRF(inputs [][]byte, info []byte, client *oprf.Client, server *oprf.Server) {
	// Request generates a request for server passing an array of inputs to be evaluated by server.
	clientRequest, err := client.Request(inputs)
	if err != nil {
		log.Println(err)
	}

	log.Println(clientRequest.BlindedElements())
	// Evaluate evaluates a set of blinded inputs from the client.
	// BlindedElements returns the serialized blinded elements produced for the client request.
	// TODO: Send the blinded inputs to the server.
	evaluation, err := server.Evaluate(clientRequest.BlindedElements(), info)
	if err != nil || evaluation == nil {
		log.Println(err)
	}

	// Finalize computes the signed token from the server Evaluation and returns the output of the
	// OPRF protocol. The function uses server's public key to verify the proof in verifiable mode.
	// TODO: send the evaluation to the client
	clientOutputs, err := client.Finalize(clientRequest, evaluation, info)
	if err != nil || clientOutputs == nil {
		log.Println(err)
	}

	for index := range inputs {
		input := inputs[index]
		output := clientOutputs[index]

		log.Println("Input :")
		utils.PrintByteArray(input)
		log.Println("Output :")
		utils.PrintByteArray(output)

		// VerifyFinalize performs a full OPRF protocol (ie. calls FullEvaluate) and returns true if the
		// output matches the expected output.
		valid := server.VerifyFinalize(input, info, output)
		log.Println("Verification from the server is valid", valid)
	}
}

func main() {
	suite := oprf.OPRFP256
	mode := oprf.BaseMode

	// INITIALIZATION (on server launch)
	// Generate the private key(s) on the server (or recover it(them) from previous server)
	privateKey, err := oprf.GenerateKey(suite, rand.Reader)
	if err != nil {
		log.Println(err)
	}

	utils.PrintKey(privateKey)
	utils.PrintKey(privateKey.Public())

	// On the server
	// TODO: create #{suite} * #{mode} servers with different private keys? or with same private key?
	server := GetServer(suite, mode, privateKey)
	publicKey := privateKey.Public()
	// For one server :
	// TODO: public information generation (https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-voprf#section-5.2)
	// TODO: send the public key (only for a verifiable protocol ie. 1 or (#{suite} * #{mode}) / 2 public keys)
	// TODO: send the public information to the client (use always the same information like in the draft test vectors?, 7465737420696e666f)

	// On the client
	// TODO: receive the public keys and the public information
	// END OF INITIALIZATION
	// TODO: choose the suite and the mode, send it to the server that choose the right instance
	client := GetClient(suite, mode, publicKey)

	ServerSideOPRF([]byte("hello world"), []byte("test info"), server)
	// In the base mode, a client and server interact to
	// compute output = F(skS, input, info), where input is the client's
	// private input, skS is the server's private key, info is the public
	// input (or metadata), and output is the POPRF output.  The client
	// learns output and the server learns nothing.  In the verifiable mode,
	// the client also receives proof that the server used skS in computing
	// the function.
	ClientServerOPRF([][]byte{{0x00}, {0xFF}}, []byte("7465737420696e666f"), client, server)
}
