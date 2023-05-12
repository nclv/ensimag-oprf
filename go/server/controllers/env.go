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

func GetEnvPrivateKeySuiteMap() map[string]string {
	return map[string]string{
		EnvP256PrivateKey: oprf.SuiteP256.Identifier(),
		EnvP384PrivateKey: oprf.SuiteP384.Identifier(),
		EnvP521PrivateKey: oprf.SuiteP521.Identifier(),
	}
}

// LoadPrivateKeysFromEnv load the base64 serialized private keys from the environment variables.
func LoadPrivateKeysFromEnv() SerializedBase64KeyMap {
	serializedBase64KeyMap := make(SerializedBase64KeyMap)
	envPrivateKeySuiteMap := GetEnvPrivateKeySuiteMap()

	for envPrivateKey, suiteID := range envPrivateKeySuiteMap {
		serializedBase64PrivateKey, ok := os.LookupEnv(envPrivateKey)
		if ok {
			serializedBase64KeyMap[suiteID] = serializedBase64PrivateKey
		}
	}

	return serializedBase64KeyMap
}
