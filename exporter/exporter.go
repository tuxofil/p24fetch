package exporter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/tuxofil/p24fetch/config"
	"github.com/tuxofil/p24fetch/schema"
)

type Exporter struct {
	// Configuration used to create the instance.
	config config.Config
}

// Create new exporter instance.
func New(cfg *config.Config) (*Exporter, error) {
	return &Exporter{config: *cfg}, nil
}

// Export transaction log to external storage.
func (e *Exporter) Export(
	trans []schema.Transaction,
	writer io.Writer,
) error {
	switch e.config.ExportFormat {
	case schema.JSON:
		return json.NewEncoder(writer).Encode(trans)
	case schema.QIF:
		return ExportToQIF(trans, e.config.SrcAccountName,
			e.config.ComissionAccountName, writer)
	}
	return fmt.Errorf("not implemented: %s", e.config.ExportFormat)
}
