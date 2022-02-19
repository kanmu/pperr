.PHONY: test
test: vet
	go test ./...

.PHONY: vet
vet:
	go vet ./...

