package sorter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
)

type Rules struct {
	// Mapping: ShortID -> GnuCash Account ID
	Accounts map[string]string `json:"accounts"`
	// Matcher rules. Transactions matching one of these
	// patterns will be silently ignored.
	Ignore []string `json:"ignore"`
	// Matcher rules. Every element is a mapping:
	//  ShortID -> list of patterns
	Rules []map[string][]string `json:"rules"`
	// Compiled regexps cache
	regexps map[string]*regexp.Regexp
}

// Read rules from file and validate it.
func ReadRules(path string) (*Rules, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	rules := &Rules{}
	if err := json.Unmarshal(data, rules); err != nil {
		return nil, fmt.Errorf("parse JSON: %w", err)
	}
	if err := rules.Validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}
	if err := rules.CompilePatterns(); err != nil {
		return nil, fmt.Errorf("compile regexps: %w", err)
	}
	return rules, nil
}

// Compile all regexp patterns.
func (r *Rules) CompilePatterns() error {
	r.regexps = make(map[string]*regexp.Regexp)
	for _, pattern := range r.Ignore {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("ignore: invalid pattern %#v: %w",
				pattern, err)
		}
		r.regexps[pattern] = re
	}
	for _, rule := range r.Rules {
		for name, patterns := range rule {
			for _, pattern := range patterns {
				re, err := regexp.Compile(pattern)
				if err != nil {
					return fmt.Errorf(
						"rule %#v: invalid pattern %#v: %w",
						name, pattern, err)
				}
				r.regexps[pattern] = re
			}
		}
	}
	return nil
}

// Validate rules.
func (r *Rules) Validate() error {
	if len(r.Accounts) == 0 {
		return errors.New("no accounts defined")
	}
	if len(r.Rules) == 0 {
		return errors.New("no rules defined")
	}
	for _, rule := range r.Rules {
		for name := range rule {
			if _, ok := r.Accounts[name]; !ok {
				return fmt.Errorf("rules: undefined account: %#v", name)
			}
		}
	}
	return nil
}

// Ignore returns true when given string were matched to
// one of 'ignore' rules.
func (r *Rules) IsIgnored(s string) bool {
	for _, pattern := range r.Ignore {
		if r.regexps[pattern].MatchString(s) {
			return true
		}
	}
	return false
}

// Traverse matching rules for GnuCash Account ID.
func (r *Rules) Map(s string) string {
	for _, rule := range r.Rules {
		for shortID, patterns := range rule {
			for _, pattern := range patterns {
				if r.regexps[pattern].MatchString(s) {
					return r.Accounts[shortID]
				}
			}
		}
	}
	return ""
}
