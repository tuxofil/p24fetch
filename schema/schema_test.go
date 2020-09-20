package schema

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseCurrency(t *testing.T) {
	type Expect struct {
		Value Currency
		Error bool
	}
	testset := []struct {
		Subject string
		Expect  Expect
	}{
		{"uah", Expect{Value: UAH}},
		{"UAH", Expect{Value: UAH}},
		{"UaH", Expect{Value: UAH}},
		{"USD", Expect{Value: USD}},
		{"EUR", Expect{Value: EUR}},
		{"", Expect{Error: true}},
		{"ZZZ", Expect{Error: true}},
	}
	for n, test := range testset {
		c, e := ParseCurrency(test.Subject)
		assert.Equal(t, test.Expect, Expect{c, e != nil},
			"test case #%d: %+v", n, test)
	}
}

func TestParseAmount(t *testing.T) {
	type Expect struct {
		Amount   float32
		Currency Currency
		Error    bool
	}
	testset := []struct {
		Subject string
		Expect  Expect
	}{
		{"0 UAH", Expect{Amount: 0, Currency: UAH}},
		{"1 UAH", Expect{Amount: 1, Currency: UAH}},
		{"3.1 UAH", Expect{Amount: 3.1, Currency: UAH}},
		{"0.11 USD", Expect{Amount: 0.11, Currency: USD}},
		{"-2.11 EUR", Expect{Amount: -2.11, Currency: EUR}},
	}
	for n, test := range testset {
		a, c, e := ParseAmount(test.Subject)
		assert.Equal(t, test.Expect, Expect{a, c, e != nil},
			"test case #%d: %+v", n, test)
	}
}

func TestParseTime(t *testing.T) {
	type Expect struct {
		Value time.Time
		Error bool
	}
	testset := []struct {
		Date, Time string
		Expect     Expect
	}{
		{"", "", Expect{Error: true}},
		{"2020-09-19", "12:13:14",
			Expect{Value: time.Date(2020, 9, 19, 12, 13, 14, 0, time.UTC)}},
	}
	for n, test := range testset {
		d, e := ParseTime(test.Date, test.Time)
		assert.Equal(t, test.Expect, Expect{d, e != nil},
			"test case #%d: %+v", n, test)
	}
}

func TestComission(t *testing.T) {
	testset := []struct {
		Tran   Transaction
		Expect float32
	}{
		{Transaction{}, 0},
		{Transaction{SrcVal: -2.25, SrcCur: UAH, DstVal: 2.0, DstCur: UAH}, 0.25},
		{Transaction{SrcVal: -2.00, SrcCur: UAH, DstVal: 2.0, DstCur: UAH}, 0.0},
		{Transaction{SrcVal: -2.01, SrcCur: UAH, DstVal: 2.0, DstCur: USD}, 0.0},
	}
	for n, test := range testset {
		assert.Equal(t, test.Expect, test.Tran.Comission(),
			"test case #%d: %+v", n, test)
	}
}
