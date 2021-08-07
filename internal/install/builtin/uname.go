// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builtin

import (
	"context"
	"fmt"
	"runtime"

	"github.com/spf13/pflag"
)

var unamedOS = map[string]string{
	"darwin": "Darwin",
}

var unamedArch = map[string]string{
	"amd64": "x86_64",
}

func Uname(ctx context.Context, host Host, ios IOs, args []string) error {
	// TODO: uname could prompt for non m and non s args?
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	arch := fs.BoolP("m", "m", false, "")
	os := fs.BoolP("s", "s", false, "")
	err := fs.Parse(args)
	if err != nil {
		fmt.Println(err)
		return err
	}

	host.Log("uname")

	// TODO: args can be supplied at once and have a specific order for printing
	switch {
	case len(args) == 0, *os:
		output := unamedOS[runtime.GOOS]
		_, err := fmt.Fprint(ios.Out, output)
		return err
	case *arch:
		output := unamedArch[runtime.GOARCH]
		_, err := fmt.Fprint(ios.Out, output)
		return err
	}

	return nil
}

func init() { Builtin["uname"] = Uname }
