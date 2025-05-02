package goversion

import (
	"runtime/debug"
)

func Version () string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "1.22" // fallback version
	}

	return info.GoVersion[2:]
}