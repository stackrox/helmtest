// +build tools

package toolimports

import (
	// Tool imports to make `go install` possible with go modules.
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)
