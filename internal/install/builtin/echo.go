package builtin

import (
	"context"
	"fmt"
	"strings"

	"github.com/jbowes/welk/internal/install/devnull"
)

func Echo(ctx context.Context, host Host, ios IOs, args []string) error {

	// TODO: can we tell when this is a pipe for command stuff? fmt.Println(ios.Out)

	// TODO: only log output when ios.Out goes nowhere (set this up with interp.Exec)
	if devnull.IsDevNull(ios.Out) {
		host.Log("echo", strings.Join(args, " "))
	} else {
		fmt.Fprint(ios.Out, strings.Join(args, " "))
		host.Log("echo")
	}

	return nil
}

func init() { Builtin["welk-echo"] = Echo }
