package builtin

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/pflag"
)

func Command(ctx context.Context, host Host, ios IOs, args []string) error {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	v := fs.BoolP("v", "v", false, "")
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	host.Log("command")

	if !*v {
		// TODO: implement this.
		return errors.New("command: requires -v")
	}

	if _, ok := Builtin[fs.Arg(0)]; ok {
		ios.Out.Write([]byte(fmt.Sprintf("/usr/bin/%s", fs.Arg(0))))
	} else {
		return errors.New("couldn't find command")
	}

	return nil
}

func init() { Builtin["sumdog-command"] = Command }