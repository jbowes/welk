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
			host.MkDir(ctx, n)
		}

		return nil
	}

	if len(fs.Args()) == 1 {
		return NewExitStatusError(1)
	}

	if len(fs.Args()) == 2 {
		b := host.File(ctx, fs.Arg(0))

		host.Move(ctx, fs.Arg(0), fs.Arg(1))

		// TODO: awkward, but works with move semantics to dir.

		o := host.Write(ctx, fs.Arg(0))
		defer o.Close()

		o.Write(b)
	}

	// TODO: copy many files to directory

	return nil
}

func init() { Builtin["install"] = Install }
