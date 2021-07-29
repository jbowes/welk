package builtin

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/pflag"
)

func Tr(ctx context.Context, host Host, ios IOs, args []string) error {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	delete := fs.BoolP("", "d", false, "")
	err := fs.Parse(args)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if *delete && len(fs.Args()) != 1 {
		return errors.New("only one arg expected with -d")
	} else if !*delete && len(fs.Args()) != 2 {
		return errors.New("2 args required for tr")
	}

	// XXX: convert escape chars, not fully accurate
	s1 := fmt.Sprint(fs.Arg(0))
	s2 := fmt.Sprint(fs.Arg(1))

	// if s2 < s1, repeat last char for replacement.
	replace := map[string]string{}
	for i, ci := range s1 {
		// TODO: not utf8 safe
		var r string
		if len(s2) <= i {
			r = ""
		} else {
			r = string(s2[i])
		}
		replace[string(ci)] = r
	}

	host.Log("tr")

	buf, err := io.ReadAll(ios.In)
	if err != nil {
		return err
	}

	out := string(buf)
	for k, v := range replace {
		out = strings.ReplaceAll(out, k, v)
	}

	_, err = ios.Out.Write([]byte(out))

	return err
}

func init() { Builtin["tr"] = Tr }
