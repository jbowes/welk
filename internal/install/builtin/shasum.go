// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builtin

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"io"

	"github.com/spf13/pflag"
)

func Shasum(ctx context.Context, host Host, ios IOs, args []string) error {
	fs := pflag.NewFlagSet("", pflag.ContinueOnError)
	alg := fs.StringP("a", "a", "1", "")
	err := fs.Parse(args)
	if err != nil {
		return err
	}

	var h hash.Hash
	switch *alg {
	case "1":
		h = sha1.New()
	case "224":
		h = sha256.New224()
	case "256":
		h = sha256.New()
	case "384":
		h = sha512.New384()
	case "512":
		h = sha512.New()
	case "512224":
		h = sha512.New512_224()
	case "512256":
		h = sha512.New512_256()
	default:
		return NewExitStatusError(1)
	}

	host.Log("shasum", "-a", *alg)

	name := "-"
	var b []byte
	if len(fs.Args()) == 0 {
		b, err = io.ReadAll(ios.In)
		if err != nil {
			return err
		}
	} else {
		b = host.File(ctx, fs.Arg(0))
		name = fs.Arg(0)
	}

	if _, err := h.Write(b); err != nil {
		return err
	}

	s := h.Sum(nil)

	_, err = fmt.Fprintf(ios.Out, "%x  %s", s, name)
	return err
}

func init() { Builtin["shasum"] = Shasum }
