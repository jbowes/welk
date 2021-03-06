// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"path/filepath"
	"runtime/debug"

	"github.com/adrg/xdg"
	"github.com/jbowes/welk/internal/db"
	"github.com/spf13/cobra"
)

var listCmd = cobra.Command{
	Use:     "list",
	Short:   "list installed and known curl | sh style packages",
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			fmt.Printf("Failed to read build info")
			return
		}

		fmt.Println("I am", bi.Main)

		d := db.DB{Root: filepath.Join(xdg.DataHome, "welk", "installed")}
		err := d.List(func(m *db.Manifest) error {
			fmt.Println(m.URL, m.State)
			return nil
		})
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(&listCmd)
}
