package main

import (
	"log"

	"github.com/cloudflare/circl/oprf"
)

// NewOPRFClient create an OPRF client with the provided suite, mode and public key.
// No public key is needed for oprf.BaseMode.
func NewOPRFClient(suite oprf.SuiteID, mode oprf.Mode, pkS *oprf.PublicKey) *oprf.Client {
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
