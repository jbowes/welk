package builtin

import (
	"context"
	"strings"

	"github.com/benhoyt/goawk/interp"
)

func Awk(ctx context.Context, host Host, ios IOs, args []string) error {
	// TODO: support awk argument parsing.

	return interp.Exec(strings.Join(args, " "), " ", ios.In, ios.Out)
}

func init() { Builtin["awk"] = Awk }
