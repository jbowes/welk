package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/jbowes/welk/internal/diagnostics"
	"github.com/spf13/cobra"
)

var versionAll bool

var versionCmd = cobra.Command{
	Use:   "version",
	Short: "print version information",
	Run: func(cmd *cobra.Command, args []string) {
		d := diagnostics.New()

		if !versionAll {
			fmt.Println(d.Version)
			return
		}

		// TODO: colorize because that would be fun.
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
		fmt.Fprintf(w, "%s:\t%s\n", "version", d.Version)
		fmt.Fprintf(w, "%s:\t%s\n", "build date", d.BuildDate)
		fmt.Fprintf(w, "%s:\t%s\n", "built by", d.BuiltBy)
		w.Flush()
		fmt.Fprintf(w, "%s:\t%s\n", "GOOS", d.Goos)
		fmt.Fprintf(w, "%s:\t%s\n", "GOARCH", d.Goarch)
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(&versionCmd)
	// TODO: --all isn't a great name but it works for now.
	versionCmd.Flags().BoolVar(&versionAll, "all", false, "show additional version information")
}
