package cmd

import (
	"fmt"

	"github.com/jbowes/welk/internal/diagnostics"
	"github.com/spf13/cobra"
)

var versionCmd = cobra.Command{
	Use:   "version",
	Short: "print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(diagnostics.Version())
	},
}

func init() {
	rootCmd.AddCommand(&versionCmd)
}
