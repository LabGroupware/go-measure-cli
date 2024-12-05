package config

type BatchConfig struct {
	Test    BatchTestConfig `mapstructure:"test"`
	Metrics MetricsConfig   `mapstructure:"metrics"`
}

type BatchTestConfig struct {
	Path      string               `mapstructure:"path"`
	MassQuery BatchMassQueryConfig `mapstructure:"massquery"`
}

type BatchMassQueryConfig struct {
	Output string `mapstructure:"output"`
}

type MetricsConfig struct {
	Output string `mapstructure:"output"`
}
