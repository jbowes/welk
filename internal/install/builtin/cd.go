package builtin

import (
	"context"
)

func Cd(ctx context.Context, host Host, ios IOs, args []string) error {
	host.Log("cd", args[0])
	host.ChDir(args[0])
	return nil
}

func init() { Builtin["sumdog-cd"] = Cd }
