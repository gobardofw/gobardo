package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// VersionCommand get cli version
var VersionCommand = &cobra.Command{
	Use:   "version",
	Short: "get gobardo cli version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version 0.0.1")
	},
}
