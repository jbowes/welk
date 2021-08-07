// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builtin

import (
	"context"
	"fmt"

	"github.com/spf13/pflag"
)

func Mkdir(ctx context.Context, host Host, ios IOs, args []string) error {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	_ = fs.BoolP("", "p", true, "")
	err := fs.Parse(args)
	if err != nil {
		fmt.Println(err)
		return err
	}

	host.Log("mkdir", fs.Arg(0))

	host.MkDir(ctx, fs.Arg(0))
	return nil
}

func init() { Builtin["mkdir"] = Mkdir }
