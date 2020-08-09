.PHONY: download
download:
	@echo Download go.mod dependencies
	@go mod download

# https://marcofranssen.nl/manage-go-tools-via-go-modules/
.PHONY: install-tool
install-tools: download
	@echo Installing tools from scripts/tools.go
	@cat scripts/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

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

.PHONY: test
test:
	gotestsum --format testname ./...

.PHONY: cover
cover:
	go test -race -covermode atomic -coverprofile=coverage.out ./...
