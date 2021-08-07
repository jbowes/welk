// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builtin

import (
	"context"
	"io"
	"strings"

	"github.com/rwtodd/Go.Sed/sed"
	"github.com/spf13/pflag"
)

func Sed(ctx context.Context, host Host, ios IOs, args []string) error {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	quiet := fs.BoolP("", "n", false, "")
	cmds := fs.StringArrayP("e", "e", nil, "")
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	host.Log("sed")

	var cmd string
	if len(*cmds) == 0 {
		cmd = strings.Join(fs.Args(), " ")
	} else {
		cmd = (*cmds)[0]
	}

	var eng *sed.Engine
	if *quiet {
		eng, err = sed.NewQuiet(strings.NewReader(cmd))
	} else {
		eng, err = sed.New(strings.NewReader(cmd))
	}

	if err != nil {
		return err
	}

	_, err = io.Copy(ios.Out, eng.Wrap(ios.In))
	return err
}

func init() { Builtin["sed"] = Sed }
