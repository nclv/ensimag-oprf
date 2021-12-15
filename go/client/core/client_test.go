package core

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/oprf/go/common"
	"log"
	"testing"

	"github.com/cloudflare/circl/oprf"
)

func setup(mode oprf.Mode, suite oprf.SuiteID) *Client {
	client := NewClient("http://localhost:1323")
	client.SetupOPRFClient(suite, mode)

	return client
}

func exchange(client *Client, mode oprf.Mode, suite oprf.SuiteID) [][]byte {
	clientRequest := client.CreateRequest([][]byte{[]byte("dead3eef")})

	token := make([]byte, 256)
	if _, err := rand.Read(token); err != nil {
		log.Println(err)

		return nil
	}

	info := hex.EncodeToString(token)
	// log.Println("Public information : ", info)

	evaluationRequest := common.NewEvaluationRequest(
		suite, mode, info, clientRequest.BlindedElements(),
	)
	evaluation := client.EvaluateRequest(evaluationRequest)

	return client.Finalize(clientRequest, evaluation, info)
}

func benchmarkClient(b *testing.B, mode oprf.Mode, suite oprf.SuiteID) {
	b.Helper()

	client := setup(mode, suite)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = exchange(client, mode, suite)
	}
}

func BenchmarkClientBaseModeOPRFP256(b *testing.B) {
	benchmarkClient(b, oprf.BaseMode, oprf.OPRFP256)
}

func BenchmarkClientVerifiableModeOPRFP256(b *testing.B) {
	benchmarkClient(b, oprf.VerifiableMode, oprf.OPRFP256)
}

func BenchmarkClientBaseModeOPRFP384(b *testing.B) {
	benchmarkClient(b, oprf.BaseMode, oprf.OPRFP384)
}

func BenchmarkClientVerifiableModeOPRFP384(b *testing.B) {
	benchmarkClient(b, oprf.VerifiableMode, oprf.OPRFP384)
}

func BenchmarkClientBaseModeOPRFP521(b *testing.B) {
	benchmarkClient(b, oprf.BaseMode, oprf.OPRFP521)
}

func BenchmarkClientVerifiableModeOPRF521(b *testing.B) {
	benchmarkClient(b, oprf.VerifiableMode, oprf.OPRFP521)
}
