## Server

The server provides 2 endpoints :
- `/request_public_keys` to retrieve the server's public keys for each encryption suite (P256, P384, P521),
- `/evaluate` evaluates an array of blinded element.

The information must be private to avoid non-deterministic results. It needs to be shared between the client (Finalization) and the server (Evaluation).

```bash
# Makefile commands
make {build,run,load-test,clean}
```

### Launch the server

```bash
make run
```

### Using cURL

```bash
# Get the public keys
curl -X GET http://localhost:1323/request_public_keys
["Aw8w56VYF4ejVfxCWt91AjPzdimuqqONpIkSrO74c4Ga","A9pxkw7jys6VmafHG1bhHOCd0b9nakuxZzHgQmDeiN8DtyemjeinyjtSNxdZPI50dQ==","AwF+WC+bWEBW1GT9wownSD7UokFge1BM7OMXAlzx9KgC4B+HMZxKgHN/FMXm9dmHaYUWXEDk4W13w2xwJGAbu1LmGw=="]

# Evaluate the blinded elements
curl -X POST http://localhost:1323/evaluate -H 'Content-Type: application/json' -d '{"suite": 3, "mode": 1, "info": "7465737420696e666f", "blinded_elements": [[2, 99, 233, 95, 211, 165, 194, 204, 118, 22, 17, 134, 162, 84, 135, 138, 180, 7, 229, 225, 238, 137, 138, 247, 196, 178, 119, 121, 218, 135, 36, 201, 132],[2, 61, 128, 127, 32, 157, 20, 86, 131, 22, 159, 225, 197, 38, 118, 154, 158, 71, 70, 50, 188, 116, 40, 80, 108, 72, 139, 91, 98, 146, 135, 105, 40]]}' # blinded elements of [][]byte{{0x00}, {0xFF}}

{"Elements":["AnzOnrnGUiaNurfXL3HXR9u7IQfQHMJ0T7alfEVn4339","A0jpFesUdIFhySiR2u9+FKAJSkGCrKyI7X8w7B2GurbA"],"Proof":null}

{"suite":3,"mode":1,"blinded_elements":["MTIzNA==","MjMz"]}  # Base64 encoded strings
```

### Load testing
`/evaluate` endpoint load testing with `ali` :

```bash
make load-test
```

## Client

The client 

```bash
# Makefile commands
make {build build-cmd build-wasm run clean clean-perfs clean-binary test-bench profile-bench}
```

### Compile the .wasm binary into the server public/ directory

```bash
make build-wasm
```

### Launch the client CLI

```bash
# Run the test with default mode and suite
make run
# Run the client with specific mode and suite on default input data [][]byte{{0x00}, {0xFF}}
go run ./cmd/ -mode=1 -suite=5
# Run the client on a list of inputs ["deadbeef", "one", "My name is"]
go run ./cmd/ -mode=1 -suite=4 deadbeef one "My name is"
```

### Benchmarks

```bash
# Run all benchmarks
make test-bench
# Generate the memory profile, cpu profile and traces for each benchmarks in perfs/
make profile-bench
# Analyse one of the generated profiles
go tool pprof -http=:8080 <profile>.pprof
```