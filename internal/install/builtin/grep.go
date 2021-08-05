package builtin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/spf13/pflag"
)

func Grep(ctx context.Context, host Host, ios IOs, args []string) error {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	err := fs.Parse(args)
	if err != nil {
		fmt.Println(err)
		return err
	}

	host.Log("grep", fs.Arg(0))

	re, err := regexp.Compile(fs.Arg(0))
	if err != nil {
		return err
	}

	var buf []byte
	if len(fs.Args()) == 1 {
		buf, err = io.ReadAll(ios.In)
		if err != nil {
			return err
		}
	} else {
		buf = host.File(ctx, fs.Arg(1))
	}

	found := false
	parts := strings.Split(string(buf), "\n")
	for _, p := range parts {
		if re.MatchString(p) {
			found = true
			ios.Out.Write([]byte(p))
			ios.Out.Write([]byte("\n"))
		}
	}

	if !found {
		return errors.New("no match") // 1 exit code?
	}
	return nil
}

func init() { Builtin["grep"] = Grep }
