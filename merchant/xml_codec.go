package merchant

import "github.com/tuxofil/p24fetch/schema"

type xmlRequest struct {
	Version  string         `xml:"version,attr"`
	Merchant xmlMerchant    `xml:"merchant"`
	Data     xmlRequestData `xml:"data"`
}

type xmlMerchant struct {
	ID        string `xml:"id"`
	Signature string `xml:"signature"`
}

type xmlRequestData struct {
	Oper    string `xml:"oper"`
	Wait    int    `xml:"wait"`
	Test    int    `xml:"test"`
	Payment struct {
		ID   string    `xml:"id,attr"`
		Prop []xmlProp `xml:"prop"`
	} `xml:"payment"`
}

type xmlProp struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
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
