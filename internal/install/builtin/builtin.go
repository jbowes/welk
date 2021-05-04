package builtin

import (
	"context"
	"io"
)

type Host interface {
	Log(tag string, message ...string)

	ChDir(path string)
	Write(path string) io.WriteCloser
	MkDir(path string)
}

type IOs struct {
	Out io.Writer
}

type BuiltinFunc func(ctx context.Context, host Host, ios IOs, args []string) error

var Builtin = make(map[string]BuiltinFunc)
