package builtin

import (
	"context"
	"strings"
)

func Echo(ctx context.Context, host Host, ios IOs, args []string) error {

	// TODO: can we tell when this is a pipe for command stuff? fmt.Println(ios.Out)

	host.Log("echo", strings.Join(args, " "))

	return nil
}

func init() { Builtin["sumdog-echo"] = Echo }
