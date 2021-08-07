package cmd

import "github.com/spf13/cobra"

// TODO: have this run install by default
var rootCmd = cobra.Command{
	Use:   "welk",
	Short: "Manage and verify curl | sh style scripts",
}

// TODO: handle err + panics here.
func Execute() { _ = rootCmd.Execute() }
