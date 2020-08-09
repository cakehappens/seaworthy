.PHONY: download
dl:
	@echo Download go.mod dependencies
	@go mod download

# https://marcofranssen.nl/manage-go-tools-via-go-modules/
.PHONY: install-tool
tools: dl
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
	mkdir -p ./dist
	go build -o ./dist/seaworthy ./cmd/seaworthy
	chmod +x ./dist/seaworthy

.PHONY: clean
clean:
	rm -fr ./dist

.PHONY: test
test:
	gotestsum --format testname ./...

.PHONY: cover
cover:
	go test -race -covermode atomic -coverprofile coverage.out ./...

# PRE-COMMIT & GITHOOKS
# ---------------------
pre-commit.install:
	pre-commit install --install-hooks

pre-commit.run:
	pre-commit run --all-files

release: clean
	@echo "--skip-publish, as we will use github actions to do this"
	git-chglog -o CHANGELOG.md
	goreleaser --snapshot --skip-publish --rm-dist --release-notes CHANGELOG.md
