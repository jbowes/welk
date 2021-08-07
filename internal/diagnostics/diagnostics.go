package diagnostics

import "runtime/debug"

var (
	version = "(devel)"
	//buildDate = ""
	// builtBy   = ""
)

func Version() string {
	// If the version was previously set here, or set with -ldflags -X, leave it be.
	if version != "(devel)" {
		return version
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
		}
	}

	return version
}
