// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builtin

import (
	"context"
)

func Ln(ctx context.Context, host Host, ios IOs, args []string) error {
	host.Log("ln")

	// TODO: actually implement this
	return nil
}

func init() { Builtin["ln"] = Ln }
