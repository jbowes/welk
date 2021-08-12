// Copyright (c) 2021 James Bowes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/lipgloss"
	"github.com/jbowes/welk/internal/diagnostics"
	"github.com/jbowes/whatsnew"
	"github.com/spf13/cobra"
)

// TODO: have this run install by default
var rootCmd = cobra.Command{
	Use:   "welk",
	Short: "Manage and verify curl | sh style scripts",
}

// TODO: add a -v, --version flag (though v could be for verbose)

// TODO: handle err + panics here.
func Execute() {
	// TODO: make this as middleware and don't put it on every command.
	d := diagnostics.New()
	fut := whatsnew.Check(context.TODO(), &whatsnew.Options{
		Slug:    "jbowes/welk",
		Cache:   filepath.Join(xdg.CacheHome, "welk", "release-check.json"),
		Version: d.Version,
	})

	_ = rootCmd.Execute()

	if v, _ := fut.Get(); v != "" {

		bold := lipgloss.NewStyle().Bold(true) // TODO: this doesn't work with non-terminals :(
		s := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF9999"))

		fmt.Println(s.Render("💾 new release available:"), v)
		fmt.Println("download from", bold.Render(
			fmt.Sprintf("https://github.com/jbowes/welk/releases/tag/%s", v),
		))
	}
}
