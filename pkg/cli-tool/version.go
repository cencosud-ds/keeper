package cli_tool

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version prints out the current app version
var version = &cobra.Command{
	Use:   "version",
	Short: "Shows app version",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("version 0.1.0")
	},
}
