/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fionna",
	Short: "A brief description of your application",
	Long: `
███████╗ ██╗  ██████╗  ███╗   ██╗ ███╗   ██╗  █████╗ 
██╔════╝ ██║ ██╔═══██╗ ████╗  ██║ ████╗  ██║ ██╔══██╗
█████╗   ██║ ██║   ██║ ██╔██╗ ██║ ██╔██╗ ██║ ███████║
██╔══╝   ██║ ██║   ██║ ██║╚██╗██║ ██║╚██╗██║ ██╔══██║
██║      ██║ ╚██████╔╝ ██║ ╚████║ ██║ ╚████║ ██║  ██║
╚═╝      ╚═╝  ╚═════╝  ╚═╝  ╚═══╝ ╚═╝  ╚═══╝ ╚═╝  ╚═╝
                Author: omsage org
                 License: AGPL v3
                                                `,
	RunE: func(cmd *cobra.Command, args []string) error {
		runWeb()
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
