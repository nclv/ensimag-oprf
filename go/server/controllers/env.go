package controllers

import (
	"os"

	"github.com/cloudflare/circl/oprf"
)

const (
	EnvP256PrivateKey = "P256_PRIVATE_KEY"
	EnvP384PrivateKey = "P384_PRIVATE_KEY"
	EnvP521PrivateKey = "P521_PRIVATE_KEY"
)

// LoadPrivateKeysFromEnv load the base64 serialized private keys from the environment variables.
func LoadPrivateKeysFromEnv() SerializedBase64KeyMap {
	serializedBase64KeyMap := make(SerializedBase64KeyMap)

	serializedBase64P256PrivateKey, ok := os.LookupEnv(EnvP256PrivateKey)
	if ok {
		serializedBase64KeyMap[oprf.OPRFP256] = serializedBase64P256PrivateKey
	}

	serializedBase64P384PrivateKey, ok := os.LookupEnv(EnvP384PrivateKey)
	if ok {
		serializedBase64KeyMap[oprf.OPRFP384] = serializedBase64P384PrivateKey
	}

	serializedBase64P521PrivateKey, ok := os.LookupEnv(EnvP521PrivateKey)
	if ok {
		serializedBase64KeyMap[oprf.OPRFP521] = serializedBase64P521PrivateKey
	}

	return serializedBase64KeyMap
}
