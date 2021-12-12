package main

import (
	"crypto/rand"
	"github.com/cloudflare/circl/oprf"
	"log"
)

// generatePrivateKey generate a private key for the encryption suite
func generatePrivateKey(suite oprf.SuiteID) *oprf.PrivateKey {
	privateKey, err := oprf.GenerateKey(suite, rand.Reader)
	if err != nil {
		log.Println(err)
	}

	// printKey(privateKey)
	// printKey(privateKey.Public())

	return privateKey
}

// serializePublicKey is a wrapper to serialize a public key
func serializePublicKey(key *oprf.PrivateKey) []byte {
	publicKey, err := key.Public().Serialize()
	if err != nil {
		log.Println(err)
	}

	return publicKey
}

// createServer create a server with the provided suite, mode and private key
func createServer(suite oprf.SuiteID, mode oprf.Mode, privateKey *oprf.PrivateKey) *oprf.Server {
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
