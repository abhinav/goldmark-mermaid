//go:build tools
// +build tools

package tools

import (
	_ "github.com/mgechev/revive"
	_ "golang.org/x/tools/cmd/stringer"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
