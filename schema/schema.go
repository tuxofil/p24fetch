package schema

import (
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"
)

type Currency string

// Known currencies
const (
	UAH Currency = "UAH"
	USD Currency = "USD"
	EUR Currency = "EUR"
)

type Transaction struct {
	// Transaction date and time
	Date time.Time
	// From account name
	Src string
	// From amount
	SrcVal float32
	// From currency
	SrcCur Currency
	// To account name
	Dst string
	// To amount
	DstVal float32
	// To currency
	DstCur Currency
	// Transaction note
	Note string
	// Not nil on XML parse error
	Error error `json:"Error,omitempty"`
	// XML Transaction. Set only when Error is not nil.
	Raw *XMLTransaction `json:"Raw,omitempty"`
}

// XMLTransaction.
// This is a part of Privat24 API spec.
type XMLTransaction struct {
	Card        string `xml:"card,attr"`
	AppCode     string `xml:"appcode,attr"`
	TranDate    string `xml:"trandate,attr"`
	TranTime    string `xml:"trantime,attr"`
	Amount      string `xml:"amount,attr"`
	CardAmount  string `xml:"cardamount,attr"`
	Rest        string `xml:"rest,attr"`
	Terminal    string `xml:"terminal,attr"`
	Description string `xml:"description,attr"`
}

// Parse currency
func ParseCurrency(s string) (Currency, error) {
	c := Currency(strings.ToUpper(s))
	switch c {
	case UAH, USD, EUR:
		return c, nil
	}
	return Currency(""), fmt.Errorf("invalid currency: %#v", s)
}

// ParseAmount parses amount and currency.
func ParseAmount(s string) (float32, Currency, error) {
	tokens := strings.SplitN(strings.Trim(s, " \t\n\r"), " ", 2)
	if n := len(tokens); n != 2 {
		return 0, Currency(""),
			fmt.Errorf("invalid tokens count: %d", n)
	}
	amount, err := strconv.ParseFloat(tokens[0], 32)
	if err != nil {
		return 0, Currency(""),
			fmt.Errorf("invalid number: %w", err)
	}
	currency, err := ParseCurrency(tokens[1])
	if err != nil {
		return 0, Currency(""),
			fmt.Errorf("invalid currency: %w", err)
	}
	return float32(amount), currency, nil
}

// ParseTime parses date and time from two strings.
// Note: Privat24 API returns date&time in Kyiv time zone.
func ParseTime(d, t string) (time.Time, error) {
	res, err := time.Parse("2006-01-02T15:04:05", d+"T"+t)
	return res.UTC(), err
}

// Convert transaction object received from the Privat24 API to
// more convenient, Golang native format.
func ParseTransaction(xmlTran XMLTransaction) Transaction {
	date, err := ParseTime(xmlTran.TranDate, xmlTran.TranTime)
	if err != nil {
		return Transaction{
			Error: fmt.Errorf("parse time: %w", err),
			Raw:   &xmlTran,
		}
	}
	fromAmount, fromCurrency, err := ParseAmount(xmlTran.CardAmount)
	if err != nil {
		return Transaction{
			Error: fmt.Errorf("parse src amount: %w", err),
			Raw:   &xmlTran,
		}
	}
	toAmount, toCurrency, err := ParseAmount(xmlTran.Amount)
	if err != nil {
		return Transaction{
			Error: fmt.Errorf("parse dst amount: %w", err),
			Raw:   &xmlTran,
		}
	}
	return Transaction{
		Date:   date,
		Src:    xmlTran.Card,
		SrcVal: fromAmount,
		SrcCur: fromCurrency,
		Dst:    html.UnescapeString(xmlTran.Terminal),
		DstVal: toAmount,
		DstCur: toCurrency,
		Note:   html.UnescapeString(xmlTran.Description),
	}
}

// Comission returns the value of comission charged.
func (t *Transaction) Comission() float32 {
	if t.SrcVal >= 0 || t.SrcCur != t.DstCur {
		return 0
	}
	return -(t.SrcVal + t.DstVal)
}
