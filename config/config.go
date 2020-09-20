package config

import (
	"errors"
	"fmt"

	"github.com/kelseyhightower/envconfig"

	"github.com/tuxofil/p24fetch/schema"
)

type Config struct {
	// Descriptive name of the merchant.
	// Used for logging/messaging.
	MerchantName string `split_words:"true" required:"true"`
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
	// Export format
	ExportFormat schema.Format `split_words:"true" required:"true"`
	// Mandatory for QIF export format.
	// Source account name.
	SrcAccountName string `split_words:"true"`
	// Mandatory for QIF export format.
	// Account name for comissions.
	ComissionAccountName string `split_words:"true"`
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
	if c.MerchantName == "" {
		return errors.New("no merchant name")
	}
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
	switch c.ExportFormat {
	case schema.JSON:
	case schema.QIF:
		if c.SrcAccountName == "" {
			return errors.New("no source account name")
		}
		if c.ComissionAccountName == "" {
			return errors.New("no comission account name")
		}
	default:
		return fmt.Errorf("invalid export format: %s", c.ExportFormat)
	}
	return nil
}
