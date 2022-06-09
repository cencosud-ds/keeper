package cli_tool

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows app version",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("version 0.1.0")
	},
}
