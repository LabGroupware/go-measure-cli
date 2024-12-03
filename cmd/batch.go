/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// batchCmd represents the batch command
var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "Batch operations",
	Long: `Batch operations are used to perform multiple operations at once.
For example, you can use the batch command to upload multiple files at once.`,
}

func init() {
	rootCmd.AddCommand(batchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// batchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// batchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
