package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	PublicKeysEndpoint = "/request_public_keys"
	EvaluateEndpoint   = "/evaluate"
)

type HTTPClient struct {
	client    *http.Client
	serverURL string
}

func NewHttpClient(serverURL string) *HTTPClient {
	return &HTTPClient{
		serverURL: serverURL,
		client:    &http.Client{Timeout: 15 * time.Second},
	}
}

// GetPublicKeys returns the public keys from the server
func (c *HTTPClient) GetPublicKeys() (map[string][]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.serverURL+PublicKeysEndpoint, http.NoBody)
	if err != nil {
		log.Println("HTTP NewRequest error :", err)

		return nil, fmt.Errorf("HTTP NewRequest error : %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Println("HTTP Do :", err)

		return nil, fmt.Errorf("HTTP Do error : %w", err)
	}
	defer resp.Body.Close()

	var publicKeys map[string][]byte
	if err := json.NewDecoder(resp.Body).Decode(&publicKeys); err != nil {
		log.Println("JSON decoder error :", err)

		return nil, fmt.Errorf("JSON decoder error : %w", err)
	}

	// log.Println(publicKeys, base64.StdEncoding.EncodeToString(publicKeys[common.OPRFP256]))

	return publicKeys, nil
}

// EvaluateRequest evaluate an EvaluationRequest into an EvaluationResponse
func (c *HTTPClient) EvaluateRequest(evaluationRequest *EvaluationRequest) (*EvaluationResponse, error) {
	data, err := json.Marshal(&evaluationRequest)
	if err != nil {
		log.Println("evaluation request marshalling error :", err)

		return nil, fmt.Errorf("evaluation request marshalling error : %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.serverURL+EvaluateEndpoint, bytes.NewBuffer(data))
	if err != nil {
		log.Println("HTTP NewRequest error :", err)

		return nil, fmt.Errorf("HTTP NewRequest error : %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
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
