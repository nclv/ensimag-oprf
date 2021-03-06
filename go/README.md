# A simple pseudonymization service

This repository contains the implementation of a simple pseudonymization service based on the version 8 of the draft [_Oblivious Pseudorandom Functions (OPRFs) using Prime-Order Groups_](https://github.com/cfrg/draft-irtf-cfrg-voprf) from [cloudflare/circl](https://github.com/cloudflare/circl/tree/master/oprf).

---

From the draft:
> An Oblivious Pseudorandom Function (OPRF) is a two-party protocol between client and server for computing the output of a Pseudorandom Function (PRF). The server provides the PRF secret key, and the client provides the PRF input. At the end of the protocol, **the client learns the PRF output without learning anything about the PRF secret key, and the server learns neither the PRF input nor output**. A Partially-Oblivious PRF (POPRF) is an OPRF that allows client and server to provide public input to the PRF. OPRFs and POPRFs can also satisfy a notion of 'verifiability'. In this setting, clients can verify that the server used a specific private key during the execution of the protocol.

The choosen implementation of the draft supports the **base** and **verifiable** modes, and uses multiplicative blinding (see [multiplicative-vs-additive-blinding](https://github.com/bytemare/voprf#multiplicative-vs-additive-blinding)). The supported ciphersuites are (P-256, SHA-256), (P-384, SHA-384) and (P-512, SHA-512).

---

Some notes about the protocol :
- **If the public information is fixed the protocol is deterministic**. It is trivial to recover the pseudonyms corresponding to some input data. The public information needs to be shared between the client (Finalization) and the server (Evaluation).

---

We provide :
- a [public API](https://ensimag-oprf.vercel.app/) to perform the evaluation with our server's secret key,
- a website that performs the full protocol : [ensimag-oprf.vercel.app](https://ensimag-oprf.vercel.app/) or [ensimag-oprf-nclv.vercel.app](https://ensimag-oprf-nclv.vercel.app),
- a command-line tool to generate the private keys for the supported ciphersuites,
- a local client to test the full protocol with a specific server and selected mode, ciphersuite and public information,
- a local server to test the full protocol with a specific client and selected or random secret keys.

The code is in [Go](https://go.dev/dl/). Install the newest version to use the tools.

## Server

The server provides 2 endpoints :
- `/api/request_public_keys` to retrieve the server's public keys for each encryption suite (P256, P384 and P521),
- `/api/evaluate` evaluates an array of blinded element.

---

The [API documentation](https://app.swaggerhub.com/apis-docs/nclv/ensimag-oprf) is generated by [swagger](https://app.swaggerhub.com/apis/nclv/ensimag-oprf/). Use the [inspector](https://inspector.swagger.io/builder) to add new endpoints and [swagger generator](https://roger13.github.io/SwagDefGen/) to generate the JSON response.

---

We provide a `Makefile` :
```bash
# Build the local server executable in /bin
make build

# Launch the local server on http://localhost:1323
make run-server 

# Generate a new secret key for P256 cipher suite
make run-key-gen 

# Performs real-time load testing with ali (https://github.com/nakabonne/ali)
# 500 requests per second for 10 seconds on the /api/evaluate endpoint
# The JSON request is perfs/evaluate.json
make load-test

# Delete the binary and clean the go build cache
make clean
```

### Deployment

Our website is deployed on [Vercel](https://vercel.com/nclv/ensimag-oprf).

The client Go code is compiled to WebAssembly to be used in the browser.

We used some tools to evaluate our website :
- Performance tests with [PageSpeed](https://pagespeed.web.dev/report?url=https%3A%2F%2Fensimag-oprf.vercel.app%2F) (use `preload` and `defer` to load the `client.wasm` file without blocking the main thread),
- Security tests with [Security Headers](https://securityheaders.com/?q=https%3A%2F%2Fensimag-oprf.vercel.app) (configure the `vercel.json` file with the secure headers).

---

**Centralized vs. decentralized**

There are some issues when doing the finalization because of the [decentralized architecture of Vercel](https://vercel.com/docs/concepts/functions/conceptual-model). The secret/public key pair changes at each request so the evaluation returned by the server correspond to a different public key that the key queried from `/api/request_public_keys`. 

The public key is only used client-side for the verifiable mode finalization step. However, because of the serverless architecture of the deployed solution, the secret key changes at each request. A consequence is that **the outputs will always be non-deterministic**, event if the public information is fixed. Moreover, **we can't associate the input data (d) to the pseudonymized data (p)** for the resolution i.e. do Resolve(p) -> d.  
A solution would be to share the secret key as a common secret. The user could do the resolution providing the knowledge of its input data and the public information, mode and suite used when generating the pseudonymized data by sending an evaluation request to the server. Note that the server's secret key should be the same one used when doing the pseudonimization. The user would also be able to retrieve the list of pseudonyms associated to the same input data by providing the input data and the list of public information used to generate the pseudonyms. The resolution with pseudonym requires the user to store all the public informations and input data.

We can no longer execute the protocol with the verifiable mode because the public key queried at the beginning of the protocol does not correspond to the secret key used in the protocol.  
The public key is not used when generating a request or blinding the inputs, so we could send the public key with the response of `/api/evaluate`. We then update the verifiable client's public key before the finalization. This step is only needed in verifiable mode as no public key is needed for the finalization in base mode.

_Note that there are no issues with the keys if we use a centralized server_ except that it will no longer be possible to perform the resolution if the server's secret key is lost.

---

**Deterministic vs. non-deterministic**

If we always provide the same public information the protocol is deterministic. We decided to follow the draft recommendations :

> It is RECOMMENDED that this metadata be constructed with some type of higher-level domain separation to avoid cross protocol attacks or related issues. For example, protocols using this construction might ensure that the metadata uses a unique, prefix-free encoding. Any system which has multiple POPRF applications should distinguish client inputs to ensure the POPRF results are separate.

**We generate a random, 256-byte length public information at each evaluation request for non-deterministic results.** However you can still use the public API to pseudonymize your data with a selected public information.

---

_The following are notes about deployment on Vercel_.

See [vercel.json](https://vercel.com/docs/cli#project-configuration/), [CDN caching](https://vercel.com/docs/concepts/edge-network/caching).

Project settings on _vercel.com_ :
- Set the root directory to `go/server/`,
- Do not override any command. The output directory is by default `server/public/`. It contains the pages and static files.

`vercel.json` parameters :
- `"cleanUrls": true` : all HTML files and Serverless Functions will have their extension removed. When visiting a path that ends with the extension, a 308 response will redirect the client to the extensionless path. Similarly, a Serverless Function named `api/index.go` will be served when visiting `/api/index`. Visiting `/api/index.go` will redirect to `/api/index`.
- `"trailingSlash": false` : visiting a path that ends with a forward slash will respond with a 308 status code and redirect to the path without the trailing slash. For example, the `/api/` path will redirect to `/api`. 
- The response header of all the Serverless Functions in `/api` is `Content-Type : application/json`.
- The rewrite `"source": "/api/(.*)", "destination": "/api"` redirects all requests to `/api/*` to `/api` i.e. to `/api/index.go` Serverless Function. In this function we instantiate the router that serves all the `/api` endpoints.

### Launch the local server

Some use cases of the provided key generation tool and local server :

```bash
# Generate a new private key for P256 cipher suite
make run-key-gen 
>go run ./key-gen
>2021/12/23 16:58:49 3
>2021/12/23 16:58:49 AtzyGS8NoBjEjqbhwdGY/zWyqdFkJghyTttoIGq4UoM=

# You can choose another cipher suite
go run ./key-gen -suite 4
>2021/12/23 18:19:07 4
>2021/12/23 18:19:07 DryE+vL9Q8ciTNy6TvC5c7iXgOmzwkhqHpzAuPAXBRL8uNczSINCqt3crXNjncIW

# Launch the server with a pre-computed private key for the P256 and P384 cipher suites
# The other private keys (P521) are generated at the server creation.
P256_PRIVATE_KEY=AtzyGS8NoBjEjqbhwdGY/zWyqdFkJghyTttoIGq4UoM= P384_PRIVATE_KEY=DryE+vL9Q8ciTNy6TvC5c7iXgOmzwkhqHpzAuPAXBRL8uNczSINCqt3crXNjncIW make run-server

# Launch the server with new private keys
make run-server
```

### Using cURL

After launching the server, you can test the API endpoints :

```bash
# Get the static keys
curl -X GET http://localhost:1323/api/request_public_keys
["Aw8w56VYF4ejVfxCWt91AjPzdimuqqONpIkSrO74c4Ga","A9pxkw7jys6VmafHG1bhHOCd0b9nakuxZzHgQmDeiN8DtyemjeinyjtSNxdZPI50dQ==","AwF+WC+bWEBW1GT9wownSD7UokFge1BM7OMXAlzx9KgC4B+HMZxKgHN/FMXm9dmHaYUWXEDk4W13w2xwJGAbu1LmGw=="]

# Evaluate the blinded elements
curl -X POST http://localhost:1323/api/evaluate -H 'Content-Type: application/json' -d '{"suite": 3, "mode": 1, "info": "7465737420696e666f", "blinded_elements": [[2, 99, 233, 95, 211, 165, 194, 204, 118, 22, 17, 134, 162, 84, 135, 138, 180, 7, 229, 225, 238, 137, 138, 247, 196, 178, 119, 121, 218, 135, 36, 201, 132],[2, 61, 128, 127, 32, 157, 20, 86, 131, 22, 159, 225, 197, 38, 118, 154, 158, 71, 70, 50, 188, 116, 40, 80, 108, 72, 139, 91, 98, 146, 135, 105, 40]]}' # blinded elements of [][]byte{{0x00}, {0xFF}}
{"Elements":["AnzOnrnGUiaNurfXL3HXR9u7IQfQHMJ0T7alfEVn4339","A0jpFesUdIFhySiR2u9+FKAJSkGCrKyI7X8w7B2GurbA"],"Proof":null}

# evaluation request with base64 encoded blinded elements is NOT IMPLEMENTED
{"suite":3,"mode":1,"blinded_elements":["MTIzNA==","MjMz"]}  # Base64 encoded strings
```

### Load testing
`/api/evaluate` endpoint load testing with `ali` :

```bash
# Performs real-time load testing with ali (https://github.com/nakabonne/ali)
# 500 requests per second for 10 seconds on the /api/evaluate endpoint
# The JSON request is perfs/evaluate.json
make load-test
```

The server evaluation takes about 10ms for each request.

## Client

The client is composed of a CLI for command-line interaction with the server (`/cmd` directory). The WebAssembly binary used on the deployed website is generated from the code into `/wasm` to the `/server/public/static` directory.

---

We provide a `Makefile` :
```bash
# Call build-cmd and build-wasm
make build

# Build the local client executable in /bin
make build-cmd  

# Build the WASM binary in /server/public/static
make build-wasm  

# Launch the local client with a server on http://localhost:1323
make run

# Call clean-perfs and clean-binary
make clean

# Delete benchmarks and profiling results in /perfs
make clean-perfs

# Delete the binary and clean the go build cache
make clean-binary  

# Run all the client benchmarks
make test-bench  

# Save the profiling results (CPU, memory, trace) of each benchmark (ciphersuites x modes) in /perfs
make profile-bench
```

---

_The following are notes about WASM_.

The client frontend may lag if we spam the _Send_ button. The error `scheduleTimeoutEvent: missed timeout event` is thrown from `wasm_exec.js`. There is a bug on the second _Send_ after reloading the page. The WASM client doesn't seem to load. The fix for both issues is to load the wasm instance only once.

### Launch the client CLI

After launching the local server you can run the local client :
```bash
# Run the client with default mode and suite
make run
# Run the client with specific mode and suite on default input data [][]byte{{0x00}, {0xFF}}
go run ./cmd/ -mode=1 -suite=5
# Run the client on a list of inputs ["deadbeef", "one", "My name is"]
go run ./cmd/ -mode=1 -suite=4 deadbeef one "My name is"
```

### Benchmarks

We provide six benchmarks of the full protocol (blinding, random information generation, evaluation and finalization) for each mode (base and verifiable) with each ciphersuite (P-256, P-384 and P-512).

```bash
# Run all benchmarks
make test-bench
# Generate the memory profile, cpu profile and traces for each benchmarks in perfs/
make profile-bench

# Analyse one of the generated profiles
go tool pprof -http=:8080 <profile>.pprof
```

Looking at the CPU profiles for the base and verifiable modes, we can see that the verifiable mode contributes to a 50% overhead on the client side during the finalization.
