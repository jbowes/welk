// Package sham provides sham implementations of external commands.
//
// Any functions contained in this package are less general-purpose commands,
// usually for a single script, to allow for more builtins.
// Perhaps in time they can be converted to builtins (python interperter in go anyone?), or replaced
// with an inclusion list, or a per-script plugin mechanism (or all 3).
package sham

import (
	"github.com/jbowes/welk/internal/install/builtin"
)

// Shams use the BuiltinFunc type, but are not expected to deal with their supplied args,
// as they are already matched.
var Sham = make(map[string]builtin.BuiltinFunc)
