/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Go Measure TUI",
	Long:  `All software has versions. This is Go Measure TUI's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Go Measure TUI v0.0.1 -- HEAD")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
