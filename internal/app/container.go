package app

import (
	"fmt"
	"time"

	"github.com/LabGroupware/go-measure-tui/internal/app/i18n"
	"github.com/LabGroupware/go-measure-tui/internal/auth"
	"github.com/LabGroupware/go-measure-tui/internal/clock"
	"github.com/LabGroupware/go-measure-tui/internal/clock/fakeclock"
	"github.com/LabGroupware/go-measure-tui/internal/config"
	"github.com/LabGroupware/go-measure-tui/internal/logger"
)

// Container holds the dependencies for the application
type Container struct {
	Clocker    clock.Clock
	Translator i18n.Translation
	Config     config.Config
	Logger     logger.Logger
	AuthToken  *auth.AuthToken
}

// NewContainer creates a new Container
func NewContainer() *Container {
	return &Container{}
}

// Init initializes the Container
func (c *Container) Init(cfg config.Config) error {
	var err error

	// ----------------------------------------
	// Set Config
	// ----------------------------------------
	c.Config = cfg

	// ----------------------------------------
	// Set Default Language
	// ----------------------------------------
	switch c.Config.Lang {
	case "en":
		i18n.Default = i18n.English
	case "ja":
		i18n.Default = i18n.Japanese
	case "":
		fmt.Println("No language specified. Defaulting to English.")
		c.Config.Lang = "en"
		i18n.Default = i18n.English
	default:
		fmt.Println("Invalid language specified. Defaulting to English.")
		c.Config.Lang = "en"
		i18n.Default = i18n.English
	}

	// ----------------------------------------
	// Set Clock
	// ----------------------------------------
	if _, err = time.Parse(c.Config.Clock.Format, c.Config.Clock.Format); err != nil {
		fmt.Println("Invalid clock format. Defaulting to 2006-01-02 15:04:05.\n Error:", err)
		c.Config.Clock.Format = "2006-01-02 15:04:05"
	}

	clk := clock.New()
	if cfg.Clock.Fake.Enabled {
		fakeClk := fakeclock.New(cfg.Clock.Fake.Time)
		clk = fakeClk
	}
	c.Clocker = clk

	//----------------------------------------
	// Set Translator
	//----------------------------------------
	c.Translator, err = i18n.NewTranslator()
	if err != nil {
		return fmt.Errorf("failed to create translator: %w", err)
	}

	// ----------------------------------------
	// Set Logger
	// ----------------------------------------
	c.Logger = logger.NewSlogLogger()
	if err := c.Logger.SetupLogger(&cfg.Logging); err != nil {
		return fmt.Errorf("failed to setup logger: %w", err)
	}

	return nil
}

// Close closes the Container
func (c *Container) Close() error {
	c.Logger.Close()
	return nil
}
