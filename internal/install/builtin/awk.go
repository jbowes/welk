package builtin

import (
	"context"
	"fmt"
	"strings"

	"github.com/benhoyt/goawk/interp"
	"github.com/spf13/pflag"
)

func Awk(ctx context.Context, host Host, ios IOs, args []string) error {
	// TODO: support awk -v arg parsing.
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	sep := fs.StringP("", "F", " ", "")
	err := fs.Parse(args)
	if err != nil {
		fmt.Println(err)
		return err
	}

	host.Log("awk", fs.Args()...)

	return interp.Exec(strings.Join(fs.Args(), " "), *sep, ios.In, ios.Out)
}

func init() { Builtin["awk"] = Awk }
