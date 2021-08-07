// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builtin

import (
	"context"
)

func Mv(ctx context.Context, host Host, ios IOs, args []string) error {
	host.Log("mv", args[0], args[1])
	return host.Move(ctx, args[0], args[1])
}

func init() { Builtin["mv"] = Mv }
