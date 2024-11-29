// Package utils provides utility functions for the application
package utils

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/rivo/tview"
)

// SetFocus sets the focus to the specified primitive
func SetFocus(app *tview.Application, primitive tview.Primitive) {
	app.SetFocus(primitive)
}

// GenerateUniqueID generates a unique ID
func GenerateUniqueID() string {
	return uuid.New().String()
}

// Contains checks if a slice contains a specific element
func Contains[T comparable](slice []T, elem T) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}

// AnyContains checks if any of the slices contains a specific element
func AnyContains[T comparable](slices1 []T, slices2 []T) bool {
	for _, e1 := range slices1 {
		for _, e2 := range slices2 {
			if e1 == e2 {
				return true
			}
		}
	}
	return false
}

// AllContains checks if all of the slices contain a specific element
func AllContains[T comparable](slices1 []T, slices2 []T) bool {
	for _, e2 := range slices2 {
		found := false
		for _, e1 := range slices1 {
			if e1 == e2 {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// UnmarshalJSON unmarshals JSON data into a struct
func UnmarshalJSON(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
