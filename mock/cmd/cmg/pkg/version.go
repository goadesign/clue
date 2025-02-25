package cluemockgen

import (
	"runtime/debug"
)

// Version returns the complete version number.
func Version() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "(unknown)"
	}
	return bi.Main.Version
}
