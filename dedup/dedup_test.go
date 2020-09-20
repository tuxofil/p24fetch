package dedup

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tuxofil/p24fetch/config"
	"github.com/tuxofil/p24fetch/schema"
)

func TestStateIsZero(t *testing.T) {
	testset := []struct {
		State  state
		Expect bool
	}{
		{state{}, true},
		{state{Date: "2020-09-20"}, true},
		{state{Time: "12:17:00"}, true},
		{state{Date: "2020-09-20", Time: "12:17:00"}, false},
	}
	for n, test := range testset {
		assert.Equal(t, test.Expect, test.State.IsZero(),
			"test case #%d: %+v", n, test)
	}
}

func TestStateFilter(t *testing.T) {
	testset := []struct {
		State  state
		Tran   schema.XMLTransaction
		Expect bool
	}{
		{state{}, schema.XMLTransaction{}, false},
		{state{Date: "2020-09-20", Time: "12:17:00"},
			schema.XMLTransaction{
				TranDate: "",
				TranTime: "",
			}, false},
		{state{Date: "2020-09-20", Time: "12:17:00"},
			schema.XMLTransaction{
				TranDate: "2020-09-20",
				TranTime: "12:16:59",
			}, true},
		{state{Date: "2020-09-20", Time: "12:17:00"},
			schema.XMLTransaction{
				TranDate: "2020-09-20",
				TranTime: "12:17:00",
			}, true},
		{state{Date: "2020-09-20", Time: "12:17:00"},
			schema.XMLTransaction{
				TranDate: "2020-09-20",
				TranTime: "12:17:01",
			}, false},
		{state{Date: "2020-09-20", Time: "12:17:25"},
			schema.XMLTransaction{
				TranDate: "2020-09-20",
				TranTime: "12:17:25",
			}, true},
		{state{Date: "2020-09-20", Time: "12:17:25"},
			schema.XMLTransaction{
				TranDate: "2020-09-20",
				TranTime: "12:17:30",
			}, false},
	}
	for n, test := range testset {
		assert.Equal(t, test.Expect, test.State.Filter(test.Tran),
			"test case #%d: %+v", n, test)
	}
}

func TestDedup(t *testing.T) {
	cfg := &config.Config{
		CardNumber: "abcd",
		DedupDir:   "testdata/run/dedup",
	}
	require.NoError(t, os.RemoveAll(cfg.DedupDir))
	dedup, err := New(cfg)
	require.NoError(t, err)
	require.NotNil(t, dedup)

	trans := []schema.XMLTransaction{
		{TranDate: "2020-09-20", TranTime: "12:17:25"},
		{TranDate: "2020-09-20", TranTime: "12:17:30"},
		{TranDate: "2020-09-20", TranTime: "12:17:35"},
	}
	require.Equal(t, trans, dedup.Filter(trans))

	require.NoError(t, dedup.Update(trans[0]))
	require.Equal(t, trans[1:], dedup.Filter(trans))

	require.NoError(t, dedup.Update(trans[1]))
	require.Equal(t, trans[2:], dedup.Filter(trans))

	require.NoError(t, dedup.Update(trans[2]))
	require.Equal(t, []schema.XMLTransaction(nil), dedup.Filter(trans))
}
