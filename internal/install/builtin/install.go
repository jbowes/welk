package builtin

import (
	"context"

	"github.com/spf13/pflag"
)

func Install(ctx context.Context, host Host, ios IOs, args []string) error {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)

	d := fs.BoolP("d", "d", false, "")
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	host.Log("install", args...)

	if *d {
		for _, n := range fs.Args() {
			host.MkDir(n)
		}

		return nil
	}

	// TODO: copy files

	return nil
}

func init() { Builtin["install"] = Install }
