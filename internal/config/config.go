// Package config provides configuration for the application.
package config

// Config represents the application configuration
type Config struct {
	Logging    LoggingConfig    `mapstructure:"logging"`
	Lang       string           `mapstructure:"lang"`
	Clock      ClockConfig      `mapstructure:"clock"`
	Auth       AuthConfig       `mapstructure:"auth"`
	Credential CredentialConfig `mapstructure:"credential"`
	Web        WebConfig        `mapstructure:"web"`
	View       ViewConfig       `mapstructure:"view"`
	Batch      BatchConfig      `mapstructure:"batch"`
}
