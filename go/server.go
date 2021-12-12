package main

import (
	"net/http"
	"sync"

	"github.com/cloudflare/circl/oprf"
	"github.com/labstack/echo/v4"
)

const (
	HOST = "localhost"
	PORT = "1323"
)

type (
	KeyMap    map[oprf.SuiteID]*oprf.PrivateKey
	ServerMap map[oprf.Mode]map[oprf.SuiteID]*oprf.Server
)

type Server struct {
	// Private keys
	keys   KeyMap
	keysMu sync.RWMutex
	// Base and verifiable servers for the allowed encryption suites
	servers   ServerMap
	serversMu sync.RWMutex
}

func NewServer() *Server {
	return &Server{
		keys:    make(KeyMap),
		servers: make(ServerMap),
	}
}

// Initialize the server
func (s *Server) Initialize() {
	s.generatePrivateKeys()
	s.initializeServers()
}

// generatePrivateKeys Initialize the private keys
func (s *Server) generatePrivateKeys() {
	s.keysMu.Lock()
	s.keys[oprf.OPRFP256] = generatePrivateKey(oprf.OPRFP256)
	s.keys[oprf.OPRFP384] = generatePrivateKey(oprf.OPRFP384)
	s.keys[oprf.OPRFP521] = generatePrivateKey(oprf.OPRFP521)
	s.keysMu.Unlock()
}

// createServersSuite create the base and verifiable servers for a provided encryption suite.
func (s *Server) createServersSuite(suite oprf.SuiteID) {
	s.keysMu.RLock()
	privateKey := s.keys[suite]
	s.keysMu.RUnlock()

	s.servers[oprf.BaseMode][suite] = createServer(suite, oprf.BaseMode, privateKey)
	s.servers[oprf.VerifiableMode][suite] = createServer(suite, oprf.VerifiableMode, privateKey)
}

// initializeServers Initialize the servers for the allowed encryption suite
func (s *Server) initializeServers() {
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
Server endpoints
*/

type EvaluationRequest struct {
	Suite           oprf.SuiteID   `json:"suite"`
	Mode            oprf.Mode      `json:"mode"`
	Info            string         `json:"info"`
	BlindedElements []oprf.Blinded `json:"blinded_elements"` // or use []string
}

// getKeys is an endpoint returning the public keys
func (s *Server) getKeys(c echo.Context) error {
	s.keysMu.RLock()
	keys := make(map[oprf.SuiteID][]byte)

	for suiteID, privateKey := range s.keys {
		keys[suiteID] = serializePublicKey(privateKey)
	}
	s.keysMu.RUnlock()

	return c.JSON(http.StatusOK, &keys)
}

// evaluate is an endpoint that evaluate an EvaluationRequest.
// It returns an HTTP 400 Bad Request Error on incorrect input and
// an HTTP 500 Internal Server Error if the evaluation fails.
// For instance :
// curl -X POST http://localhost:1323/evaluate -H 'Content-Type: application/json' -d \
// '{"suite": 3, "mode": 1, "info": "7465737420696e666f", "blinded_elements": \
// [[2, 99, 233, 95, 211, 165, 194, 204, 118, 22, 17, 134, 162, 84, 135, 138, 180, 7, \
// 229, 225, 238, 137, 138, 247, 196, 178, 119, 121, 218, 135, 36, 201, 132],[2, 61, 128, \
// 127, 32, 157, 20, 86, 131, 22, 159, 225, 197, 38, 118, 154, 158, 71, 70, 50, 188, 116, \
// 40, 80, 108, 72, 139, 91, 98, 146, 135, 105, 40]]}'
func (s *Server) evaluate(c echo.Context) error {
	evaluationRequest := new(EvaluationRequest)
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

	return c.JSON(http.StatusOK, evaluation)
}

func main() {
	e := echo.New()

	// TODO: https://echo.labstack.com/cookbook/auto-tls/

	server := NewServer()
	server.Initialize()

	e.GET("/request_public_keys", server.getKeys)
	e.POST("/evaluate", server.evaluate)

	e.Logger.Fatal(e.Start(HOST + ":" + PORT))
}
