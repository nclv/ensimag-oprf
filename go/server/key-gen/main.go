package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cloudflare/circl/oprf"
)

var (
	suiteID string
	help    bool
)

func init() {
	flag.StringVar(&suiteID, "suite", "P256-SHA256", "Cipher suite : P256-SHA256, P384-SHA384 or P521-SHA512")
	flag.BoolVar(&help, "help", false, "Show the usage")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	// Handle required flags
	if help {
		flag.Usage()
		os.Exit(0)
	}

	suite, err := oprf.GetSuite(suiteID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Suite :", suite)

	// Generate and serialize the private key
	privateKey, err := oprf.GenerateKey(suite, rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	serializedKey, err := privateKey.MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}

	// Show the base64 encoded key
	log.Println("Base64 encode private key :", base64.StdEncoding.EncodeToString(serializedKey))
}
