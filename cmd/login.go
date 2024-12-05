/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/LabGroupware/go-measure-tui/internal/auth"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

func createOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     container.Config.Auth.ClientID,
		ClientSecret: container.Config.Auth.ClientSecret,
		RedirectURL: fmt.Sprintf("%s:%s%s",
			container.Config.Auth.RedirectHost,
			container.Config.Auth.RedirectPort,
			container.Config.Auth.RedirectPath),
		Scopes: []string{"openid", "profile", "email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   container.Config.Auth.AuthURL,
			TokenURL:  container.Config.Auth.TokenURL,
			AuthStyle: oauth2.AuthStyleAutoDetect,
		},
	}
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the application",
	Long:  `Responds to the login command. This command is used to login to the application.`,
	Run: func(cmd *cobra.Command, args []string) {
		oauthConf := createOAuthConfig()
		if authToken, err := auth.StartOAuthFlow(
			container.Ctx,
			*oauthConf,
			container.Config.Auth.RedirectPort,
			container.Config.Auth.RedirectPath,
			container.Config.Credential.Path,
		); err != nil {
			red := color.New(color.FgRed).SprintFunc()
			fmt.Println(red(fmt.Sprintf("Failed to start OAuth flow: %v", err)))
			os.Exit(1)
		} else {
			green := color.New(color.FgGreen).SprintFunc()
			fmt.Println(green("Successfully authenticated"))
			container.AuthToken = authToken
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
