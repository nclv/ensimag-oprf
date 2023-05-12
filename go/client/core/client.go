package core

import (
	"fmt"
	"log"

	"github.com/cloudflare/circl/oprf"
)

// Client regroup an HTTP client and an OPRF client
type Client struct {
	httpClient *HTTPClient
	oprfClient OprfClientInterface
	// suite:public key
	publicKeys map[string]*oprf.PublicKey
}

// NewClient returns a new HTTP + OPRF client with the provided server URL, suite, mode and static key.
// No static key is needed for oprf.BaseMode.
func NewClient(serverURL string, suite oprf.Suite, mode oprf.Mode) *Client {
	client := &Client{
		httpClient: NewHttpClient(serverURL),
	}
	if err := client.setupOPRFClient(suite, mode); err != nil {
		log.Fatal(err)
	}

	return client
}

// SetupOPRFClient retrieve the server's public keys and create the OPRF client.
func (c *Client) setupOPRFClient(suite oprf.Suite, mode oprf.Mode) error {
	serializedPublicKeys, err := c.httpClient.GetPublicKeys()
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

	publicKey := publicKeys[suite.Identifier()]
	c.oprfClient = NewOPRFClient(suite, mode, publicKey)

	return nil
}

// Blind generates a request for the server by passing an array of inputs to be evaluated by server.
func (c *Client) Blind(inputs [][]byte) (*oprf.FinalizeData, *oprf.EvaluationRequest, error) {
	return c.oprfClient.Blind(inputs)
}

// DeterministicBlind generates a request for the server by passing an array of inputs and serialized blinds to be evaluated by server.
// Use FinalizedData.CopyBlinds() to get the serialized blinds from a previous call to client.Blind(inputs)
func (c *Client) DeterministicBlind(inputs [][]byte, blinds []oprf.Blind) (*oprf.FinalizeData, *oprf.EvaluationRequest, error) {
	return c.oprfClient.DeterministicBlind(inputs, blinds)
}

func (c *Client) EvaluateRequest(evaluationRequest *EvaluationRequest) (*EvaluationResponse, error) {
	return c.httpClient.EvaluateRequest(evaluationRequest)
}

// Finalize computes the signed token from the server Evaluation and returns the output of the
// OPRF protocol. The function uses server's static key to verify the proof in verifiable mode.
func (c *Client) Finalize(finalizeData *oprf.FinalizeData,
	evaluation *oprf.Evaluation, info string,
) ([][]byte, error) {
	clientOutputs, err := c.oprfClient.Finalize(finalizeData, evaluation, []byte(info))
	if err != nil || clientOutputs == nil {
		log.Println("Finalize error :", err, clientOutputs)

		return nil, fmt.Errorf("finalize error : %w", err)
	}

	return clientOutputs, nil
}
