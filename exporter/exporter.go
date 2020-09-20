package exporter

import (
	"encoding/json"
	"os"

	"github.com/tuxofil/p24fetch/config"
	"github.com/tuxofil/p24fetch/schema"
)

type Exporter struct{}

// Create new exporter instance.
func New(*config.Config) (*Exporter, error) {
	// TODO:
	return &Exporter{}, nil
}

// Export transaction log to external storage.
func (*Exporter) Export(log []schema.Transaction) error {
	return json.NewEncoder(os.Stdout).Encode(log)
}
