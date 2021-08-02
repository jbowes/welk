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
	squeeze := fs.BoolP("s", "s", false, "")
	delete := fs.BoolP("d", "d", false, "")
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	// TODO: -d should support 2 args
	if *delete && len(fs.Args()) != 1 {
		return errors.New("only one arg expected with -d")
	} else if !*delete && len(fs.Args()) != 2 {
		return errors.New("2 args required for tr")
	}

	// XXX: convert escape chars, not fully accurate
	// TODO: this is totally wrong
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

	// TODO: hacky special case
	if !*delete && s1 == "[:upper:]" && s2 == "[:lower:]" {
		out = strings.ToLower(out)
	} else {

		for k, v := range replace {
			out = strings.ReplaceAll(out, k, v)
		}
	}

	if *squeeze {
		smap := replace // just rely on the keys
		if len(fs.Args()) == 2 {
			smap = make(map[string]string)
			for _, r := range s2 {
				smap[string(r)] = ""
			}
		}

		var last rune
		b := &strings.Builder{}
		for _, c := range out {
			if _, ok := smap[string(c)]; !ok || c != last {
				b.WriteRune(c)
				last = c
			} // skip
		}

		out = b.String()
	}

	_, err = ios.Out.Write([]byte(out))

	return err
}

func init() { Builtin["tr"] = Tr }
