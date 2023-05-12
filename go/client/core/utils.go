package core

import (
	"log"

	"github.com/cloudflare/circl/group"
	"github.com/cloudflare/circl/oprf"
)

// DeserializePublicKey deserialize a public key
func DeserializePublicKey(suite oprf.Suite, serializedPublicKey []byte) (*oprf.PublicKey, error) {
	publicKey := new(oprf.PublicKey)
	if err := publicKey.UnmarshalBinary(suite, serializedPublicKey); err != nil {
		log.Println(err)

		return nil, err
	}

	return publicKey, nil
}

// DeserializePublicKeys deserialize the server's public keys
func DeserializePublicKeys(serializedPublicKeys map[string][]byte) (map[string]*oprf.PublicKey, error) {
	publicKeys := make(map[string]*oprf.PublicKey)

	for suiteID, serializedPublicKey := range serializedPublicKeys {
		// log.Println(
		// 	"Suite ID:", suiteID,
		//	", Public key : ", base64.StdEncoding.EncodeToString(serializedPublicKey),
		// )

		suite, err := oprf.GetSuite(suiteID)
		if err != nil {
			return nil, err
		}

		publicKey, err := DeserializePublicKey(suite, serializedPublicKey)
		if err != nil {
			return nil, err
		}

		publicKeys[suiteID] = publicKey
	}

	return publicKeys, nil
}

func SerializeElements(elements []group.Element) ([][]byte, error) {
	blindedElements := make([][]byte, len(elements))

	for index, element := range elements {
		blindedElement, err := element.MarshalBinary()
		if err != nil {
			log.Println(err)

			return nil, err
		}

		blindedElements[index] = blindedElement
	}

	return blindedElements, nil
}
