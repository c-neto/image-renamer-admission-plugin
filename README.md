WIP!!!

---

# Image Renamer Admission Controller

This project is an admission controller for Kubernetes that automatically prefixes container images with a specified registry and repository. It is implemented in Go and uses the Kubernetes API to mutate admission requests.

## Project Structure

```bash
.
├── cmd
│   └── image-renamer-admission-plugin
│       └── main.go
├── pkg
│   └── admission
│       ├── handler.go
│       └── handler_test.go
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

- [`go.mod`](go.mod): Go module file that defines the module path and dependencies.
- [`go.sum`](go.sum): Go checksum file that ensures the integrity of the dependencies.
- [`cmd/image-renamer-admission-plugin/main.go`](cmd/image-renamer-admission-plugin/main.go): Main application file that starts the server.
- [`pkg/admission/handler.go`](pkg/admission/handler.go): Contains the admission controller logic.
- [`pkg/admission/handler_test.go`](pkg/admission/handler_test.go): Test file that contains unit tests for the admission controller logic.
- [`Makefile`](Makefile): Makefile to run the application.

## Prerequisites

- Go 1.24.1 or later
- Kubernetes cluster

## Installation

1. Clone the repository:

```sh
git clone https://github.com/c-neto/image-renamer-admission-plugin.git
cd image-renamer-admission-plugin
```

2. Install dependencies:

```sh
go mod tidy
```

## Usage

To run the application, use the following command:

```sh
make run
```

This will start the server on port 8080.

## Endpoints

- `/mutate`: Admission controller endpoint that mutates the container images.
- `/healthz`: Health check endpoint that returns "ok".
- `/readyz`: Readiness check endpoint that returns "ready".

## Testing

To run the tests, use the following command:

```sh
go test ./...
```

## License

This project is licensed under the MIT License. See the LICENSE file for details.
