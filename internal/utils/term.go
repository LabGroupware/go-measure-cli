package utils

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// GetTerminalSize returns the width and height of the terminal.
func GetTerminalSize() (int, int, error) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get terminal size: %w", err)
	}

	return width, height, nil
}
