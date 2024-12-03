package config

type BatchConfig struct {
	Test BatchTestConfig `mapstructure:"test"`
}

type BatchTestConfig struct {
	Path string `mapstructure:"path"`
}
