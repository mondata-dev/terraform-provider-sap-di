default: testacc

# Run acceptance tests
.PHONY: testacc
testacc: start-mock-server
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Run mock server
.PHONY: start-mock-server
start-mock-server:
	cd mock-server && docker-compose up --build --detach

.PHONY: stop-mock-server
stop-mock-server:
	cd mock-server && docker-compose down
