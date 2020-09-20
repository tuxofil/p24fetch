package sorter

import (
	"github.com/tuxofil/p24fetch/config"
	"github.com/tuxofil/p24fetch/schema"
)

type Sorter struct{}

func New(*config.Config) (*Sorter, error) {
	// TODO:
	return &Sorter{}, nil
}

func (*Sorter) Sort(log []schema.Transaction) ([]schema.Transaction, error) {
	// TODO:
	return log, nil
}
