package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"

	"github.com/tuxofil/p24fetch/schema"
)

type Config struct {
	// Descriptive name of the merchant.
	// Used for logging/messaging.
	MerchantName string `split_words:"true" json:"merchant_name"`

	// Privat24 Merchant ID
	MerchantID int `split_words:"true" json:"merchant_id"`
	// Privat34 Merchant Password
	MerchantPassword string `split_words:"true" json:"merchant_password"`
	// Bank card number
	CardNumber string `split_words:"true" json:"card_number"`
	// Fetch transaction history for this number of days
	Days int `split_words:"true" json:"days"`

	// Deduplicator state directory
	DedupDir string `split_words:"true" json:"dedup_dir"`
	// Path to a JSON file with sorting rules.
	RulesPath string `split_words:"true" json:"rules_path"`

	// Export format
	ExportFormat schema.Format `split_words:"true" json:"export_format"`
	// Mandatory for QIF export format.
	// Source account name.
	SrcAccountName string `split_words:"true" json:"src_account_name"`
	// Mandatory for QIF export format.
	// Account name for comissions.
	ComissionAccountName string `split_words:"true" json:"comission_account_name"`

	// Token used to authenticate to Slack API
	SlackToken string `split_words:"true" json:"slack_token"`
	// Slack channel ID to write messages to.
	SlackChannel string `split_words:"true" json:"slack_channel"`

	// Logging interface. Optional.
	Logger *log.Logger
}

// Read and parse configurations.
func New() (*Config, error) {
	cfg := &Config{Logger: log.New(os.Stderr, "", 0)}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("parse env vars: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	return cfg, nil
}

// SetDefaultsFrom copies missing values from another Config instance
func (c *Config) SetDefaultsFrom(d Config) {
	if c.Days == 0 {
		c.Days = d.Days
	}
	if c.DedupDir == "" {
		c.DedupDir = d.DedupDir
	}
	if c.RulesPath == "" {
		c.RulesPath = d.RulesPath
	}
	if c.ExportFormat == schema.Format("") {
		c.ExportFormat = d.ExportFormat
	}
	if c.SrcAccountName == "" {
		c.SrcAccountName = d.SrcAccountName
	}
	if c.ComissionAccountName == "" {
		c.ComissionAccountName = d.ComissionAccountName
	}
	if c.SlackToken == "" {
		c.SlackToken = d.SlackToken
	}
	if c.SlackChannel == "" {
		c.SlackChannel = d.SlackChannel
	}
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
	if fd, err := os.Open(c.RulesPath); err != nil {
		return fmt.Errorf("invalid path to rules file: %w", err)
	} else {
		_ = fd.Close()
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

	if s := c.SlackToken; s != "" && !strings.HasPrefix(s, "xoxp-") {
		return errors.New("invalid Slack token")
	}
	return nil
}

func (c *Config) Logf(format string, v ...interface{}) {
	if c.Logger != nil {
		c.Logger.Printf(format, v...)
	}
}
