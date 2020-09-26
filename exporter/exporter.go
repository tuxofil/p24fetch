package exporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/tuxofil/p24fetch/config"
	"github.com/tuxofil/p24fetch/schema"
)

type Exporter struct {
	// Configuration used to create the instance.
	config config.Config
}

// Create new exporter instance.
func New(cfg *config.Config) *Exporter {
	return &Exporter{config: *cfg}
}

// Export transaction log to external storage.
func (e *Exporter) Export(trans []schema.Transaction) error {
	if len(trans) == 0 {
		return nil
	}
	if err := os.MkdirAll(e.config.ResultsDir, 0700); err != nil {
		return fmt.Errorf("create results dir: %w", err)
	}
	filePath := path.Join(e.config.ResultsDir, e.config.CardNumber)

	switch e.config.ExportFormat {
	case schema.JSON:
		var buf bytes.Buffer
		encoder := json.NewEncoder(&buf)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(trans); err != nil {
			return fmt.Errorf("encode tran: %w", err)
		}
		filePath += "-" + time.Now().Format("2006-01-02T15-04-05") + ".json"
		if err := ioutil.WriteFile(filePath, buf.Bytes(), 0600); err != nil {
			return fmt.Errorf("write file: %w", err)
		}
		return nil
	case schema.QIF:
		filePath += ".qif"
		return ExportToQIF(trans, e.config.SrcAccountName,
			e.config.ComissionAccountName, filePath)
	}
	return fmt.Errorf("not implemented: %s", e.config.ExportFormat)
}
