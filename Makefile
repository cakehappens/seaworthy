.PHONY: gen
install-deps:
	go install github.com/gotestyourself/gotestsum

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

.PHONY: clean
clean:
	rm -f ./seaworthy

.PHONY: clean
test:
	gotestsum --format short-verbose ./...

