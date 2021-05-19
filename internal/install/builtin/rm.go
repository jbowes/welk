package builtin

import (
	"context"
	"fmt"

	"github.com/spf13/pflag"
)

func Rm(ctx context.Context, host Host, ios IOs, args []string) error {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	_ = fs.BoolP("", "r", true, "")
	err := fs.Parse(args)
	if err != nil {
		fmt.Println(err)
		return err
	}

	host.Log("rm", fs.Arg(0))

	host.Remove(fs.Arg(0))
	return nil
}

func init() { Builtin["rm"] = Rm }
