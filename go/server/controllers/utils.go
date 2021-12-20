package controllers

import (
	"crypto/rand"
	"log"

	"github.com/cloudflare/circl/oprf"
)

// GeneratePrivateKey generate a private key for the encryption suite
func GeneratePrivateKey(suite oprf.SuiteID) *oprf.PrivateKey {
	privateKey, err := oprf.GenerateKey(suite, rand.Reader)
	if err != nil {
		log.Println(err)
	}

	return privateKey
}

// SerializePublicKey is a wrapper to serialize a public key
func SerializePublicKey(key *oprf.PrivateKey) []byte {
	publicKey, err := key.Public().Serialize()
	if err != nil {
		log.Println(err)
	}

	return publicKey
}

// NewOPRFServer create an OPRF server with the provided suite, mode and private key
func NewOPRFServer(suite oprf.SuiteID, mode oprf.Mode, privateKey *oprf.PrivateKey) *oprf.Server {
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
