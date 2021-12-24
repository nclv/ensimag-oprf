package controllers

import (
	"encoding/base64"
	"github.com/cloudflare/circl/oprf"
	"testing"
)

func TestLoadPrivateKey(t *testing.T) {
	serializedBase64Key := "AtzyGS8NoBjEjqbhwdGY/zWyqdFkJghyTttoIGq4UoM="
	privateKey, err := LoadPrivateKey(oprf.OPRFP256, serializedBase64Key)
	if err != nil {
		t.Error(err)
	}

	serializedKey, err := privateKey.Serialize()
	if err != nil {
		t.Error(err)
	}

	if base64.StdEncoding.EncodeToString(serializedKey) != serializedBase64Key {
		t.Errorf("the loaded private key does not match")
	}
}
