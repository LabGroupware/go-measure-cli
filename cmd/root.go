/*
Package cmd provides the command line interface for the application
Copyright © 2024 NAME HERE k.hayashi@cresplanex.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/LabGroupware/go-measure-tui/internal/app"
	"github.com/LabGroupware/go-measure-tui/internal/config"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var container = app.NewContainer()

var ctx context.Context

// var ctx context.Context

var rootCmd = &cobra.Command{
	Use:   "Go Measure TUI",
	Short: "A simple TUI for measuring web performance",
	Long: `For measuring web performance, Go Measure TUI is a simple terminal user interface that
allows you to measure the performance of a website by providing the URL and the number of requests
to be made to the website. It then displays the average time taken to make the requests.`,
}

// Execute runs the root command
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("failed to execute root command: %w", err)
	}

	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.nova-measure/config.yaml)")
	if err := viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config")); err != nil {
		fmt.Printf("Error binding flag: %v\n", err)
		os.Exit(1)
	}
}

func initConfig() {
	ctx = context.Background()

	configFile := viper.GetString("config")
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Printf("Failed to get home directory: %v\n", err)
			os.Exit(1)
		}
		viper.AddConfigPath(homeDir + "/.nova-measure")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// 環境変数も読み込む
	viper.AutomaticEnv()
	// 環境変数のプレフィックスを設定
	viper.SetEnvPrefix("NOVA_MEASURE")
	// 環境変数名のセパレーターを変換（例: "NOVA_MEASURE_LOGGING_LEVEL" -> "logging.level"）
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 設定ファイルを読み込む
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		os.Exit(1)
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg, func(m *mapstructure.DecoderConfig) {
		m.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToIPNetHookFunc(),
			mapstructure.StringToIPHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		)
	}); err != nil {
		fmt.Printf("Error unmarshalling config: %v\n", err)
		os.Exit(1)
	}

	if err := container.Init(cfg); err != nil {
		fmt.Printf("Error initializing container: %v\n", err)
		os.Exit(1)
	}

	if expire := container.AuthToken.IsExpired(); expire {
		fmt.Println("Token has expired. Refreshing token...")
		if err := container.AuthToken.Refresh(ctx, createOAuthConfig(), container.Config.Credential.Path); err != nil {
			fmt.Printf("Failed to refresh token: %v\n", err)
			fmt.Println("You may need to re-authenticate, if want to access the credential API.")
		} else {
			fmt.Println("Token refreshed.")
		}
	}

	// container.Logger.Info(ctx, "Container initialized",
	// 	logger.Value("LANG", container.Config.Lang),
	// 	logger.Group("CLOCK",
	// 		logger.Value("FORMAT", container.Config.Clock.Format),
	// 		logger.Group("FAKE",
	// 			logger.Value("TIME", container.Config.Clock.Fake.Time),
	// 			logger.Value("ENABLE", container.Config.Clock.Fake.Enabled),
	// 		),
	// 	),
	// )
}
