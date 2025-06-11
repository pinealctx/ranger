package jda

import (
	"bytes"
	"fmt"
)

type TierPrice struct {
	TierId      string   `json:"tierId"`
	SymbolId    string   `json:"symbolId"`
	SendingTime TimeNano `json:"sendingTime"`
	Bid         float64  `json:"bid"`
	Ask         float64  `json:"ask"`
}

func (x *TierPrice) String() string {
	buf := make([]byte, 0, 100)
	sb := bytes.NewBuffer(buf)
	sb.Reset()
	sb.WriteString("TierPrice{")
	sb.WriteString("TierId: ")
	sb.WriteString(x.TierId)
	sb.WriteString(", SymbolId: ")
	sb.WriteString(x.SymbolId)
	sb.WriteString(", SendingTime: ")
	sb.WriteString(x.SendingTime.String())
	sb.WriteString(", Bid: ")
	sb.WriteString(fmt.Sprintf("%f", x.Bid))
	sb.WriteString(", Ask: ")
	sb.WriteString(fmt.Sprintf("%f", x.Ask))
	sb.WriteString("}")
	return sb.String()
}

type TierPriceEod struct {
	Date       string       `json:"date"`
	TierPrices []*TierPrice `json:"tierPrices"`
}
