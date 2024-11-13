/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display the configuration settings",
	Long: `This command displays the configuration settings that are currently in use by the application.
It displays the settings that are read from the configuration file and the environment variables.`,
	Run: func(_ *cobra.Command, args []string) {
		settings := viper.AllSettings()
		exConfigSettings := map[string]any{}
		for key, value := range settings {
			if key != "config" {
				exConfigSettings[key] = value
			}
		}
		yamlData, err := yaml.Marshal(exConfigSettings)
		if err != nil {
			log.Fatalf("Error converting configuration to YAML: %v", err)
		}
		fmt.Println("Current configuration:")
		fmt.Println(string(yamlData))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
