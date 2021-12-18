package main

import (
	"encoding/json"
	"errors"
	"github.com/cloudflare/circl/oprf"
)

// WrappedPseudonimizeRequest allows to partially parse the JSON input
type WrappedPseudonimizeRequest struct {
	Data  []json.RawMessage `json:"data"`
	Mode  oprf.Mode         `json:"mode"`
	Suite oprf.SuiteID      `json:"suite"`
}

// PseudonimizeRequest holds the data and the client setup parameters.
type PseudonimizeRequest struct {
	Data  [][]byte     `json:"data"`
	Mode  oprf.Mode    `json:"mode"`
	Suite oprf.SuiteID `json:"suite"`
}

// ValidateMode validates the client mode (Base or Verifiable).
func (p *PseudonimizeRequest) ValidateMode() error {
	switch p.Mode {
	case oprf.BaseMode, oprf.VerifiableMode:
		return nil
	}

	return errors.New("invalid mode")
}

// ValidateSuite validates the encryption suite (P256, P384, P521).
func (p *PseudonimizeRequest) ValidateSuite() error {
	switch p.Suite {
	case oprf.OPRFP256, oprf.OPRFP384, oprf.OPRFP521:
		return nil
	}

	return errors.New("invalid suite")
}

// UnmarshalJSON first parse the JSON input into a WrappedPseudonimizeRequest
// and then convert the string into a byte array into PseudonimizeRequest.
func (p *PseudonimizeRequest) UnmarshalJSON(data []byte) error {
	var wp WrappedPseudonimizeRequest
	if err := json.Unmarshal(data, &wp); err != nil {
		return err
	}

	p.Data = make([][]byte, len(wp.Data))

	var input string
	for index, stringInput := range wp.Data {
		if err := json.Unmarshal(stringInput, &input); err != nil {
			return err
		}
		p.Data[index] = []byte(input)
	}

	p.Suite = wp.Suite
	p.Mode = wp.Mode

	return nil
}
