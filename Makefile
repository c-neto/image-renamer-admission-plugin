.PHONY: run test

run:
	go run ./cmd/image-renamer-admission-plugin

test:
	go test ./...
