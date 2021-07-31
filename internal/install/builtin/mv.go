package builtin

import (
	"context"
)

func Mv(ctx context.Context, host Host, ios IOs, args []string) error {
	host.Log("mv", args[0], args[1])
	return host.Move(args[0], args[1])
}

func init() { Builtin["mv"] = Mv }
