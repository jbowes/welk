package builtin

import (
	"context"
)

func Mktemp(ctx context.Context, host Host, ios IOs, args []string) error {

	// TODO: better (real) pattern
	fname := "/temp/faketmp/...."
	host.Log("mktemp", fname)

	host.MkDir(fname)

	_, err := ios.Out.Write([]byte(fname))
	return err
}

func init() { Builtin["mktemp"] = Mktemp }
