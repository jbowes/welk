// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package devnull

import "io"

type devNull struct{}

func (devNull) Read(p []byte) (int, error)  { return 0, io.EOF }
func (devNull) Write(p []byte) (int, error) { return len(p), nil }
func (devNull) Close() error                { return nil }

func New() io.ReadWriteCloser      { return devNull{} }
func IsDevNull(i interface{}) bool { _, ok := i.(devNull); return ok }
