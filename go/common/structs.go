package common

import "github.com/cloudflare/circl/oprf"

// EvaluationRequest represents an evaluation requests
type EvaluationRequest struct {
	Suite           oprf.SuiteID   `json:"suite"`
	Mode            oprf.Mode      `json:"mode"`
	Info            string         `json:"info"`
	BlindedElements []oprf.Blinded `json:"blinded_elements"` // or use []string
}

func NewEvaluationRequest(suite oprf.SuiteID, mode oprf.Mode, info string,
	blindedElements []oprf.Blinded) *EvaluationRequest {
	return &EvaluationRequest{
		Suite:           suite,
		Mode:            mode,
		Info:            info,
		BlindedElements: blindedElements,
	}
}
