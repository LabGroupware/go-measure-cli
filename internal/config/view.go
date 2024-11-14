package config

// ViewConfig represents the view configuration
type ViewConfig struct {
	Color ColorViewConfig `mapstructure:"color"`
	Theme string          `mapstructure:"theme"`
}

// ColorViewConfig represents the color view configuration
type ColorViewConfig struct {
	Primary   string `mapstructure:"primary"`
	Secondary string `mapstructure:"secondary"`
	Success   string `mapstructure:"success"`
	Danger    string `mapstructure:"danger"`
	Warning   string `mapstructure:"warning"`
	Info      string `mapstructure:"info"`
	Light     string `mapstructure:"light"`
	Dark      string `mapstructure:"dark"`
}
