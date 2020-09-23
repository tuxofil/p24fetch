package sorter

import (
	"fmt"

	"github.com/tuxofil/p24fetch/config"
	"github.com/tuxofil/p24fetch/schema"
)

type Sorter struct {
	// Configuration used to create the instance.
	config config.Config
	// Matching rules
	rules *Rules
}

// Create new Sorter instance
func New(cfg *config.Config) (*Sorter, error) {
	s := &Sorter{config: *cfg}
	rules, err := ReadRules(cfg.RulesPath)
	if err != nil {
		return nil, fmt.Errorf("create rules: %w", err)
	}
	s.rules = rules
	return s, nil
}

// Sort transactions according to rules.
func (s *Sorter) Sort(trans []schema.Transaction) (
	ignore []schema.Transaction,
	mapped []schema.Transaction,
	unmapped []schema.Transaction,
) {
	var (
		good []schema.Transaction
		bad  []schema.Transaction
	)
	for _, tran := range trans {
		// Check transaction
		if tran.Raw != nil {
			bad = append(bad, tran)
			continue
		} else if tran.Error != "" {
			bad = append(bad, tran)
			continue
		} else if tran.SrcVal >= 0 {
			tran.Error = "deposits are not implemented"
			bad = append(bad, tran)
			continue
		} else if tran.SrcCur != tran.DstCur {
			tran.Error = "currencies differ"
			bad = append(bad, tran)
			continue
		}

		if s.rules.IsIgnored(tran.Dst) || s.rules.IsIgnored(tran.Note) {
			ignore = append(ignore, tran)
			continue
		}

		// Map transaction
		var dstAcc string
		if n := s.rules.Map(tran.Dst); n != "" {
			dstAcc = n
		} else if n := s.rules.Map(tran.Note); n != "" {
			dstAcc = n
		} else {
			bad = append(bad, tran)
			continue
		}

		// Convert transaction
		tran.Note = fmt.Sprintf("%s: %s", tran.Dst, tran.Note)
		tran.Dst = dstAcc
		good = append(good, tran)
	}
	return ignore, good, bad
}
