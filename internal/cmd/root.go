package cmd

import "github.com/spf13/cobra"

// TODO: have this run install by default
var rootCmd = cobra.Command{
	Use:   "sumdog",
	Short: "Manage and verify curl | sh style scripts",
}

func Execute() error { return rootCmd.Execute() }
