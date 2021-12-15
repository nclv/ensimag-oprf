.PHONY: all clean clean-perfs clean-binary run build test-bench profile-bench
all: build

BINARY_DIR = ./bin
PROFILE_DIR = ./perfs

BENCHMARKS = BenchmarkClientBaseModeOPRFP256 BenchmarkClientVerifiableModeOPRFP256 \
	BenchmarkClientBaseModeOPRFP384 BenchmarkClientVerifiableModeOPRFP384 \
	BenchmarkClientBaseModeOPRFP521 BenchmarkClientVerifiableModeOPRF521

run:
	go run .

build:
	go build -o ${BINARY_DIR}/client .


test-bench:
	go test -bench .

profile-bench:
	go install golang.org/x/perf/cmd/benchstat@latest

	$(foreach bench,$(BENCHMARKS), \
		go test -bench=$(bench) -count 5 -benchmem \
			-memprofile ${PROFILE_DIR}/memprofile_$(bench).pprof \
			-cpuprofile ${PROFILE_DIR}/cpuprofile_$(bench).pprof \
			-trace=${PROFILE_DIR}/trace_$(bench).pprof \
			| tee ${PROFILE_DIR}/$(bench).txt; \
		benchstat ${PROFILE_DIR}/$(bench).txt \
	;)

clean : clean-binary clean-perfs

clean-perfs:
	rm ${PROFILE_DIR}/*.pprof ${PROFILE_DIR}/*.txt

clean-binary:
	go clean
	rm ${BINARY_DIR}/*

