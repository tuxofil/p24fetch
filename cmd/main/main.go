package main

import (
	"context"
	"fmt"
	"log"
	"os"

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
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("create config: %w", err)
	}
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
	exporter, err := exporter.New(cfg)
	if err != nil {
		return fmt.Errorf("create exporter: %w", err)
	}

	ctx := context.TODO()
	// Fetch transaction log
	xmlTrans, err := merchant.FetchLog(ctx)
	if err != nil {
		return fmt.Errorf("fetch log: %w", err)
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
	sortedTrans, unsortedTrans := sorter.Sort(trans)

	// Send Slack notifications for unsorted transactions
	slack.ReportUnsorted(unsortedTrans)

	// Export transactions
	if err := exporter.Export(sortedTrans, os.Stdout); err != nil {
		return fmt.Errorf("export: %w", err)
	}

	// Update deduplicator state
	if err := dedup.Update(lastTran); err != nil {
		return fmt.Errorf("update dedup: %w", err)
	}
	return nil
}
