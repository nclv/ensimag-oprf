package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/cloudflare/circl/oprf"
)

// GeneratePrivateKey generate a private key for the encryption suite
func GeneratePrivateKey(suite oprf.SuiteID) (*oprf.PrivateKey, error) {
	privateKey, err := oprf.GenerateKey(suite, rand.Reader)
	if err != nil {
		log.Println(err)

		return nil, err
	}

	return privateKey, nil
}

// LoadPrivateKey decode the base64 serialized private key and deserialized it.
func LoadPrivateKey(suiteID oprf.SuiteID, serializedBase64Key string) (*oprf.PrivateKey, error) {
	serializedKey, err := base64.StdEncoding.DecodeString(serializedBase64Key)
	if err != nil {
		log.Println(err)

		return nil, fmt.Errorf("couldn't load the private key : %w", err)
	}

	privateKey, err := DeserializePrivateKey(suiteID, serializedKey)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// LoadOrGenerateKey tries to load the key from SerializedBase64KeyMap. If there is no entry for the
// provided cipher suite the key is generated.
func LoadOrGenerateKey(suite oprf.SuiteID,
	serializedBase64KeyMap SerializedBase64KeyMap) (*oprf.PrivateKey, error) {
	var (
		privateKey *oprf.PrivateKey
		err        error
	)

	serializedBase64Key, ok := serializedBase64KeyMap[suite]
	if ok {
		privateKey, err = LoadPrivateKey(suite, serializedBase64Key)
		if err != nil {
			return nil, err
		}
	} else {
		privateKey, err = GeneratePrivateKey(suite)
		if err != nil {
			return nil, err
		}
	}

	return privateKey, nil
}

// SerializePublicKey is a wrapper to serialize a public key
func SerializePublicKey(key *oprf.PrivateKey) []byte {
	publicKey, err := key.Public().Serialize()
	if err != nil {
		log.Println(err)
	}

	return publicKey
}

// DeserializePrivateKey deserialize a private key
func DeserializePrivateKey(suiteID oprf.SuiteID, serializedPrivateKey []byte) (*oprf.PrivateKey, error) {
	privateKey := new(oprf.PrivateKey)
	if err := privateKey.Deserialize(suiteID, serializedPrivateKey); err != nil {
		log.Println(err)

		return nil, err
	}

	return privateKey, nil
}

// NewOPRFServer create an OPRF server with the provided suite, mode and private key
func NewOPRFServer(suite oprf.SuiteID, mode oprf.Mode, privateKey *oprf.PrivateKey) (*oprf.Server, error) {
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

		return nil, err
	}

	return server, nil
}
