package jda

import (
	"encoding/json"
	"regexp"
	"strings"
)

var (
	marginAllowReg  = regexp.MustCompile(`allow:(\{.*?\}),`)
	marginAccPosReg = regexp.MustCompile(`account position:(\{.*\})`)
)

type MarginAllow struct {
	Allow                 bool    `json:"allow"`
	ErrorCode             int     `json:"errorCode"`
	BidPrice              float64 `json:"bidPrice"`
	AskPrice              float64 `json:"askPrice"`
	Leverage              float64 `json:"leverage"`
	TradePrice            float64 `json:"tradePrice"`
	TradeQty              float64 `json:"tradeQty"`
	BeforeQuotePnl        float64 `json:"beforeQuotePnl"`
	BeforeQuoteNotion     float64 `json:"beforeQuoteNotion"`
	AfterQuotePnl         float64 `json:"afterQuotePnl"`
	AfterQuoteNotion      float64 `json:"afterQuoteNotion"`
	QuoteBeforeNotionRate float64 `json:"quoteBeforeNotionRate"`
	QuoteAfterNotionRate  float64 `json:"quoteAfterNotionRate"`
	QuotePnlRate          float64 `json:"quotePnlRate"`
	MarginLevel           float64 `json:"marginLevel"`
	LeftQty               float64 `json:"leftQty"`
	LeftPrice             float64 `json:"leftPrice"`
	TradeQuotePnl         float64 `json:"tradeQuotePnl"`
	TradeQuoteRate        float64 `json:"tradeQuoteRate"`
	TradeSide             string  `json:"tradeSide"`
	LeftSide              string  `json:"leftSide"`
}

func FigureMarginAllowAndAccPos(line string) (*MarginAllow, *CustomerAccountPositions, error) {
	if strings.Contains(line, `allow:{"allow":`) {
		allowMatches := marginAllowReg.FindStringSubmatch(line)
		var allowJson, accountJson string
		if len(allowMatches) > 1 {
			allowJson = allowMatches[1]
		}
		accountMatches := marginAccPosReg.FindStringSubmatch(line)
		if len(accountMatches) > 1 {
			accountJson = accountMatches[1]
		}
		if allowJson == "" || accountJson == "" {
			panic("Margin allow or account position JSON not found in line: " + line)
		}
		allow := &MarginAllow{}
		err := json.Unmarshal([]byte(allowJson), allow)
		if err != nil {
			return nil, nil, err
		}
		account := &CustomerAccountPositions{}
		err = json.Unmarshal([]byte(accountJson), account)
		if err != nil {
			return nil, nil, err
		}
		return allow, account, nil
	}
	return nil, nil, nil
}

func (x *MarginAllow) Delta() (quoteNotion, quotePnl float64) {
	// Calculate the delta in quote notion and PnL
	quoteNotion = x.AfterQuoteNotion - x.BeforeQuoteNotion
	quotePnl = x.AfterQuotePnl + x.TradeQuotePnl - x.BeforeQuotePnl
	return quoteNotion, quotePnl
}
