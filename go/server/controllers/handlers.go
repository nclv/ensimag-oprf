package controllers

import (
	"net/http"

	"github.com/cloudflare/circl/oprf"
	"github.com/labstack/echo/v4"
)

// EvaluationRequest represents an evaluation requests
type EvaluationRequest struct {
	Suite           oprf.SuiteID   `json:"suite"`
	Mode            oprf.Mode      `json:"mode"`
	Info            string         `json:"info"`
	BlindedElements []oprf.Blinded `json:"blinded_elements"` // or use []string
}

// IndexHandler handles the index.html template
func IndexHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil) //nolint:wrapcheck
}

// GetKeysHandler is an endpoint returning the static keys
func (s *OPRFServerController) GetKeysHandler(c echo.Context) error {
	s.keysMu.RLock()
	keys := make(map[oprf.SuiteID][]byte)

	for suiteID, privateKey := range s.keys {
		keys[suiteID] = SerializePublicKey(privateKey)
	}
	s.keysMu.RUnlock()

	return c.JSON(http.StatusOK, &keys)
}

// EvaluateHandler is an endpoint that evaluate an EvaluationRequest.
// It returns an HTTP 400 Bad Request Error on incorrect input and
// an HTTP 500 Internal OPRFServerController Error if the evaluation fails.
// For instance :
// curl -X POST http://localhost:1323/evaluate -H 'Content-Type: application/json' -d \
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

	evaluation, err := server.Evaluate(evaluationRequest.BlindedElements, []byte(evaluationRequest.Info))
	if err != nil || evaluation == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, evaluation) //nolint:wrapcheck
}
