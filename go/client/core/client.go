package core

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/oprf/go/common"

	"github.com/cloudflare/circl/oprf"
)

// Client regroup an HTTP client and an OPRF client
type Client struct {
	serverURL string

	httpClient *http.Client
	oprfClient *oprf.Client

	publicKeys map[oprf.SuiteID]*oprf.PublicKey
}

// NewClient returns a new client
func NewClient(serverURL string) *Client {
	// TODO: https client with http.Transport
	return &Client{
		serverURL:  serverURL,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// SetupOPRFClient retrieve the server's public keys and create the OPRF client.
func (c *Client) SetupOPRFClient(suite oprf.SuiteID, mode oprf.Mode) error {
	serializedPublicKeys := c.GetPublicKeys()
	publicKeys := DeserializePublicKeys(serializedPublicKeys)

	c.publicKeys = publicKeys

	oprfClient, err := NewOPRFClient(suite, mode, publicKeys[suite])
	if err != nil {
		log.Println("error when setting up the OPRF client", err)

		return err
	}

	c.oprfClient = oprfClient

	return nil
}

// GetPublicKeys returns the public keys from the server
func (c *Client) GetPublicKeys() map[oprf.SuiteID][]byte {
	req, err := http.NewRequest("GET", c.serverURL+PublicKeysEndpoint, nil)
	if err != nil {
		log.Println("HTTP NewRequest error :", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println("HTTP Do request :", err)
	}
	defer resp.Body.Close()

	var publicKeys map[oprf.SuiteID][]byte
	if err := json.NewDecoder(resp.Body).Decode(&publicKeys); err != nil {
		log.Println("JSON decoder error :", err)
	}

	// log.Println(publicKeys, base64.StdEncoding.EncodeToString(publicKeys[common.OPRFP256]))

	return publicKeys
}

// CreateRequest generates a request for the server by passing an array of inputs to be evaluated by server.
func (c *Client) CreateRequest(inputs [][]byte) *oprf.ClientRequest {
	clientRequest, err := c.oprfClient.Request(inputs)
	if err != nil {
		log.Println("OPRF client request creation error :", err)
	}

	return clientRequest
}

// EvaluateRequest evaluate a common.EvaluationRequest into an oprf.Evaluation
func (c *Client) EvaluateRequest(evaluationRequest *common.EvaluationRequest) *oprf.Evaluation {
	data, err := json.Marshal(&evaluationRequest)
	if err != nil {
		log.Println("evaluation request marshalling error :", err)
	}

	req, err := http.NewRequest("POST", c.serverURL+EvaluateEndpoint, bytes.NewBuffer(data))
	if err != nil {
		log.Println("HTTP NewRequest error :", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println("HTTP Do request :", err)
	}
	defer resp.Body.Close()

	var evaluation oprf.Evaluation
	if err = json.NewDecoder(resp.Body).Decode(&evaluation); err != nil {
		log.Println("JSON decoder error :", err)
	}

	return &evaluation
}

// Finalize computes the signed token from the server Evaluation and returns the output of the
// OPRF protocol. The function uses server's public key to verify the proof in verifiable mode.
func (c *Client) Finalize(clientRequest *oprf.ClientRequest,
	evaluation *oprf.Evaluation, info string) [][]byte {
	clientOutputs, err := c.oprfClient.Finalize(clientRequest, evaluation, []byte(info))
	if err != nil || clientOutputs == nil {
		log.Println("Finalize error :", err, clientOutputs)
	}

	return clientOutputs
}
