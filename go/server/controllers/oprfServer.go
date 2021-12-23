package controllers

import (
	"sync"

	"github.com/cloudflare/circl/oprf"
)

type (
	SerializedBase64KeyMap map[oprf.SuiteID]string
	KeyMap                 map[oprf.SuiteID]*oprf.PrivateKey
	ServerMap              map[oprf.Mode]map[oprf.SuiteID]*oprf.Server
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
func (s *OPRFServerController) Initialize(serializedBase64KeyMap SerializedBase64KeyMap) error {
	// get the private key map
	privateKeyMap, err := s.getPrivateKeys(serializedBase64KeyMap)
	if err != nil {
		return err
	}
	// set the keys
	s.setKeys(privateKeyMap)

	// initialize the servers with the private keys
	if err := s.initializeServers(); err != nil {
		return err
	}

	return nil
}

func (s *OPRFServerController) setKeys(privateKeyMap KeyMap) {
	s.keysMu.Lock()
	s.keys[oprf.OPRFP256] = privateKeyMap[oprf.OPRFP256]
	s.keys[oprf.OPRFP384] = privateKeyMap[oprf.OPRFP384]
	s.keys[oprf.OPRFP521] = privateKeyMap[oprf.OPRFP521]
	s.keysMu.Unlock()
}

// getPrivateKeys return the private key map
func (s *OPRFServerController) getPrivateKeys(serializedBase64KeyMap SerializedBase64KeyMap) (KeyMap, error) {
	privateKeyMap := make(KeyMap)

	p256PrivateKey, err := s.loadOrGenerateKey(oprf.OPRFP256, serializedBase64KeyMap)
	if err != nil {
		return nil, err
	}
	privateKeyMap[oprf.OPRFP256] = p256PrivateKey

	p384PrivateKey, err := s.loadOrGenerateKey(oprf.OPRFP384, serializedBase64KeyMap)
	if err != nil {
		return nil, err
	}
	privateKeyMap[oprf.OPRFP384] = p384PrivateKey

	p521PrivateKey, err := s.loadOrGenerateKey(oprf.OPRFP521, serializedBase64KeyMap)
	if err != nil {
		return nil, err
	}
	privateKeyMap[oprf.OPRFP521] = p521PrivateKey

	return privateKeyMap, nil
}

// loadOrGenerateKey tries to load the key from SerializedBase64KeyMap. If there is no entry for the
// provided cipher suite the key is generated.
func (s *OPRFServerController) loadOrGenerateKey(suite oprf.SuiteID,
	serializedBase64KeyMap SerializedBase64KeyMap) (*oprf.PrivateKey, error) {
	var (
		privateKey *oprf.PrivateKey
		err        error
	)

	serializedBase64Key, ok := serializedBase64KeyMap[suite]
	if ok {
		privateKey, err = LoadPrivateKey(suite, serializedBase64Key)
		if err != nil {
			return nil, err
		}
	} else {
		privateKey, err = GeneratePrivateKey(suite)
		if err != nil {
			return nil, err
		}
	}

	return privateKey, nil
}

// createServersSuite create the base and verifiable servers for a provided encryption suite.
func (s *OPRFServerController) createServersSuite(suite oprf.SuiteID) error {
	s.keysMu.RLock()
	privateKey := s.keys[suite]
	s.keysMu.RUnlock()

	baseModeServer, err := NewOPRFServer(suite, oprf.BaseMode, privateKey)
	if err != nil {
		return err
	}

	verifiableModeServer, err := NewOPRFServer(suite, oprf.VerifiableMode, privateKey)
	if err != nil {
		return err
	}

	s.servers[oprf.BaseMode][suite] = baseModeServer
	s.servers[oprf.VerifiableMode][suite] = verifiableModeServer

	return nil
}

// initializeServers initialize the servers for the allowed encryption suite
func (s *OPRFServerController) initializeServers() error {
	s.serversMu.Lock()
	// Create the nested sub-maps
	s.servers[oprf.BaseMode] = make(map[oprf.SuiteID]*oprf.Server)
	s.servers[oprf.VerifiableMode] = make(map[oprf.SuiteID]*oprf.Server)

	if err := s.createServersSuite(oprf.OPRFP256); err != nil {
		return err
	}

	if err := s.createServersSuite(oprf.OPRFP384); err != nil {
		return err
	}

	if err := s.createServersSuite(oprf.OPRFP521); err != nil {
		return err
	}

	s.serversMu.Unlock()

	return nil
}
