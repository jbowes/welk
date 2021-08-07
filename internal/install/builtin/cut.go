// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builtin

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/pflag"
)

func Cut(ctx context.Context, host Host, ios IOs, args []string) error {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	delim := fs.StringP("", "d", "\t", "")
	field := fs.IntP("f", "f", 1, "") // TODO: multi field support
	err := fs.Parse(args)
	if err != nil {
		fmt.Println(err)
		return err
	}

	host.Log("cut")

	buf, err := io.ReadAll(ios.In)
	if err != nil {
		return err
	}

	parts := strings.Split(string(buf), *delim)
	_, err = ios.Out.Write([]byte(parts[*field-1]))
	return err
}

func init() { Builtin["cut"] = Cut }
