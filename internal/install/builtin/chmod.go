// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builtin

import (
	"context"
)

func Chmod(ctx context.Context, host Host, ios IOs, args []string) error {
	host.Log("chmod", args[1])
	// TODO: update vfs
	return nil
}

func init() { Builtin["chmod"] = Chmod }
