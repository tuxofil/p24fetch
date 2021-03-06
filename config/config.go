package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tuxofil/p24fetch/schema"
)

type Config struct {
	// Descriptive name of the merchant.
	// Used for logging/messaging.
	MerchantName string `json:"merchant_name"`

	// Privat24 Merchant ID
	MerchantID int `json:"merchant_id"`
	// Privat34 Merchant Password
	MerchantPassword string `json:"merchant_password"`
	// Bank card number
	CardNumber string `json:"card_number"`
	// Fetch transaction history for this number of days
	Days int `json:"days"`

	// Deduplicator state directory
	DedupDir string `json:"dedup_dir"`
	// Path to a JSON file with sorting rules.
	RulesPath string `json:"rules_path"`
	// Path to a directory to write exported files
	ResultsDir string `json:"results_dir"`

	// Export format
	ExportFormat schema.Format `json:"export_format"`
	// Mandatory for QIF export format.
	// Source account name -- GnuCash Account ID.
	SrcAccountName string `json:"src_account_name"`
	// Mandatory for QIF export format.
	// Account name for comissions -- GnuCash Account ID.
	ComissionAccountName string `json:"comission_account_name"`

	// Token used to authenticate to Slack API
	SlackToken string `json:"slack_token"`
	// Slack channel ID to write messages to.
	SlackChannel string `json:"slack_channel"`

	// Logging interface. Optional.
	Logger *log.Logger
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
	if c.ResultsDir == "" {
		c.ResultsDir = d.ResultsDir
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
	if c.ResultsDir == "" {
		return fmt.Errorf("invalid results dir: %#v", c.ResultsDir)
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
