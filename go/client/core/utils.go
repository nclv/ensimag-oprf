package core

import (
	"log"

	"github.com/cloudflare/circl/oprf"
)

// NewOPRFClient create an OPRF client with the provided suite, mode and static key.
// No static key is needed for oprf.BaseMode.
func NewOPRFClient(suite oprf.SuiteID, mode oprf.Mode, pkS *oprf.PublicKey) (*oprf.Client, error) {
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
		log.Println("error when creating OPRF client", err)

		return nil, err
	}

	return client, nil
}

// DeserializePublicKeys deserialize the server's static keys
func DeserializePublicKeys(serializedPublicKeys map[oprf.SuiteID][]byte) map[oprf.SuiteID]*oprf.PublicKey {
	publicKeys := make(map[oprf.SuiteID]*oprf.PublicKey)

	for suiteID, serializedPublicKey := range serializedPublicKeys {
		// log.Println(
		// 	"Suite ID:", suiteID,
		//	", Public key : ", base64.StdEncoding.EncodeToString(serializedPublicKey),
		// )

		publicKey := new(oprf.PublicKey)
		if err := publicKey.Deserialize(suiteID, serializedPublicKey); err != nil {
			log.Println(err)
		}

		publicKeys[suiteID] = publicKey
	}

	return publicKeys
}
