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
		fromDate = time.Now().UTC().Add(-time.Hour * 24 *
			time.Duration(m.config.Days))
		toDate = time.Now().UTC()
	)

	// Generate request body
	xmlData := xmlRequestData{
		Oper: "cmt",
		Wait: 10, // in seconds
		Payment: xmlPayment{
			ID: fmt.Sprintf("%x", rand.Uint64()),
			Prop: []xmlProp{
				{Name: "sd", Value: fromDate.Format("02.01.2006")},
				{Name: "ed", Value: toDate.Format("02.01.2006")},
				{Name: "card", Value: m.config.CardNumber},
			},
		},
	}
	encodedXMLData, err := xml.Marshal(xmlData)
	if err != nil {
		return nil, fmt.Errorf("marshal data: %w", err)
	}
	xmlReq := xmlRequest{
		Version: "1.0",
		Merchant: xmlMerchant{
			ID: m.config.MerchantID,
			Signature: sha1hex(md5hex(string(encodedXMLData) +
				m.config.MerchantPassword)),
		},
		Data: xmlData,
	}
	encodedXMLReq, err := xml.Marshal(xmlReq)
	if err != nil {
		return nil, fmt.Errorf("marshal req: %w", err)
	}
	encodedXMLReq = append(
		[]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"),
		encodedXMLReq...)
	reqBuf := bytes.NewBuffer(encodedXMLReq)

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
	if reason := parsedXML.Data.Error.Message; reason != "" {
		return nil, fmt.Errorf("API error: %s", reason)
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
