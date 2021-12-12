package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cloudflare/circl/oprf"
)

const (
	ServerUrl = "http://localhost:1323"
)

type Client struct {
	httpClient *http.Client
	oprfClient *oprf.Client
}

func NewClient() *Client {
	// TODO: https client with http.Transport
	return &Client{httpClient: &http.Client{Timeout: 15 * time.Second}}
}

func (c *Client) GetPublicKeys() map[oprf.SuiteID][]byte {
	req, err := http.NewRequest("GET", ServerUrl+"/request_public_keys", nil)
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
		fmt.Println(err)
	}

	// log.Println(publicKeys, base64.StdEncoding.EncodeToString(publicKeys[oprf.OPRFP256]))

	return publicKeys
}

func (c *Client) SetOPRFClient(suite oprf.SuiteID, mode oprf.Mode, pkS *oprf.PublicKey) {
	var (
		client *oprf.Client
		err    error
	)

	if mode == oprf.BaseMode {
		client, err = oprf.NewClient(suite)
	} else if mode == oprf.VerifiableMode {
		client, err = oprf.NewVerifiableClient(suite, pkS)
	}

	if err != nil {
		log.Println(err)
	}

	c.oprfClient = client
}

// CreateRequest generates a request for server passing an array of inputs to be evaluated by server.
func (c *Client) CreateRequest(inputs [][]byte) *oprf.ClientRequest {
	clientRequest, err := c.oprfClient.Request(inputs)
	if err != nil {
		log.Println(err)
	}

	return clientRequest
}

type EvaluationRequest struct {
	Suite           oprf.SuiteID   `json:"suite"`
	Mode            oprf.Mode      `json:"mode"`
	Info            string         `json:"info"`
	BlindedElements []oprf.Blinded `json:"blinded_elements"` // or use []string
}

func (c *Client) EvaluateRequest(evaluationRequest *EvaluationRequest) *oprf.Evaluation {
	data, err := json.Marshal(&evaluationRequest)
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest("POST", ServerUrl+"/evaluate", bytes.NewBuffer(data))
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
		fmt.Println(err)
	}

	return &evaluation
}

// Finalize computes the signed token from the server Evaluation and returns the output of the
// OPRF protocol. The function uses server's public key to verify the proof in verifiable mode.
func (c Client) Finalize(clientRequest *oprf.ClientRequest,
	evaluation *oprf.Evaluation, info string) [][]byte {
	clientOutputs, err := c.oprfClient.Finalize(clientRequest, evaluation, []byte(info))
	if err != nil || clientOutputs == nil {
		log.Println(err)
	}

	return clientOutputs
}

var (
	mode  oprf.Mode
	suite oprf.SuiteID
)

func init() {
	modeFlag := flag.Int("mode", int(oprf.BaseMode), "mode")
	suiteFlag := flag.Int("suite", int(oprf.OPRFP256), "cipher suite")

	flag.Parse()

	mode = oprf.Mode(*modeFlag)
	suite = oprf.SuiteID(*suiteFlag)
}

func main() {
	log.Println(mode, suite)

	client := NewClient()

	publicKeys := client.GetPublicKeys()
	log.Println(publicKeys)

	publicKey := new(oprf.PublicKey)
	if err := publicKey.Deserialize(suite, publicKeys[suite]); err != nil {
		log.Println(err)
	}

	client.SetOPRFClient(suite, mode, publicKey)
	clientRequest := client.CreateRequest([][]byte{{0x00}, {0xFF}})

	info := "7465737420696e666f"
	evaluationRequest := &EvaluationRequest{
		Suite:           suite,
		Mode:            mode,
		Info:            info,
		BlindedElements: clientRequest.BlindedElements(),
	}
	evaluation := client.EvaluateRequest(evaluationRequest)
	log.Println(evaluation)

	outputs := client.Finalize(clientRequest, evaluation, info)
	log.Println(outputs)
}
