// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import "github.com/spf13/cobra"

// TODO: have this run install by default
var rootCmd = cobra.Command{
	Use:   "welk",
	Short: "Manage and verify curl | sh style scripts",
}

// TODO: add a -v, --version flag (though v could be for verbose)

// TODO: handle err + panics here.
func Execute() { _ = rootCmd.Execute() }
