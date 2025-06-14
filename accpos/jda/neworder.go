package jda

import (
	"encoding/json"
	"regexp"
	"strings"
)

var (
	newOrderReg = regexp.MustCompile(`event:(\{.*\})`)
)

type NewOrder struct {
	ClOrdID  string        `json:"clOrdID"`
	Symbol   string        `json:"symbol"`
	OrdType  string        `json:"ordType"` // MARKET or others
	Side     string        `json:"side"`
	Price    NumericString `json:"price"`    // 0.00 for market orders
	OrderQty NumericString `json:"orderQty"` // e.g. "1K" for 1000 units
}

func FigureNewOrder(line string) (*NewOrder, error) {
	if strings.Contains(line, "NewOrderSingle Pre check:") {
		strData := extractNewOrderJson(line)
		order := &NewOrder{}
		err := json.Unmarshal([]byte(strData), order)
		if err != nil {
			return nil, err
		}
		return order, nil
	}
	return nil, nil
}

func extractNewOrderJson(line string) string {
	matches := newOrderReg.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	panic("NewOrder JSON not found in line: " + line)
}
