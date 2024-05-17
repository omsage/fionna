package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version code of sas",
	Long:  "Version code of sas",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("FIONNA_VERSION")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
