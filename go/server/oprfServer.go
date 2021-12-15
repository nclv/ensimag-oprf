package main

import (
	"github.com/oprf/go/common"
	"net/http"
	"sync"

	"github.com/cloudflare/circl/oprf"
	"github.com/labstack/echo/v4"
)

type (
	KeyMap    map[oprf.SuiteID]*oprf.PrivateKey
	ServerMap map[oprf.Mode]map[oprf.SuiteID]*oprf.Server
)

type OPRFServerManager struct {
	// Private keys
	keys   KeyMap
	keysMu sync.RWMutex
	// Base and verifiable servers for the allowed encryption suites
	servers   ServerMap
	serversMu sync.RWMutex
}

func NewOPRFServerManager() *OPRFServerManager {
	return &OPRFServerManager{ //nolint:exhaustivestruct
		keys:    make(KeyMap),
		servers: make(ServerMap),
	}
}

// Initialize generate private keys and initialize the encryption's suite servers
func (s *OPRFServerManager) Initialize() {
	s.generatePrivateKeys()
	s.initializeServers()
}

// generatePrivateKeys generate the private keys
func (s *OPRFServerManager) generatePrivateKeys() {
	s.keysMu.Lock()
	s.keys[oprf.OPRFP256] = GeneratePrivateKey(oprf.OPRFP256)
	s.keys[oprf.OPRFP384] = GeneratePrivateKey(oprf.OPRFP384)
	s.keys[oprf.OPRFP521] = GeneratePrivateKey(oprf.OPRFP521)
	s.keysMu.Unlock()
}

// createServersSuite create the base and verifiable servers for a provided encryption suite.
func (s *OPRFServerManager) createServersSuite(suite oprf.SuiteID) {
	s.keysMu.RLock()
	privateKey := s.keys[suite]
	s.keysMu.RUnlock()

	s.servers[oprf.BaseMode][suite] = NewOPRFServer(suite, oprf.BaseMode, privateKey)
	s.servers[oprf.VerifiableMode][suite] = NewOPRFServer(suite, oprf.VerifiableMode, privateKey)
}

// initializeServers initialize the servers for the allowed encryption suite
func (s *OPRFServerManager) initializeServers() {
	s.serversMu.Lock()
	// Create the nested sub-maps
	s.servers[oprf.BaseMode] = make(map[oprf.SuiteID]*oprf.Server)
	s.servers[oprf.VerifiableMode] = make(map[oprf.SuiteID]*oprf.Server)

	s.createServersSuite(oprf.OPRFP256)
	s.createServersSuite(oprf.OPRFP384)
	s.createServersSuite(oprf.OPRFP521)

	s.serversMu.Unlock()
}

/*
OPRFServerManager endpoints
*/

// getKeysHandler is an endpoint returning the public keys
func (s *OPRFServerManager) getKeysHandler(c echo.Context) error {
	s.keysMu.RLock()
	keys := make(map[oprf.SuiteID][]byte)

	for suiteID, privateKey := range s.keys {
		keys[suiteID] = SerializePublicKey(privateKey)
	}
	s.keysMu.RUnlock()

	return c.JSON(http.StatusOK, &keys)
}

// evaluateHandler is an endpoint that evaluateHandler an EvaluationRequest.
// It returns an HTTP 400 Bad Request Error on incorrect input and
// an HTTP 500 Internal OPRFServerManager Error if the evaluation fails.
// For instance :
// curl -X POST http://localhost:1323/evaluate -H 'Content-Type: application/json' -d \
// '{"suite": 3, "mode": 1, "info": "7465737420696e666f", "blinded_elements": \
// [[2, 99, 233, 95, 211, 165, 194, 204, 118, 22, 17, 134, 162, 84, 135, 138, 180, 7, \
// 229, 225, 238, 137, 138, 247, 196, 178, 119, 121, 218, 135, 36, 201, 132],[2, 61, 128, \
// 127, 32, 157, 20, 86, 131, 22, 159, 225, 197, 38, 118, 154, 158, 71, 70, 50, 188, 116, \
// 40, 80, 108, 72, 139, 91, 98, 146, 135, 105, 40]]}'
func (s *OPRFServerManager) evaluateHandler(c echo.Context) error {
	evaluationRequest := new(common.EvaluationRequest)
	if err := c.Bind(evaluationRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	s.serversMu.RLock()
	server := s.servers[evaluationRequest.Mode][evaluationRequest.Suite]
	s.serversMu.RUnlock()

	if server == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "No server")
	}

	evaluation, err := server.Evaluate(evaluationRequest.BlindedElements, []byte(evaluationRequest.Info))
	if err != nil || evaluation == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, evaluation) //nolint:wrapcheck
}
