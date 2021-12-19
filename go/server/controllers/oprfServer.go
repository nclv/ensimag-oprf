package controllers

import (
	"sync"

	"github.com/cloudflare/circl/oprf"
)

type (
	KeyMap    map[oprf.SuiteID]*oprf.PrivateKey
	ServerMap map[oprf.Mode]map[oprf.SuiteID]*oprf.Server
)

type OPRFServerController struct {
	// Private keys
	keys   KeyMap
	keysMu sync.RWMutex
	// Base and verifiable servers for the allowed encryption suites
	servers   ServerMap
	serversMu sync.RWMutex
}

func NewOPRFServerController() *OPRFServerController {
	return &OPRFServerController{ //nolint:exhaustivestruct
		keys:    make(KeyMap),
		servers: make(ServerMap),
	}
}

// Initialize generate private keys and initialize the encryption's suite servers
func (s *OPRFServerController) Initialize() {
	s.generatePrivateKeys()
	s.initializeServers()
}

// generatePrivateKeys generate the private keys
func (s *OPRFServerController) generatePrivateKeys() {
	s.keysMu.Lock()
	s.keys[oprf.OPRFP256] = GeneratePrivateKey(oprf.OPRFP256)
	s.keys[oprf.OPRFP384] = GeneratePrivateKey(oprf.OPRFP384)
	s.keys[oprf.OPRFP521] = GeneratePrivateKey(oprf.OPRFP521)
	s.keysMu.Unlock()
}

// createServersSuite create the base and verifiable servers for a provided encryption suite.
func (s *OPRFServerController) createServersSuite(suite oprf.SuiteID) {
	s.keysMu.RLock()
	privateKey := s.keys[suite]
	s.keysMu.RUnlock()

	s.servers[oprf.BaseMode][suite] = NewOPRFServer(suite, oprf.BaseMode, privateKey)
	s.servers[oprf.VerifiableMode][suite] = NewOPRFServer(suite, oprf.VerifiableMode, privateKey)
}

// initializeServers initialize the servers for the allowed encryption suite
func (s *OPRFServerController) initializeServers() {
	s.serversMu.Lock()
	// Create the nested sub-maps
	s.servers[oprf.BaseMode] = make(map[oprf.SuiteID]*oprf.Server)
	s.servers[oprf.VerifiableMode] = make(map[oprf.SuiteID]*oprf.Server)

	s.createServersSuite(oprf.OPRFP256)
	s.createServersSuite(oprf.OPRFP384)
	s.createServersSuite(oprf.OPRFP521)

	s.serversMu.Unlock()
}
