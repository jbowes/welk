package builtin

import (
	"context"
	"fmt"
	"strings"
)

func Echo(ctx context.Context, host Host, ios IOs, args []string) error {

	fmt.Println(ios.Out)

	host.Log("echo", strings.Join(args, " "))

	return nil
}

func init() { Builtin["sumdog-echo"] = Echo }
