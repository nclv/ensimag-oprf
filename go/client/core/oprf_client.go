package core

import "github.com/cloudflare/circl/oprf"

type OprfClientInterface interface {
	Blind(inputs [][]byte) (*oprf.FinalizeData, *oprf.EvaluationRequest, error)
	DeterministicBlind(inputs [][]byte, blinds []oprf.Blind) (*oprf.FinalizeData, *oprf.EvaluationRequest, error)
	Finalize(f *oprf.FinalizeData, e *oprf.Evaluation, info []byte) (outputs [][]byte, err error)
}

type OprfBaseClient struct {
	oprf.Client
}

func (c OprfBaseClient) Finalize(f *oprf.FinalizeData, e *oprf.Evaluation, info []byte) ([][]byte, error) {
	return c.Client.Finalize(f, e)
}

type OprfVerifiableClient struct {
	oprf.VerifiableClient
}

func (c OprfVerifiableClient) Finalize(f *oprf.FinalizeData, e *oprf.Evaluation, info []byte) ([][]byte, error) {
	return c.VerifiableClient.Finalize(f, e)
}

type OprfPartialObliviousClient struct {
	oprf.PartialObliviousClient
}

// NewOPRFClient create an OPRF client with the provided suite, mode and static key.
// No static key is needed for oprf.BaseMode.
func NewOPRFClient(suite oprf.Suite, mode oprf.Mode, pkS *oprf.PublicKey) OprfClientInterface {
	var client OprfClientInterface

	switch mode {
	case oprf.BaseMode:
		client = OprfBaseClient{oprf.NewClient(suite)}
	case oprf.VerifiableMode:
		client = OprfVerifiableClient{oprf.NewVerifiableClient(suite, pkS)}
	case oprf.PartialObliviousMode:
		client = OprfPartialObliviousClient{oprf.NewPartialObliviousClient(suite, pkS)}
	}

	return client
}
