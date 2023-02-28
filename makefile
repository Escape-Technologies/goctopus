mock = mockoon-cli start --data 
mock-p-name = --pname "goctopus-mock"
goctopus = go run cmd/goctopus/goctopus.go

goctopus:
	$(goctopus)

mock-latency: stop-mock
	$(mock) ./test/timeout/mock.json $(mock-p-name)
	
stop-mock:
	mockoon-cli stop mockoon-goctopus-mock

test-latency: mock-latency
	time $(goctopus) -i test/timeout/input.txt -v