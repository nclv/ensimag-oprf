package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/cloudflare/circl/group"
	"github.com/cloudflare/circl/oprf"
)

// LoadPrivateKey decode the base64 serialized private key and deserialized it.
func LoadPrivateKey(suite oprf.Suite, serializedBase64Key string) (*oprf.PrivateKey, error) {
	serializedKey, err := base64.StdEncoding.DecodeString(serializedBase64Key)
	if err != nil {
		log.Println(err)

		return nil, fmt.Errorf("couldn't load the private key : %w", err)
	}

	privateKey := new(oprf.PrivateKey)
	if err := privateKey.UnmarshalBinary(suite, serializedKey); err != nil {
		log.Println(err)

		return nil, err
	}

	return privateKey, nil
}

// LoadOrGenerateKey tries to load the key from SerializedBase64KeyMap. If there is no entry for the
// provided cipher suite the key is generated.
func LoadOrGenerateKey(suite oprf.Suite, serializedBase64KeyMap SerializedBase64KeyMap) (*oprf.PrivateKey, error) {
	var (
		privateKey *oprf.PrivateKey
		err        error
	)

	serializedBase64Key, ok := serializedBase64KeyMap[suite.Identifier()]
	if ok {
		// load a private key for the encryption suite
		privateKey, err = LoadPrivateKey(suite, serializedBase64Key)
	} else {
		// generate a private key for the encryption suite
		privateKey, err = oprf.GenerateKey(suite, rand.Reader)
	}

	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// SerializePublicKey is a wrapper to serialize a public key
func SerializePublicKey(key *oprf.PrivateKey) []byte {
	publicKey, err := key.Public().MarshalBinary()
	if err != nil {
		log.Println(err)
	}

	return publicKey
}

func DeserializeElements(blindedElements [][]byte, suiteGroup group.Group) ([]group.Element, error) {
	elements := make([]group.Element, len(blindedElements))

	for index, blindedElement := range blindedElements {
		element := suiteGroup.NewElement()

		if err := element.UnmarshalBinary(blindedElement); err != nil {
			log.Println(err)

			return nil, err
		}

		elements[index] = element
	}

	return elements, nil
}
