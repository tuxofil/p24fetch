package sorter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRulesValidate(t *testing.T) {
	testset := []struct {
		Subject       Rules
		ExpectSuccess bool
	}{
		{ExpectSuccess: false},
		{Subject: Rules{
			Accounts: map[string]string{},
			Rules:    []map[string][]string{}},
			ExpectSuccess: false},
		{Subject: Rules{
			Accounts: map[string]string{
				"acc1": "name1",
			},
			Rules: []map[string][]string{}},
			ExpectSuccess: false},
		{Subject: Rules{
			Accounts: map[string]string{
				"acc1": "name1",
			},
			Rules: []map[string][]string{
				map[string][]string{
					"acc2": []string{"pat1", "pat2"},
				},
			}},
			ExpectSuccess: false},
		{Subject: Rules{
			Accounts: map[string]string{
				"acc1": "name1",
			},
			Rules: []map[string][]string{
				map[string][]string{
					"acc1": []string{"pat1", "pat2"},
				},
			}},
			ExpectSuccess: true},
	}
	for n, test := range testset {
		err := test.Subject.Validate()
		if test.ExpectSuccess {
			assert.NoError(t, err, "test #%d: %+v", n, test)
		} else {
			assert.Error(t, err, "test #%d: %+v", n, test)
		}
	}
}

func TestRulesMap(t *testing.T) {
	testset := []struct {
		Rules   Rules
		Subject string
		Expect  string
	}{
		{Rules: Rules{
			Accounts: map[string]string{
				"acc1": "name1",
			},
			Rules: []map[string][]string{
				map[string][]string{
					"acc1": []string{"pat1", "pat2"},
				},
			}},
			Subject: "pat1",
			Expect:  "name1"},

		{Rules: Rules{
			Accounts: map[string]string{
				"acc1": "name1",
			},
			Rules: []map[string][]string{
				map[string][]string{
					"acc1": []string{"pat1", "pat2"},
				},
			}},
			Subject: "pat2",
			Expect:  "name1"},

		{Rules: Rules{
			Accounts: map[string]string{
				"acc1": "name1",
				"acc2": "name2",
			},
			Rules: []map[string][]string{
				map[string][]string{
					"acc1": []string{"pat1", "pat2"},
				},
				map[string][]string{
					"acc2": []string{"pat1", "pat3"},
				},
			}},
			Subject: "pat1",
			Expect:  "name1"},

		{Rules: Rules{
			Accounts: map[string]string{
				"acc1": "name1",
				"acc2": "name2",
			},
			Rules: []map[string][]string{
				map[string][]string{
					"acc1": []string{"pat1", "pat2"},
				},
				map[string][]string{
					"acc2": []string{"pat1", "pat3"},
				},
			}},
			Subject: "pat3",
			Expect:  "name2"},
	}
	for n, test := range testset {
		require.NoError(t, test.Rules.Validate(),
			"test #%d: %+v", n, test)
		require.NoError(t, test.Rules.CompilePatterns(),
			"test #%d: %+v", n, test)
		assert.Equal(t, test.Expect, test.Rules.Map(test.Subject),
			"test #%d: %+v", n, test)
	}
}
