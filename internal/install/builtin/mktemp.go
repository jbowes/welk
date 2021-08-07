// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builtin

import (
	"context"
)

func Mktemp(ctx context.Context, host Host, ios IOs, args []string) error {

	// TODO: better (real) pattern
	fname := "/temp/faketmp/...."
	host.Log("mktemp", fname)

	host.MkDir(ctx, fname)

	_, err := ios.Out.Write([]byte(fname))
	return err
}

func init() { Builtin["mktemp"] = Mktemp }
