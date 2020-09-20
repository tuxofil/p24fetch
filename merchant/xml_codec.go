package merchant

import "github.com/tuxofil/p24fetch/schema"

type xmlMerchant struct {
	ID        string `xml:"id"`
	Signature string `xml:"signature"`
}

type xmlResponse struct {
	Merchant xmlMerchant     `xml:"merchant"`
	Data     xmlResponseData `xml:"data"`
}

type xmlResponseData struct {
	Oper string `xml:"oper"`
	Info struct {
		Statements struct {
			Status    string                  `xml:"status,attr"`
			Credit    float32                 `xml:"credit,attr"`
			Debet     float32                 `xml:"debet,attr"`
			Statement []schema.XMLTransaction `xml:"statement"`
		} `xml:"statements"`
	} `xml:"info"`
	Error struct {
		Message string `xml:"message,attr"`
	} `xml:"error"`
}
