/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// refreshCmd represents the refresh command
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Refresh access token",
	Long: `This command refreshes the access token for the application.
It reads the refresh token from the configuration file and sends a request to the authorization server to get a new access token.`,
	Run: func(_ *cobra.Command, args []string) {
		yellow := color.New(color.FgYellow).SprintFunc()
		fmt.Println(yellow("Token has expired. Refreshing token..."))
		if err := container.AuthToken.Refresh(container.Ctx, createOAuthConfig(), container.Config.Credential.Path); err != nil {
			red := color.New(color.FgRed).SprintFunc()
			fmt.Println(red(fmt.Sprintf("Failed to refresh token: %v", err)))
			fmt.Println("You may need to re-authenticate, if want to access the credential API.")
		} else {
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Println(green("Successfully refreshed token"))
		}
	},
}

func init() {
	rootCmd.AddCommand(refreshCmd)
}
