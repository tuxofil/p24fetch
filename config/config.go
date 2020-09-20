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
}

// Read and parse configurations.
func New() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("parse env vars: %w", err)
	}
	if cfg.MerchantID < 1 {
		return nil, fmt.Errorf("invalid merchant ID: %d", cfg.MerchantID)
	}
	if cfg.MerchantPassword == "" {
		return nil, errors.New("invalid merchant password")
	}
	if cfg.CardNumber == "" {
		return nil, errors.New("invalid card number")
	}
	return cfg, nil
}
