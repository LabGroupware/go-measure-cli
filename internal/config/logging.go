package config

// LoggingOutputConfig represents the configuration for logging output
type LoggingOutputConfig struct {
	Type     string `mapstructure:"type"`
	Format   string `mapstructure:"format"`
	Level    string `mapstructure:"level"`
	Filename string `mapstructure:"filename"`
	Address  string `mapstructure:"address"`
}

// LoggingConfig represents the configuration for logging
type LoggingConfig struct {
	Output []LoggingOutputConfig `mapstructure:"output"`
}
