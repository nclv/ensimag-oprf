.PHONY: all clean run-server run-key-gen build load-test
all: build

BINARY_DIR = ./bin
PROFILE_DIR = ./perfs

run-key-gen:
	go run ./key-gen

run-server:
	go run ./cmd

build:
	go build -o ${BINARY_DIR}/server ./cmd

load-test:
	ali --body-file=${PROFILE_DIR}/evaluate.json \
		--rate=500 --duration=10s \
		-H 'Content-Type: application/json' --method=POST http://127.0.0.1:1323/api/evaluate

clean:
	go clean
	rm ${BINARY_DIR}/*

