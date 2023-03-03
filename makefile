mock = mockoon-cli start --data 
mock-p-name = --pname "goctopus-mock"
goctopus = go run cmd/goctopus/goctopus.go

goctopus:
	$(goctopus)

test:
	go test ./...

build:
	go build -o goctopus cmd/goctopus/goctopus.go

# Change this to dockerhub when going public
release-docker:
	docker build -t registry.gitlab.com/escape.tech/misc/goctopus .
	docker push registry.gitlab.com/escape.tech/misc/goctopus

# LEGACY
mock-latency: stop-mock
	$(mock) ./test/timeout/mock.json $(mock-p-name)
	
stop-mock:
	mockoon-cli stop mockoon-goctopus-mock

test-latency: mock-latency
	time $(goctopus) -i test/timeout/input.txt -v


.PHONY: test