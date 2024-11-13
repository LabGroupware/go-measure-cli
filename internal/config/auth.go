package config

// AuthConfig represents the configuration for the authentication
type AuthConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectHost string `mapstructure:"redirect_host"`
	RedirectPort string `mapstructure:"redirect_port"`
	RedirectPath string `mapstructure:"redirect_path"`
	AuthURL      string `mapstructure:"auth_url"`
	TokenURL     string `mapstructure:"token_url"`
}
