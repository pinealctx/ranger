package jda

import "fmt"

type SymbolCtx struct {
	SymbolId int64
	Name     string
	Bid      float64
	Ask      float64
	Leverage float64

	Quote       SymbolQuote
	RefName     string
	RefSymbolId int64

	RefBid float64 // Bid price in quote currency
	RefAsk float64 // Ask price in quote currency
}

func (x *SymbolCtx) checkPosition(pos *Position) bool {
	isLong := pos.Type == "BUY"
	var pnl float64
	qty := pos.Qty.F64()
	match := true
	if isLong {
		if x.Bid != pos.PriceCurr {
			if pos.Qty != 0 {
				fmt.Printf("Positon currenct price is not matching symbol bid: %f != %f\n", x.Bid, pos.PriceCurr)
				match = false
			}
		}
		pnl = pos.PriceCurr*qty - pos.Price*qty
	} else {
		if x.Ask != pos.PriceCurr {
			if pos.Qty != 0 {
				fmt.Printf("Positon currenct price is not matching symbol ask: %f != %f\n", x.Ask, pos.PriceCurr)
				match = false
			}
		}
		pnl = pos.Price*qty - pos.PriceCurr*qty
	}
	qutoePnl := x.quoteAsUsd(pnl)
	quoteValue := x.quoteAsUsd(pos.PriceCurr * qty)

	if !equalsF64(pnl, pos.Pnl) {
		fmt.Printf("Position PnL mismatch: %f != %f\n", pnl, pos.Pnl)
		match = false
	}
	if !equalsUsd(qutoePnl, pos.QuotePnl) {
		fmt.Printf("Position Quote PnL mismatch: %f != %f\n", qutoePnl, pos.QuotePnl)
		match = false
	}
	if !equalsUsd(quoteValue, pos.QuoteValue) {
		fmt.Printf("Position Quote Value mismatch: %f != %f\n", quoteValue, pos.QuoteValue)
		match = false
	}
	return match
}

func (x *SymbolCtx) checkOpenOrd(symbolId int64, oid string, ord *Order) bool {
	price := ord.Price.F64()
	if price == 0 {
		isLong := ord.Side == "BUY"
		if isLong {
			price = x.Ask
		} else {
			price = x.Bid
		}
	}
	quoteValue := x.quoteAsUsd(price * (ord.Qty.F64() - ord.CumQty))
	if !equalsUsd(quoteValue, ord.QuoteValue) {
		fmt.Printf("Order SymbolId:%+v orderId: %+v Quote Value mismatch: %f != %f\n",
			symbolId, oid, quoteValue, ord.QuoteValue)
		return false
	}
	return true
}

func (x *SymbolCtx) feedRef(qt SymbolQuote, refId int64, refName string, refSnap *SymbolSnap) {
	x.Quote = qt
	x.RefName = refName
	x.RefSymbolId = refId
	x.RefBid = refSnap.Bid
	x.RefAsk = refSnap.Ask
}

func (x *SymbolCtx) quoteAsUsd(v float64) float64 {
	switch x.Quote {
	case SqQuote:
		return v
	case SqBase:
		// if v >= 0, use ask price, otherwise use bid price
		if v >= 0 {
			return v / x.Ask
		} else {
			return v / x.Bid
		}
	case SqRefQuote:
		// if v >= 0, use ref bid price, otherwise use ref ask price
		if v >= 0 {
			return v * x.RefBid
		} else {
			return v * x.RefAsk
		}
	case SqRefBase:
		// if v >= 0, use ref ask price, otherwise use ref bid price
		if v >= 0 {
			return v / x.RefAsk
		} else {
			return v / x.RefBid
		}
	}
	panic(fmt.Sprintf("Unknown quote type: %v", x.Quote))
}
