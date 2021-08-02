package builtin

import (
	"context"
	"fmt"
	"io"
)

type Host interface {
	Log(tag string, message ...string)

	File(path string) []byte
	ChDir(path string)
	Write(path string) io.WriteCloser
	MkDir(path string)
	Remove(path string)
	Move(from, to string) error
}

type IOs struct {
	In  io.Reader
	Out io.Writer
}

type BuiltinFunc func(ctx context.Context, host Host, ios IOs, args []string) error

var Builtin = make(map[string]BuiltinFunc)

type ExitStatusError struct {
	status uint8
}

func NewExitStatusError(status uint8) error { return &ExitStatusError{status: status} }

func (e *ExitStatusError) Error() string { return fmt.Sprintf("exit code %d", e.status) }
func (e ExitStatusError) Status() uint8  { return e.status }

var _ error = &ExitStatusError{}
