package main

import (
	"encoding/base64"
	"flag"
	"log"

	"github.com/cloudflare/circl/oprf"

	"github.com/ensimag-oprf/go/server/controllers"
)

var (
	suite oprf.SuiteID
)

func commandLine() {
	suiteFlag := flag.Int("suite", int(oprf.OPRFP256), "cipher suite")

	flag.Parse()

	suite = oprf.SuiteID(*suiteFlag)
}

func main() {
	commandLine()

	log.Println(suite)

	// Generate and serialize the private key
	privateKey, err := controllers.GeneratePrivateKey(suite)
	if err != nil {
		log.Println(err)

		return
	}

	serializedKey, err := privateKey.Serialize()
	if err != nil {
		log.Println(err)

		return
	}

	// Show the base64 encoded key
	log.Println(base64.StdEncoding.EncodeToString(serializedKey))
}
