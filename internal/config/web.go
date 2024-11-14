package config

// WebConfig represents the configuration for the web command
type WebConfig struct {
	WebSocket WebSocketConfig `mapstructure:"websocket"`
}

// WebSocketConfig represents the configuration for the WebSocket connection
type WebSocketConfig struct {
	Url string `mapstructure:"url"`
}
