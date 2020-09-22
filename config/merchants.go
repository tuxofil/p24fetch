package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Merchants represents an array of merchants configurations
type Merchants struct {
	// List of configured merchants
	Entries []*Config `json:"merchants"`
	// Default values for merchants configurations
	Defaults Config `json:"defaults"`
}

// NewConfigs reads merchants configurations from JSON file.
func NewConfigs(path string) ([]*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	var merchants Merchants
	if err := json.Unmarshal(data, &merchants); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}
	for i, m := range merchants.Entries {
		m.SetDefaultsFrom(merchants.Defaults)
		if err := m.Validate(); err != nil {
			return nil, fmt.Errorf("merchant #%d validate: %w", i, err)
		}
	}
	return merchants.Entries, nil
}
