package dedup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/tuxofil/p24fetch/config"
	"github.com/tuxofil/p24fetch/schema"
)

type Deduplicator struct {
	// Configuration used to create the instance
	config config.Config
	// Last processed entry
	state state
}

type state struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

// Create new deduplicator instance.
func New(cfg *config.Config) (*Deduplicator, error) {
	if err := os.MkdirAll(cfg.DedupDir, 0700); err != nil {
		return nil, fmt.Errorf("create state dir: %w", err)
	}
	dedup := &Deduplicator{config: *cfg}
	data, err := ioutil.ReadFile(dedup.stateFileName())
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("read state file: %w", err)
		}
	} else if err := json.Unmarshal(data, &dedup.state); err != nil {
		return nil, fmt.Errorf("parse state: %w", err)
	}
	return dedup, nil
}

// Filter filters transaction log according to deduplicator saved state.
func (d *Deduplicator) Filter(trans []schema.XMLTransaction) []schema.XMLTransaction {
	if d.state.IsZero() {
		return trans
	}
	var res []schema.XMLTransaction
	for _, tran := range trans {
		if !d.state.Filter(tran) {
			res = append(res, tran)
		}
	}
	return res
}

// Update saves transaction to deduplicator's internal state as the
// last processed transaction.
func (d *Deduplicator) Update(tran schema.XMLTransaction) error {
	newState := state{Date: tran.TranDate, Time: tran.TranTime}
	if newState.IsZero() {
		return fmt.Errorf("invalid state: %+v", newState)
	}
	data, err := json.Marshal(newState)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err := ioutil.WriteFile(d.stateFileName(), data, 0600); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	d.state = newState
	return nil
}

func (d *Deduplicator) stateFileName() string {
	return path.Join(d.config.DedupDir, d.config.CardNumber+".json")
}

func (s *state) IsZero() bool {
	return s.Date == "" || s.Time == ""
}

func (s *state) Filter(tran schema.XMLTransaction) bool {
	s2 := state{Date: tran.TranDate, Time: tran.TranTime}
	return !s2.IsZero() && s.String() >= s2.String()
}

func (s *state) String() string {
	return fmt.Sprintf("%sT%s", s.Date, s.Time)
}
