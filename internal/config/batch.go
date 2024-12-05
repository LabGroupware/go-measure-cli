package config

type BatchConfig struct {
	Test    BatchTestConfig `mapstructure:"test"`
	Metrics MetricsConfig   `mapstructure:"metrics"`
}

type BatchTestConfig struct {
	Path   string `mapstructure:"path"`
	Output string `mapstructure:"output"`
}

type MetricsConfig struct {
	Output string `mapstructure:"output"`
}
