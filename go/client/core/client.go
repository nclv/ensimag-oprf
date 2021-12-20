package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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

// SetOPRFClientPublicKey update the OPRF client to use the provided public key.
// The public key is only needed for the verifiable mode so this method does nothing in base mode.
func (c *Client) SetOPRFClientPublicKey(publicKey *oprf.PublicKey) error {
	if c.oprfClient.Mode == oprf.VerifiableMode {
		oprfClient, err := oprf.NewVerifiableClient(c.oprfClient.SuiteID, publicKey)
		if err != nil {
			log.Println("error when setting OPRF client public key", err)

			return err
		}

		c.oprfClient = oprfClient
	}

	return nil
}

// SetupOPRFClient retrieve the server's public keys and create the OPRF client.
func (c *Client) SetupOPRFClient(suite oprf.SuiteID, mode oprf.Mode) error {
	serializedPublicKeys, err := c.GetPublicKeys()
	if err != nil {
		log.Println("error when getting the static keys")

		return err
	}

	publicKeys, err := DeserializePublicKeys(serializedPublicKeys)
	if err != nil {
		log.Println("couldn't deserialize public keys :", err)

		return fmt.Errorf("couldn't deserialize public keys")
	}

	c.publicKeys = publicKeys
	publicKey := publicKeys[suite]

	oprfClient, err := NewOPRFClient(suite, mode, publicKey)
	if err != nil {
		log.Println("error when setting up the OPRF client")

		return err
	}

	c.oprfClient = oprfClient

	return nil
}

// GetPublicKeys returns the public keys from the server
func (c *Client) GetPublicKeys() (map[oprf.SuiteID][]byte, error) {
	req, err := http.NewRequest("GET", c.serverURL+PublicKeysEndpoint, nil)
	if err != nil {
		log.Println("HTTP NewRequest error :", err)

		return nil, fmt.Errorf("HTTP NewRequest error : %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println("HTTP Do :", err)

		return nil, fmt.Errorf("HTTP Do error : %w", err)
	}
	defer resp.Body.Close()

	var publicKeys map[oprf.SuiteID][]byte
	if err := json.NewDecoder(resp.Body).Decode(&publicKeys); err != nil {
		log.Println("JSON decoder error :", err)

		return nil, fmt.Errorf("JSON decoder error : %w", err)
	}

	// log.Println(publicKeys, base64.StdEncoding.EncodeToString(publicKeys[common.OPRFP256]))

	return publicKeys, nil
}

// CreateRequest generates a request for the server by passing an array of inputs to be evaluated by server.
func (c *Client) CreateRequest(inputs [][]byte) (*oprf.ClientRequest, error) {
	clientRequest, err := c.oprfClient.Request(inputs)
	if err != nil {
		log.Println("OPRF client request creation error :", err)

		return nil, fmt.Errorf("OPRF client request creation error : %w", err)
	}

	return clientRequest, nil
}

// EvaluateRequest evaluate an EvaluationRequest into an EvaluationResponse
func (c *Client) EvaluateRequest(evaluationRequest *EvaluationRequest) (*EvaluationResponse, error) {
	data, err := json.Marshal(&evaluationRequest)
	if err != nil {
		log.Println("evaluation request marshalling error :", err)

		return nil, fmt.Errorf("evaluation request marshalling error : %w", err)
	}

	req, err := http.NewRequest("POST", c.serverURL+EvaluateEndpoint, bytes.NewBuffer(data))
	if err != nil {
		log.Println("HTTP NewRequest error :", err)

		return nil, fmt.Errorf("HTTP NewRequest error : %w", err)

	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Println("HTTP Do :", err)

		return nil, fmt.Errorf("HTTP Do error : %w", err)
	}
	defer resp.Body.Close()

	var evaluationResponse EvaluationResponse
	if err = json.NewDecoder(resp.Body).Decode(&evaluationResponse); err != nil {
		log.Println("JSON decoder error :", err)

		return nil, fmt.Errorf("JSON decoder error : %w", err)
	}

	return &evaluationResponse, nil
}

// Finalize computes the signed token from the server Evaluation and returns the output of the
// OPRF protocol. The function uses server's static key to verify the proof in verifiable mode.
func (c *Client) Finalize(clientRequest *oprf.ClientRequest,
	evaluation *oprf.Evaluation, info string) ([][]byte, error) {
	clientOutputs, err := c.oprfClient.Finalize(clientRequest, evaluation, []byte(info))
	if err != nil || clientOutputs == nil {
		log.Println("Finalize error :", err, clientOutputs)

		return nil, fmt.Errorf("finalize error : %w", err)
	}

	return clientOutputs, nil
}
