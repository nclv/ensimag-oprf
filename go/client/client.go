package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/cloudflare/circl/oprf"
)

// EvaluationRequest
type EvaluationRequest struct {
	Suite           oprf.SuiteID   `json:"suite"`
	Mode            oprf.Mode      `json:"mode"`
	Info            string         `json:"info"`
	BlindedElements []oprf.Blinded `json:"blinded_elements"` // or use []string
}

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
func (c *Client) SetupOPRFClient(suite oprf.SuiteID, mode oprf.Mode) {
	serializedPublicKeys := c.GetPublicKeys()
	publicKeys := DeserializePublicKeys(serializedPublicKeys)

	c.publicKeys = publicKeys
	c.oprfClient = NewOPRFClient(suite, mode, publicKeys[suite])
}

// GetPublicKeys returns the public keys from the server
func (c *Client) GetPublicKeys() map[oprf.SuiteID][]byte {
	req, err := http.NewRequest("GET", c.serverURL+PublicKeysEndpoint, nil)
	if err != nil {
		log.Println(err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	var publicKeys map[oprf.SuiteID][]byte
	if err := json.NewDecoder(resp.Body).Decode(&publicKeys); err != nil {
		log.Println(err)
	}

	// log.Println(publicKeys, base64.StdEncoding.EncodeToString(publicKeys[utils.OPRFP256]))

	return publicKeys
}

// DeserializePublicKeys deserialize the server's public keys
func DeserializePublicKeys(serializedPublicKeys map[oprf.SuiteID][]byte) map[oprf.SuiteID]*oprf.PublicKey {
	publicKeys := make(map[oprf.SuiteID]*oprf.PublicKey)

	for suiteID, serializedPublicKey := range serializedPublicKeys {
		log.Println(
			"Suite ID:", suiteID,
			", Public key : ", base64.StdEncoding.EncodeToString(serializedPublicKey),
		)

		publicKey := new(oprf.PublicKey)
		if err := publicKey.Deserialize(suiteID, serializedPublicKey); err != nil {
			log.Println(err)
		}

		publicKeys[suiteID] = publicKey
	}

	return publicKeys
}

// CreateRequest generates a request for the server by passing an array of inputs to be evaluated by server.
func (c *Client) CreateRequest(inputs [][]byte) *oprf.ClientRequest {
	clientRequest, err := c.oprfClient.Request(inputs)
	if err != nil {
		log.Println(err)
	}

	return clientRequest
}

// EvaluateRequest
func (c *Client) EvaluateRequest(evaluationRequest *EvaluationRequest) *oprf.Evaluation {
	data, err := json.Marshal(&evaluationRequest)
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest("POST", c.serverURL+EvaluateEndpoint, bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	var evaluation oprf.Evaluation
	if err = json.NewDecoder(resp.Body).Decode(&evaluation); err != nil {
		log.Println(err)
	}

	return &evaluation
}

// Finalize computes the signed token from the server Evaluation and returns the output of the
// OPRF protocol. The function uses server's public key to verify the proof in verifiable mode.
func (c *Client) Finalize(clientRequest *oprf.ClientRequest,
	evaluation *oprf.Evaluation, info string) [][]byte {
	clientOutputs, err := c.oprfClient.Finalize(clientRequest, evaluation, []byte(info))
	if err != nil || clientOutputs == nil {
		log.Println(err, clientOutputs)
	}

	return clientOutputs
}
