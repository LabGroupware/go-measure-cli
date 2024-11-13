package config

import "time"

// ClockConfig represents the configuration for the clock
type ClockConfig struct {
	Format string         `mapstructure:"format"`
	Fake   FakeTimeConfig `mapstructure:"fake"`
}

// FakeTimeConfig represents the configuration for the fake time
type FakeTimeConfig struct {
	Enabled bool      `mapstructure:"enabled"`
	Time    time.Time `mapstructure:"time"`
}
