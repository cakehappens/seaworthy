// +build tools

package main

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/goreleaser/goreleaser"
	_ "gotest.tools/gotestsum"
	_ "github.com/git-chglog/git-chglog/cmd/git-chglog"
)
