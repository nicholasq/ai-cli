package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
    Use:   "version",
    Short: "Print the version number of ai CLI",
    Long:  `All software has versions. This is ai's`,
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("ai CLI v0.1")
    },
}

func init() {
    rootCmd.AddCommand(versionCmd)
}
