package core

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"testing"

	"github.com/cloudflare/circl/oprf"
)

func setup(mode oprf.Mode, suite oprf.Suite) *Client {
	return NewClient("http://localhost:1323/api", suite, mode)
}

func exchange(client *Client, mode oprf.Mode, suite oprf.Suite) [][]byte {
	finalizeData, oprfEvaluationRequest, _ := client.Blind([][]byte{[]byte("dead3eef")})

	token := make([]byte, 256)
	if _, err := rand.Read(token); err != nil {
		log.Println(err)

		return nil
	}

	info := hex.EncodeToString(token)
	// log.Println("Public information : ", info)

	blindedElements, err := SerializeElements(oprfEvaluationRequest.Elements)
	if err != nil {
		log.Fatal(err)
	}

	evaluationRequest := NewEvaluationRequest(
		suite, mode, info, blindedElements,
	)

	evaluation, err := client.EvaluateRequest(evaluationRequest)
	if err != nil {
		log.Fatal(err)
	}

	outputs, err := client.Finalize(finalizeData, evaluation.Evaluation, info)
	if err != nil {
		log.Fatal(err)
	}

	return outputs
}

func benchmarkClient(b *testing.B, mode oprf.Mode, suite oprf.Suite) {
	b.Helper()

	client := setup(mode, suite)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = exchange(client, mode, suite)
	}
}

func BenchmarkClientBaseModeOPRFP256(b *testing.B) {
	benchmarkClient(b, oprf.BaseMode, oprf.SuiteP256)
}

func BenchmarkClientVerifiableModeOPRFP256(b *testing.B) {
	benchmarkClient(b, oprf.VerifiableMode, oprf.SuiteP256)
}

func BenchmarkClientPartiallyObliviousModeOPRFP256(b *testing.B) {
	benchmarkClient(b, oprf.PartialObliviousMode, oprf.SuiteP256)
}

func BenchmarkClientBaseModeOPRFP384(b *testing.B) {
	benchmarkClient(b, oprf.BaseMode, oprf.SuiteP384)
}

func BenchmarkClientVerifiableModeOPRFP384(b *testing.B) {
	benchmarkClient(b, oprf.VerifiableMode, oprf.SuiteP384)
}

func BenchmarkClientPartiallyObliviousModeOPRFP384(b *testing.B) {
	benchmarkClient(b, oprf.PartialObliviousMode, oprf.SuiteP384)
}

func BenchmarkClientBaseModeOPRFP521(b *testing.B) {
	benchmarkClient(b, oprf.BaseMode, oprf.SuiteP521)
}

func BenchmarkClientVerifiableModeOPRF521(b *testing.B) {
	benchmarkClient(b, oprf.VerifiableMode, oprf.SuiteP521)
}

func BenchmarkClientPartiallyObliviousModeOPRFP521(b *testing.B) {
	benchmarkClient(b, oprf.PartialObliviousMode, oprf.SuiteP521)
}
