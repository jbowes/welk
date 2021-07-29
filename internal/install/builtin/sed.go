package builtin

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/rwtodd/Go.Sed/sed"
	"github.com/spf13/pflag"
)

func Sed(ctx context.Context, host Host, ios IOs, args []string) error {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	quiet := fs.BoolP("", "n", false, "")
	err := fs.Parse(args)
	if err != nil {
		fmt.Println(err)
		return err
	}

	host.Log("sed")

	var eng *sed.Engine
	if *quiet {
		eng, err = sed.NewQuiet(strings.NewReader(strings.Join(fs.Args(), " ")))
	} else {
		eng, err = sed.New(strings.NewReader(strings.Join(fs.Args(), " ")))
	}

	if err != nil {
		return err
	}

	_, err = io.Copy(ios.Out, eng.Wrap(ios.In))
	return err
}

func init() { Builtin["sed"] = Sed }
