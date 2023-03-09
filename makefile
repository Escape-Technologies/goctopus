mock = mockoon-cli start --data 
mock-p-name = --pname "goctopus-mock"
goctopus = go run cmd/goctopus/goctopus.go

goctopus:
	$(goctopus)

test:
	go test ./...

build:
	go build -o goctopus cmd/goctopus/goctopus.go

release-docker:
	docker buildx build --platform linux/amd64,linux/arm64 --push -t registry.gitlab.com/escape.tech/misc/goctopus .

.PHONY: test