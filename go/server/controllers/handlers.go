package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/cloudflare/circl/oprf"
	"github.com/labstack/echo/v4"
)

// EvaluationRequest represents an evaluation requests
type EvaluationRequest struct {
	Suite           string    `json:"suite"`
	Info            string    `json:"info"`
	BlindedElements [][]byte  `json:"blinded_elements"`
	Mode            oprf.Mode `json:"mode"`
}

type EvaluationResponse struct {
	Evaluation          *oprf.Evaluation `json:"evaluation"`
	Suite               string           `json:"suite"`
	SerializedPublicKey []byte           `json:"serialized_public_key"`
}

// WrappedEvaluationResponse allows to partially parse the JSON input
type WrappedEvaluationResponse struct {
	Evaluation          *WrappedEvaluation `json:"evaluation"`
	Suite               string             `json:"suite"`
	SerializedPublicKey []byte             `json:"serialized_public_key"`
}

type WrappedEvaluation struct {
	Proof    []byte   `json:"proof"`
	Elements [][]byte `json:"elements"`
}

// UnmarshalJSON first parse the JSON input into a WrappedEvaluationResponse
// and then convert the string into a byte array into EvaluationResponse.
func (r *EvaluationResponse) MarshalJSON() ([]byte, error) {
	elements := r.Evaluation.Elements
	rawElements := make([][]byte, len(elements))

	for index, element := range elements {
		blindedElement, err := element.MarshalBinary()
		if err != nil {
			return nil, err
		}

		rawElements[index] = blindedElement
	}

	var (
		proof []byte
		err   error
	)

	if r.Evaluation.Proof != nil {
		proof, err = r.Evaluation.Proof.MarshalBinary()
		if err != nil {
			return nil, err
		}
	}

	return json.Marshal(WrappedEvaluationResponse{
		Evaluation: &WrappedEvaluation{
			Proof:    proof,
			Elements: rawElements,
		},
		Suite:               r.Suite,
		SerializedPublicKey: r.SerializedPublicKey,
	})
}

func NewEvaluationResponse(evaluation *oprf.Evaluation, suiteID string, serializedPublicKey []byte) *EvaluationResponse {
	return &EvaluationResponse{
		Evaluation:          evaluation,
		Suite:               suiteID,
		SerializedPublicKey: serializedPublicKey,
	}
}

// GetKeysHandler is an endpoint returning the static keys
func (s *OPRFServerController) GetKeysHandler(c echo.Context) error {
	s.keysMu.RLock()
	keys := make(map[string][]byte)

	for suiteID, privateKey := range s.keys {
		keys[suiteID] = SerializePublicKey(privateKey)
	}
	s.keysMu.RUnlock()

	return c.JSON(http.StatusOK, &keys) //nolint:wrapcheck
}

// EvaluateHandler is an endpoint that evaluate an EvaluationRequest.
// It returns an HTTP 400 Bad Request Error on incorrect input and
// an HTTP 500 Internal OPRFServerController Error if the evaluation fails.
// For instance :
// curl -X POST http://localhost:1323/api/evaluate -H 'Content-Type: application/json' -d \
// '{"suite": 3, "mode": 1, "info": "7465737420696e666f", "blinded_elements": \
// [[2, 99, 233, 95, 211, 165, 194, 204, 118, 22, 17, 134, 162, 84, 135, 138, 180, 7, \
// 229, 225, 238, 137, 138, 247, 196, 178, 119, 121, 218, 135, 36, 201, 132],[2, 61, 128, \
// 127, 32, 157, 20, 86, 131, 22, 159, 225, 197, 38, 118, 154, 158, 71, 70, 50, 188, 116, \
// 40, 80, 108, 72, 139, 91, 98, 146, 135, 105, 40]]}'
func (s *OPRFServerController) EvaluateHandler(c echo.Context) error {
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

	blindedElements, err := DeserializeElements(evaluationRequest.BlindedElements, server.Suite().Group())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Couldn't decode the blinded elements")
	}

	// Calculate the evaluation
	evaluation, err := server.Evaluate(
		&oprf.EvaluationRequest{
			Elements: blindedElements,
		},
		[]byte(evaluationRequest.Info),
	)
	if err != nil || evaluation == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Send the public key to the client for the finalization (needed for Serverless Functions)
	s.keysMu.RLock()
	serializedPublicKey := SerializePublicKey(s.keys[evaluationRequest.Suite])
	s.keysMu.RUnlock()

	response := NewEvaluationResponse(evaluation, server.Suite().Identifier(), serializedPublicKey)

	return c.JSON(http.StatusOK, response) //nolint:wrapcheck
}
