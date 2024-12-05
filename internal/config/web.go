package config

// WebConfig represents the configuration for the web command
type WebConfig struct {
	QueryAPI  APIConfig       `mapstructure:"queryApi"`
	WebSocket WebSocketConfig `mapstructure:"websocket"`
}

// WebSocketConfig represents the configuration for the WebSocket connection
type WebSocketConfig struct {
	Url string `mapstructure:"url"`
}

// APIConfig represents the configuration for the API command
type APIConfig struct {
	Url string `mapstructure:"url"`
}
