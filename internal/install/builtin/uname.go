package builtin

import (
	"context"
	"fmt"
	"runtime"
)

var unamed = map[string]string{
	"darwin": "Darwin",
}

func Uname(ctx context.Context, host Host, ios IOs, args []string) error {
	// TODO: uname could prompt for non m and non s args?

	host.Log("uname")

	if len(args) == 0 || args[0] == "-s" {
		output := unamed[runtime.GOOS]
		_, err := fmt.Fprint(ios.Out, output)
		return err
	}
	return nil

}

func init() { Builtin["uname"] = Uname }
