package core

import (
	"encoding/json"

	"github.com/cloudflare/circl/group"
	"github.com/cloudflare/circl/oprf"
	"github.com/cloudflare/circl/zk/dleq"
)

// EvaluationRequest represents an evaluation requests
type EvaluationRequest struct {
	Suite           string    `json:"suite"`
	Info            string    `json:"info"`
	BlindedElements [][]byte  `json:"blinded_elements"`
	Mode            oprf.Mode `json:"mode"`
}

func NewEvaluationRequest(suite oprf.Suite, mode oprf.Mode, info string,
	blindedElements [][]byte,
) *EvaluationRequest {
	return &EvaluationRequest{
		Suite:           suite.Identifier(),
		Mode:            mode,
		Info:            info,
		BlindedElements: blindedElements,
	}
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

// EvaluationResponse contains the oprf.Evaluation and the serialized public key that will be used for the
// finalization.
type EvaluationResponse struct {
	Evaluation          *oprf.Evaluation `json:"evaluation"`
	Suite               string           `json:"suite"`
	SerializedPublicKey []byte           `json:"serialized_public_key"`
}

// UnmarshalJSON first parse the JSON input into a WrappedEvaluationResponse
// and then convert the string into a byte array into EvaluationResponse.
func (r *EvaluationResponse) UnmarshalJSON(data []byte) error {
	var wr WrappedEvaluationResponse
	if err := json.Unmarshal(data, &wr); err != nil {
		return err
	}

	r.SerializedPublicKey = wr.SerializedPublicKey
	r.Evaluation = &oprf.Evaluation{}

	suite, err := oprf.GetSuite(wr.Suite)
	if err != nil {
		return err
	}

	suiteGroup := suite.Group()

	var proof *dleq.Proof
	if len(wr.Evaluation.Proof) > 0 {
		proof = &dleq.Proof{}
		if err := proof.UnmarshalBinary(suiteGroup, wr.Evaluation.Proof); err != nil {
			return err
		}
	}

	r.Evaluation.Proof = proof

	elements := make([]group.Element, len(wr.Evaluation.Elements))

	for index, data := range wr.Evaluation.Elements {
		element := suiteGroup.NewElement()
		if err := element.UnmarshalBinary(data); err != nil {
			return err
		}

		elements[index] = element
	}

	r.Evaluation.Elements = elements

	return nil
}
