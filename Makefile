.PHONY: gen
gen:
	go generate ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: build
build: clean
	go build -o ./seaworthy ./cmd/seaworthy
	chmod +x ./seaworthy

clean:
	rm -f ./seaworthy
