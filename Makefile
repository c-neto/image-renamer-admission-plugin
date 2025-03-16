.PHONY: run test

run:
	go run ./cmd/image-renamer-admission-controller

test:
	go test ./...
