/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/LabGroupware/go-measure-tui/internal/batch/batchtest"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Batch test",
	Long:  `Batch test is a command to test the batch command.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := batchtest.BatchTest(container); err != nil {
			container.Logger.Error(cmd.Context(), "Batch test failed", logger.Value("error", err))
		}
	},
}

func init() {
	batchCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
