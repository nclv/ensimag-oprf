## Server

The server provides 2 endpoints :
- `/api/request_public_keys` to retrieve the server's public keys for each encryption suite (P256, P384, P521),
- `/api/evaluate` evaluates an array of blinded element.

The information must be private to avoid non-deterministic results. It needs to be shared between the client (Finalization) and the server (Evaluation).

```bash
# Makefile commands
make {build,run,load-test,clean}
```

### Deployment

On [Vercel](https://vercel.com/nclv/ensimag-oprf).

See [the vercel.json configuration file](https://vercel.com/docs/cli#project-configuration/), [CDN caching](https://vercel.com/docs/concepts/edge-network/caching).

**Centralized vs. decentralized**

Issue when doing the finalization because of the [architecture](https://vercel.com/docs/concepts/functions/conceptual-model). The private/public key pair changes at each request so the evaluation returned by the server correspond to a different public key that the key queried from `/api/request_public_keys`. The public key is not used when generating a request or blinding the inputs, so we could send the public key with the response of `/api/evaluate`. We then update the verifiable client's public key before the finalization. This step is only needed in verifiable mode as no public key is needed for the finalization in base mode.

The public key is only used client-side for the verifiable mode finalization step. However, because of the serverless architecture of the deployed solution, the private key changes at each request. A consequence is that the outputs will always be non-deterministic, event if the public information is fixed. Moreover, we can't associate the input data (d) to the pseudonymized data (p) for the resolution i.e. do Resolve(p) -> d.
A solution would be to share the private key as a common secret. See [storing complex secrets](https://github.com/vercel/vercel/issues/749). The user could do the resolution providing the knowledge of its input data and the public information, mode and suite used when generating the pseudonymized data by sending an evaluation request to the server. Note that the server's private key should be the same one used when doing the pseudonimization. The user would also be able to retrieve the list of pseudonyms associated to the same input data by providing the input data and the list of public information used to generate the pseudonyms. The resolution with pseudonym requires to store all the public informations.

If the public information is fixed the protocol is deterministic. It is trivial to recover the pseudonyms corresponding to some input data.

---

There are no issues with the keys if we use a centralized server.

---

Project settings on _vercel.com_ :
- Set the root directory to `go/server/`,
- Do not override any command. The output directory is by default `server/public`. It contains the pages and static files.

`vercel.json` parameters :
- `"cleanUrls": true` : all HTML files and Serverless Functions will have their extension removed. When visiting a path that ends with the extension, a 308 response will redirect the client to the extensionless path. Similarly, a Serverless Function named `api/index.go` will be served when visiting `/api/index`. Visiting `/api/index.go` will redirect to `/api/index`.
- `"trailingSlash": false` : visiting a path that ends with a forward slash will respond with a 308 status code and redirect to the path without the trailing slash. For example, the `/api/` path will redirect to `/api`. 
- The response header of all the Serverless Functions in `/api` is `Content-Type : application/json`.
- The rewrite `"source": "/api/(.*)", "destination": "/api"` redirects all requests to `/api/*` to `/api` i.e. to `/api/index.go` Serverless Function. In this function we instantiate the router that serves all the `/api` endpoints.

### Launch the server

```bash
make run
```

### Using cURL

```bash
# Get the static keys
curl -X GET http://localhost:1323/api/request_public_keys
["Aw8w56VYF4ejVfxCWt91AjPzdimuqqONpIkSrO74c4Ga","A9pxkw7jys6VmafHG1bhHOCd0b9nakuxZzHgQmDeiN8DtyemjeinyjtSNxdZPI50dQ==","AwF+WC+bWEBW1GT9wownSD7UokFge1BM7OMXAlzx9KgC4B+HMZxKgHN/FMXm9dmHaYUWXEDk4W13w2xwJGAbu1LmGw=="]

# Evaluate the blinded elements
curl -X POST http://localhost:1323/api/evaluate -H 'Content-Type: application/json' -d '{"suite": 3, "mode": 1, "info": "7465737420696e666f", "blinded_elements": [[2, 99, 233, 95, 211, 165, 194, 204, 118, 22, 17, 134, 162, 84, 135, 138, 180, 7, 229, 225, 238, 137, 138, 247, 196, 178, 119, 121, 218, 135, 36, 201, 132],[2, 61, 128, 127, 32, 157, 20, 86, 131, 22, 159, 225, 197, 38, 118, 154, 158, 71, 70, 50, 188, 116, 40, 80, 108, 72, 139, 91, 98, 146, 135, 105, 40]]}' # blinded elements of [][]byte{{0x00}, {0xFF}}

{"Elements":["AnzOnrnGUiaNurfXL3HXR9u7IQfQHMJ0T7alfEVn4339","A0jpFesUdIFhySiR2u9+FKAJSkGCrKyI7X8w7B2GurbA"],"Proof":null}

{"suite":3,"mode":1,"blinded_elements":["MTIzNA==","MjMz"]}  # Base64 encoded strings
```

### Load testing
`/api/evaluate` endpoint load testing with `ali` :

```bash
make load-test
```

## Client

The client is composed of a CLI for command-line interaction with the server (`/cmd` directory). The WASM binary is generated from the code into `/wasm` to the `/server/public/static` directory.

---

The client frontend may lag if we spam the _Send_ button. The error `scheduleTimeoutEvent: missed timeout event` is thrown from `wasm_exec.js`. There is a bug on the second _Send_ after reloading the page. The WASM client doesn't seem to load. The fix for both issues is to load the wasm instance only once.

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