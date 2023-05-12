package controllers

import (
	"sync"

	"github.com/cloudflare/circl/oprf"
)

type (
	// suite:base64 private key
	SerializedBase64KeyMap map[string]string
	// suite:private key
	KeyMap map[string]*oprf.PrivateKey
	// mode:suite:server
	ServerMap map[oprf.Mode]map[string]Server
)

type Server interface {
	Evaluate(req *oprf.EvaluationRequest, info []byte) (*oprf.Evaluation, error)
	FullEvaluate(input, info []byte) (output []byte, err error)
	VerifyFinalize(input, info, expectedOutput []byte) bool
	Suite() oprf.Suite
}

type BaseServer struct {
	oprf.Server
	suite oprf.Suite
}

func (s BaseServer) Suite() oprf.Suite {
	return s.suite
}

func (s BaseServer) Evaluate(req *oprf.EvaluationRequest, _ []byte) (*oprf.Evaluation, error) {
	return s.Server.Evaluate(req)
}

func (s BaseServer) FullEvaluate(input, _ []byte) ([]byte, error) {
	return s.Server.FullEvaluate(input)
}

func (s BaseServer) VerifyFinalize(input, _, expectedOutput []byte) bool {
	return s.Server.VerifyFinalize(input, expectedOutput)
}

type VerifiableServer struct {
	oprf.VerifiableServer
	suite oprf.Suite
}

func (s VerifiableServer) Suite() oprf.Suite {
	return s.suite
}

func (s VerifiableServer) Evaluate(req *oprf.EvaluationRequest, _ []byte) (*oprf.Evaluation, error) {
	return s.VerifiableServer.Evaluate(req)
}

func (s VerifiableServer) FullEvaluate(input, _ []byte) ([]byte, error) {
	return s.VerifiableServer.FullEvaluate(input)
}

func (s VerifiableServer) VerifyFinalize(input, _, expectedOutput []byte) bool {
	return s.VerifiableServer.VerifyFinalize(input, expectedOutput)
}

type PartialObliviousServer struct {
	oprf.PartialObliviousServer
	suite oprf.Suite
}

func (s PartialObliviousServer) Suite() oprf.Suite {
	return s.suite
}

// OPRFServerController holds the private keys and the servers
type OPRFServerController struct {
	keys      KeyMap
	servers   ServerMap
	keysMu    sync.RWMutex
	serversMu sync.RWMutex
}

func NewOPRFServerController() *OPRFServerController {
	controller := &OPRFServerController{ //nolint:exhaustivestruct
		keys:    make(KeyMap),
		servers: make(ServerMap),
	}
	controller.servers[oprf.BaseMode] = make(map[string]Server)
	controller.servers[oprf.VerifiableMode] = make(map[string]Server)
	controller.servers[oprf.PartialObliviousMode] = make(map[string]Server)

	return controller
}

// Initialize generate private keys and initialize the encryption's suite servers
func (s *OPRFServerController) Initialize(serializedBase64KeyMap SerializedBase64KeyMap) error {
	suites := []oprf.Suite{oprf.SuiteP256, oprf.SuiteP384, oprf.SuiteP521}
	for _, suite := range suites {
		suiteID := suite.Identifier()

		privateKey, err := LoadOrGenerateKey(suite, serializedBase64KeyMap)
		if err != nil {
			return err
		}

		s.keys[suiteID] = privateKey

		// create the base and verifiable servers for a provided encryption suite.
		baseModeServer := oprf.NewServer(suite, privateKey)
		verifiableModeServer := oprf.NewVerifiableServer(suite, privateKey)
		partialObliviousModeServer := oprf.NewPartialObliviousServer(suite, privateKey)

		s.servers[oprf.BaseMode][suiteID] = BaseServer{baseModeServer, suite}
		s.servers[oprf.VerifiableMode][suiteID] = VerifiableServer{verifiableModeServer, suite}
		s.servers[oprf.PartialObliviousMode][suiteID] = PartialObliviousServer{partialObliviousModeServer, suite}
	}

	return nil
}
