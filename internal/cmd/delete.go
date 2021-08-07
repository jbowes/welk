// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/jbowes/welk/internal/db"
	"github.com/jbowes/welk/internal/filesync"
	"github.com/spf13/cobra"
)

var deleteCmd = cobra.Command{
	Use:   "delete URL",
	Short: "delete a previously installed curl | sh style package",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		d := db.DB{Root: filepath.Join(xdg.DataHome, "welk", "installed")}
		m, err := d.Query(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("deleting", args[0])

		txn, err := d.Delete(m)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer txn.Rollback()

		if err := filesync.Remove(m.Files); err != nil {
			fmt.Println(err)
			return
		}

		if err := txn.Commit(); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(&deleteCmd)
}
