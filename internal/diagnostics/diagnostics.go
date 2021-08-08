// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package diagnostics

import (
	"runtime"
	"runtime/debug"
)

var (
	version   = "(devel)"
	buildTime = "unknown"
	builtBy   = "unknown"
)

func loadBuildInfo() {
	// If the version was previously set here, or set with -ldflags -X, leave it be.
	if version != "(devel)" {
		return
	}

	// See if we can get a version from module build info.
	if info, ok := debug.ReadBuildInfo(); ok {
		mod := &info.Main
		if mod.Replace != nil {
			mod = mod.Replace
		}

		// TODO: normalize form between goreleaser, git describe (maybe), and go mod.
		// this is important for anything off-tag.
		version = mod.Version

		if version == "(devel)" {
			version = "unknown"
		} else if builtBy == "unknown" {
			builtBy = "go module"
		}
	}
}

type Diagnostics struct {
	Version   string
	BuildTime string
	BuiltBy   string

	Goos   string
	Goarch string
	// TODO: include libc (eg muscl)

	// TODO: Add Checks for path info

	// TODO: Add go module information
}

func New() *Diagnostics {
	loadBuildInfo()

	return &Diagnostics{
		Version:   version,
		BuildTime: buildTime,
		BuiltBy:   builtBy,

		Goos:   runtime.GOOS,
		Goarch: runtime.GOARCH,
	}
}
