package builtin

import (
	"context"
)

func Chmod(ctx context.Context, host Host, ios IOs, args []string) error {
	host.Log("chmod", args[1])
	// TODO: update store
	return nil
}

func init() { Builtin["chmod"] = Chmod }
