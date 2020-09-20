package config

import (
	"errors"
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// Privat24 Merchant ID
	MerchantID int `split_words:"true" required:"true"`
	// Privat34 Merchant Password
	MerchantPassword string `split_words:"true" required:"true"`
	// Bank card number
	CardNumber string `split_words:"true" required:"true"`
	// Fetch transaction history for this number of days
	Days int `split_words:"true" required:"true"`
	// Deduplicator state directory
	DedupDir string `split_words:"true" required:"true"`
}

// Read and parse configurations.
func New() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("parse env vars: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	return cfg, nil
}

// Validate checks values of the configuration.
func (c *Config) Validate() error {
	if c.MerchantID < 1 {
		return fmt.Errorf("invalid merchant ID: %d", c.MerchantID)
	}
	if c.MerchantPassword == "" {
		return errors.New("invalid merchant password")
	}
	if c.CardNumber == "" {
		return errors.New("invalid card number")
	}
	if c.DedupDir == "" {
		return fmt.Errorf("invalid deduplicator dir: %#v", c.DedupDir)
	}
	if c.Days < 1 {
		return fmt.Errorf("invalid days number: %d", c.Days)
	}
	return nil
}
