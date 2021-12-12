package main

import (
	"crypto/rand"
	"log"
	"net/http"

	"github.com/cloudflare/circl/oprf"
	"github.com/labstack/echo/v4"
)

var privateKeyFP256 *oprf.PrivateKey
var privateKeyFP384 *oprf.PrivateKey
var privateKeyFP521 *oprf.PrivateKey
var serverFP256BaseMode *oprf.Server
var serverFP256VerifiableMode *oprf.Server
var serverFP384BaseMode *oprf.Server
var serverFP384VerifiableMode *oprf.Server
var serverFP521BaseMode *oprf.Server
var serverFP521VerifiableMode *oprf.Server

func init() {
	initKeys()
	initServers()
}

func initKey(suite oprf.SuiteID) *oprf.PrivateKey {
	privateKey, err := oprf.GenerateKey(suite, rand.Reader)
	if err != nil {
		log.Println(err)
	}

	printKey(privateKey)
	printKey(privateKey.Public())

	return privateKey
}

// init_keys creates list of the different key for each mode
func initKeys() {
	privateKeyFP256 = initKey(oprf.OPRFP256)
	privateKeyFP384 = initKey(oprf.OPRFP384)
	privateKeyFP521 = initKey(oprf.OPRFP521)
}

func GetServer(suite oprf.SuiteID, mode oprf.Mode, privateKey *oprf.PrivateKey) *oprf.Server {
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
	}

	return server
}

func initServers() {
	serverFP256BaseMode = GetServer(oprf.OPRFP256, oprf.BaseMode, privateKeyFP256)
	serverFP256VerifiableMode = GetServer(oprf.OPRFP256, oprf.VerifiableMode, privateKeyFP256)
	serverFP384BaseMode = GetServer(oprf.OPRFP384, oprf.BaseMode, privateKeyFP384)
	serverFP384VerifiableMode = GetServer(oprf.OPRFP384, oprf.VerifiableMode, privateKeyFP384)
	serverFP521BaseMode = GetServer(oprf.OPRFP521, oprf.BaseMode, privateKeyFP521)
	serverFP521VerifiableMode = GetServer(oprf.OPRFP521, oprf.VerifiableMode, privateKeyFP521)
}

func serialize_key(key *oprf.PrivateKey) []byte {
	publicKey, err := key.Public().Serialize()
	if err != nil {
		log.Println(err)
	}
	return publicKey
}

func getKeys(c echo.Context) error {
	keys := &[][]byte{
		serialize_key(privateKeyFP256),
		serialize_key(privateKeyFP384),
		serialize_key(privateKeyFP521),
	}

	return c.JSON(http.StatusOK, keys)
}

type EvaluationRequest struct {
	Suite           oprf.SuiteID `json:"suite"`
	Mode            oprf.Mode    `json:"mode"`
	BlindedElements [][]byte     `json:"blinded_elements"` // or use []string
}

func evaluate(c echo.Context) error {
	evaluation := new(EvaluationRequest)
	if err := c.Bind((evaluation)); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	log.Println(evaluation)

	return c.JSON(http.StatusOK, evaluation)
}

func main() {
	e := echo.New()

	e.GET("/request_public_keys", getKeys)
	e.POST("/evaluate", evaluate)

	e.Logger.Fatal(e.Start(":1323"))
}
