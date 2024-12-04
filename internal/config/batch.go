package config

type BatchConfig struct {
	Test BatchTestConfig `mapstructure:"test"`
}

type BatchTestConfig struct {
	Path      string               `mapstructure:"path"`
	MassQuery BatchMassQueryConfig `mapstructure:"massquery"`
}

type BatchMassQueryConfig struct {
	Output string `mapstructure:"output"`
}
