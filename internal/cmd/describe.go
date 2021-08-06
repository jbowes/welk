package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/jbowes/welk/internal/db"
	"github.com/spf13/cobra"
)

var describeCmd = cobra.Command{
	Use:   "describe a package",
	Short: "describe installed and known curl | sh style packages",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		d := db.DB{Root: filepath.Join(xdg.DataHome, "welk", "installed")}
		m, err := d.Query(args[0])
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("URL: %s\n", m.URL)
		fmt.Printf("State: %s\n", m.State)
		fmt.Println("Files:")
		for _, f := range m.Files {
			fmt.Printf("  %s\n", f.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(&describeCmd)
}
