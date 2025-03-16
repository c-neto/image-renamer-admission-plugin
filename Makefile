.PHONY: run test

run:
	go run ./cmd/image-renamer-admission-plugin

test:
	go test ./...

build-container:
	docker build -t image-renamer-admission-plugin .

run-container:
	docker run -p 8080:8080 image-renamer-admission-plugin
