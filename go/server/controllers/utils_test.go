package controllers

import (
	"encoding/base64"
	"testing"

	"github.com/cloudflare/circl/oprf"
)

func TestLoadPrivateKey(t *testing.T) {
	serializedBase64Key := "AtzyGS8NoBjEjqbhwdGY/zWyqdFkJghyTttoIGq4UoM="

	privateKey, err := LoadPrivateKey(oprf.SuiteP256, serializedBase64Key)
	if err != nil {
		t.Error(err)
	}

	serializedKey, err := privateKey.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	if base64.StdEncoding.EncodeToString(serializedKey) != serializedBase64Key {
		t.Errorf("the loaded private key does not match")
	}
}
