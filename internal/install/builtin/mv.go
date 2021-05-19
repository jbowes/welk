package builtin

import (
	"context"
)

func Mv(ctx context.Context, host Host, ios IOs, args []string) error {
	host.Log("mv", args[0], args[1])
	// TODO: update store
	return nil
}

func init() { Builtin["mv"] = Mv }
