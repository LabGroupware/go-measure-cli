/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/LabGroupware/go-measure-tui/internal/testprompt"
	"github.com/LabGroupware/go-measure-tui/internal/testprompt/massquery"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "For testing purposes",
	Long:  `This command is used for testing purposes. It does not have any functionality.`,
	Run: func(cmd *cobra.Command, args []string) {
		testType, err := testprompt.TestTypeSelection()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch testType {
		case testprompt.TestPromptMassiveQuery:
			massquery.MassiveQueryPrompt()
		case testprompt.TestPromptWaitSaga:
			fmt.Println("Wait Saga")
		case testprompt.TestPromptUrgeOnConsistentAfterStartSaga:
			fmt.Println("Urge on Consistent After Start Saga")
		case testprompt.TestPromptUrgeOnConsistentAfterEndSaga:
			fmt.Println("Urge on Consistent After End Saga")
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
