package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/tuxofil/p24fetch/config"
	"github.com/tuxofil/p24fetch/dedup"
	"github.com/tuxofil/p24fetch/exporter"
	"github.com/tuxofil/p24fetch/merchant"
	"github.com/tuxofil/p24fetch/schema"
	"github.com/tuxofil/p24fetch/slack"
	"github.com/tuxofil/p24fetch/sorter"
)

// Entry point
func main() {
	log.Println("started")
	if err := Main(); err != nil {
		log.Fatalf("%s", err)
	}
	log.Println("done")
}

func Main() error {
	configs, err := config.NewConfigs(os.Args[1])
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}
	for _, config := range configs {
		if err := processMerchant(config); err != nil {
			return fmt.Errorf("%s: %w", config.MerchantName, err)
		}
		time.Sleep(10 * time.Second)
	}
	return nil
}

func processMerchant(cfg *config.Config) error {
	log.Printf("processing: %s", cfg.MerchantName)
	merchant, err := merchant.New(cfg)
	if err != nil {
		return fmt.Errorf("create merchant: %w", err)
	}
	dedup, err := dedup.New(cfg)
	if err != nil {
		return fmt.Errorf("create deduplicator: %w", err)
	}
	sorter, err := sorter.New(cfg)
	if err != nil {
		return fmt.Errorf("create sorter: %w", err)
	}
	slack, err := slack.New(cfg)
	if err != nil {
		return fmt.Errorf("create Slack interface: %w", err)
	}

	ctx := context.TODO()
	// Fetch transaction log
	xmlTrans, err := merchant.FetchLog(ctx)
	if err != nil {
		return fmt.Errorf("fetch log: %w", err)
	} else if len(xmlTrans) == 0 {
		log.Printf("no transactions found")
		return nil
	}

	// Deduplicate
	newTrans := dedup.Filter(xmlTrans)
	var lastTran schema.XMLTransaction
	if len(newTrans) > 0 {
		lastTran = newTrans[len(newTrans)-1]
	} else {
		log.Printf("fetched %d transactions but no new found", len(xmlTrans))
		return nil
	}

	// Parse transactions
	trans := make([]schema.Transaction, len(newTrans))
	for i, tran := range newTrans {
		trans[i] = schema.ParseTransaction(tran)
	}

	// Sort transactions
	ignoredTrans, sortedTrans, unsortedTrans := sorter.Sort(trans)
	log.Printf("  sorted: %d; unsorted: %d; ignored: %d",
		len(sortedTrans), len(unsortedTrans), len(ignoredTrans))

	// Export sorted transactions
	if err := exporter.New(cfg).Export(sortedTrans); err != nil {
		return fmt.Errorf("export sorted: %w", err)
	}

	// Export ignored transaction as JSON
	cfg.ExportFormat = schema.JSON
	resultsDir := cfg.ResultsDir
	cfg.ResultsDir = path.Join(resultsDir, "ignored")
	if err := exporter.New(cfg).Export(ignoredTrans); err != nil {
		return fmt.Errorf("export ignored: %w", err)
	}

	// Export unsorted transactions as JSON
	cfg.ResultsDir = path.Join(resultsDir, "unsorted")
	if err := exporter.New(cfg).Export(unsortedTrans); err != nil {
		return fmt.Errorf("export unsorted: %w", err)
	}

	// Send Slack notifications for unsorted transactions
	slack.ReportUnsorted(unsortedTrans)

	// Update deduplicator state
	if err := dedup.Update(lastTran); err != nil {
		return fmt.Errorf("update dedup: %w", err)
	}
	return nil
}
