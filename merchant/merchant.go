package merchant

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/tuxofil/p24fetch/config"
	"github.com/tuxofil/p24fetch/schema"
)

type Merchant struct {
	// Configuration used to create the Merchant
	config config.Config
	// HTTP client
	httpClient *http.Client
}

// Module initialization hook.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Create new Merchant instance.
func New(cfg *config.Config) (*Merchant, error) {
	return &Merchant{
		config: *cfg,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				ResponseHeaderTimeout: 5 * time.Second,
			},
		},
	}, nil
}

// Fetch transaction log for the configured account.
func (m *Merchant) FetchLog(ctx context.Context) ([]schema.XMLTransaction, error) {
	var (
		wait      = 10 // in seconds
		test      = 0
		paymentID = fmt.Sprintf("%x", rand.Uint64())
		fromDate  = time.Now().UTC().Add(-time.Hour * 24 *
			time.Duration(m.config.Days))
		toDate = time.Now().UTC()
	)

	// Generate request body
	data := fmt.Sprintf("<oper>cmt</oper>"+
		"<wait>%d</wait>"+
		"<test>%d</test>"+
		`<payment id="%s">`+
		`  <prop name="sd" value="%s" />`+
		`  <prop name="ed" value="%s" />`+
		`  <prop name="card" value="%s" />`+
		`</payment>`,
		wait, test, paymentID,
		fromDate.Format("02.01.2006"),
		toDate.Format("02.01.2006"),
		m.config.CardNumber)
	signature := sha1hex(md5hex(data + m.config.MerchantPassword))
	reqBuf := bytes.NewBufferString(fmt.Sprintf(
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+
			`<request version="1.0">`+
			`  <merchant>`+
			`    <id>%d</id>`+
			`    <signature>%s</signature>`+
			`  </merchant>`+
			`  <data>%s</data>`+
			`</request>`,
		m.config.MerchantID, signature, data))

	// Perform HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://api.privatbank.ua/p24api/rest_fiz", reqBuf)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response status: %s", resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Parse XML from the response
	var parsedXML xmlResponse
	if err := xml.Unmarshal(body, &parsedXML); err != nil {
		return nil, fmt.Errorf("parse xml: %w", err)
	}

	// Revert the list
	var (
		xmlTrans = parsedXML.Data.Info.Statements.Statement
		reverted []schema.XMLTransaction
	)
	for i := len(xmlTrans) - 1; i >= 0; i-- {
		reverted = append(reverted, xmlTrans[i])
	}
	return reverted, nil
}

func md5hex(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func sha1hex(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}
