// QIF (for Quicken interchange format) Formatter.

package exporter

import (
	"fmt"
	"os"
	"strings"

	"github.com/tuxofil/p24fetch/schema"
)

// QIF templates
const (
	dateLayout       = "2006-01-02"
	qifHeader        = "!Account\nN%s\n^\n"
	qifSimple        = "!Type:Bank\nD%s\nT%.2f\nP%s\nS%s\n$%.2f\n^\n"
	qifWithComission = "!Type:Bank\nD%s\nT%.2f\nP%s\nS%s\n$%.2f\nS%s\n$%.2f\n^\n"
)

// Format transactions to QIF.
func ExportToQIF(
	trans []schema.Transaction,
	srcAccName string,
	comissionsAccName string,
	path string,
) error {
	fd, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	fileInfo, err := fd.Stat()
	if err != nil {
		_ = fd.Close()
		return fmt.Errorf("fstat: %w", err)
	}
	if fileInfo.Size() == 0 {
		if _, err := fd.Write([]byte(fmt.Sprintf(qifHeader, srcAccName))); err != nil {
			_ = fd.Close()
			return fmt.Errorf("write header: %w", err)
		}
	}
	for _, tran := range trans {
		var s string
		if comission := tran.Comission(); comission > 0.01 {
			s = fmt.Sprintf(qifWithComission, tran.Date.Format(dateLayout),
				tran.SrcVal, rmNLs(tran.Note), tran.Dst, tran.DstVal,
				comissionsAccName, comission)
		} else {
			s = fmt.Sprintf(qifSimple, tran.Date.Format(dateLayout),
				tran.SrcVal, rmNLs(tran.Note), tran.Dst, -tran.SrcVal)
		}
		if _, err := fd.Write([]byte(s)); err != nil {
			_ = fd.Close()
			return err
		}
	}
	if err := fd.Close(); err != nil {
		return fmt.Errorf("file close: %w", err)
	}
	return nil
}

// Replace all new line chars
func rmNLs(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\n", " "), "\r", "")
}
