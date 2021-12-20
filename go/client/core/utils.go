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

// DeserializePublicKey deserialize a public key
func DeserializePublicKey(suiteID oprf.SuiteID, serializedPublicKey []byte) (*oprf.PublicKey, error) {
	publicKey := new(oprf.PublicKey)
	if err := publicKey.Deserialize(suiteID, serializedPublicKey); err != nil {
		log.Println(err)

		return nil, err
	}

	return publicKey, nil
}

// DeserializePublicKeys deserialize the server's public keys
func DeserializePublicKeys(serializedPublicKeys map[oprf.SuiteID][]byte) (map[oprf.SuiteID]*oprf.PublicKey, error) {
	publicKeys := make(map[oprf.SuiteID]*oprf.PublicKey)

	for suiteID, serializedPublicKey := range serializedPublicKeys {
		// log.Println(
		// 	"Suite ID:", suiteID,
		//	", Public key : ", base64.StdEncoding.EncodeToString(serializedPublicKey),
		// )

		publicKey, err := DeserializePublicKey(suiteID, serializedPublicKey)
		if err != nil {
			return nil, err
		}

		publicKeys[suiteID] = publicKey
	}

	return publicKeys, nil
}
