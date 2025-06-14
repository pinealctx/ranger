package jda

import "fmt"

type MarginAudit struct {
	order  *NewOrder
	allow  *MarginAllow
	accPos *CustomerAccountPositions
	cc     *Runtime
	ctx    *SymbolCtx
	pos    *Position
}

func NewMarginAudit(order *NewOrder, allow *MarginAllow, accPos *CustomerAccountPositions) *MarginAudit {
	return &MarginAudit{
		order:  order,
		allow:  allow,
		accPos: accPos,
	}
}

func (x *MarginAudit) Validate() bool {
	ok := x.validateOrder()
	if !ok {
		fmt.Println("Order validation failed")
		return false
	}
	ok = x.validateAllowBidAsk()
	if !ok {
		fmt.Println("Allow Bid/Ask validation failed")
		return false
	}
	hasPrvPosition := x.hasPrvPosition()
	if hasPrvPosition {
		if x.order.Side == x.pos.Type {
			ok = x.validateSamePosition()
			if !ok {
				fmt.Println("Same position validation failed")
				return false
			}
		} else {
			ok = x.validatePositionClose()
			if !ok {
				fmt.Println("Position close validation failed")
				return false
			}
		}
	} else {
		ok = x.validateNoPrvPosition()
		if !ok {
			fmt.Println("No previous position validation failed")
			return false
		}
	}

	ok = x.validateMarginLevel()
	if !ok {
		fmt.Println("Margin level validation failed")
		return false
	}
	return true
}

func (x *MarginAudit) validateOrder() bool {
	if x.order.OrdType == "MARKET" && x.order.Price <= 0 {
		tdPrice := tradePrice(x.order, x.allow)
		if x.allow.TradePrice != tdPrice {
			fmt.Printf("Order price does not match trade price: %f != %f\n", tdPrice, x.allow.TradePrice)
			return false
		}
	} else {
		orderPrice := x.order.Price.F64()
		if orderPrice != x.allow.TradePrice {
			fmt.Printf("Order price does not match trade price: %f != %f\n", orderPrice, x.allow.TradePrice)
			return false
		}
	}
	ordQty := x.order.OrderQty.F64()
	if ordQty != x.allow.TradeQty {
		fmt.Printf("Order quantity does not match trade quantity: %f != %f\n", ordQty, x.allow.TradeQty)
		return false
	}
	return true
}

func (x *MarginAudit) validateAllowBidAsk() bool {
	cc := NewRuntime(x.accPos)
	x.cc = cc
	ctx := cc.buildSymbolCtxByName(x.order.Symbol)
	x.ctx = ctx
	if x.allow.BidPrice != ctx.Bid {
		fmt.Printf("BidPrice does not match context bid price: %f != %f\n", x.allow.BidPrice, ctx.Bid)
		return false
	}
	if x.allow.AskPrice != ctx.Ask {
		fmt.Printf("AskPrice does not match context ask price: %f != %f\n", x.allow.AskPrice, ctx.Ask)
		return false
	}
	return true
}

func (x *MarginAudit) validateNoPrvPosition() bool {
	if x.allow.BeforeQuotePnl != 0 {
		fmt.Printf("BeforeQuotePnl should be zero, but got: %f\n", x.allow.BeforeQuotePnl)
		return false
	}
	if x.allow.BeforeQuoteNotion != 0 {
		fmt.Printf("BeforeQuoteNotion should be zero, but got: %f\n", x.allow.BeforeQuoteNotion)
		return false
	}

	if x.allow.LeftSide != convertSideType(x.order.Side) {
		fmt.Printf("LeftSide does not match order side: %s != %s\n", x.allow.LeftSide, convertSideType(x.order.Side))
		return false
	}
	if x.allow.LeftQty != x.allow.TradeQty {
		fmt.Printf("LeftQty does not match expected value: %f != %f\n", x.allow.LeftQty, x.allow.TradeQty)
		return false
	}
	leftPrice := reverseTradePrice(x.order, x.allow)
	if x.allow.LeftPrice != leftPrice {
		fmt.Printf("LeftPrice does not match expected value: %f != %f\n", x.allow.LeftPrice, leftPrice)
		return false
	}

	return x._validateMarginAllowAfter()
}

func (x *MarginAudit) validateSamePosition() bool {
	if x.allow.BeforeQuotePnl != x.pos.QuotePnl {
		fmt.Printf("BeforeQuotePnl does not match position QuotePnl: %f != %f\n", x.allow.BeforeQuotePnl, x.pos.QuotePnl)
		return false
	}
	if x.allow.BeforeQuoteNotion != x.pos.QuoteValue {
		fmt.Printf("BeforeQuoteNotion does not match position QuoteValue: %f != %f\n", x.allow.BeforeQuoteNotion, x.pos.QuoteValue)
		return false
	}

	if x.allow.LeftSide != convertSideType(x.order.Side) {
		fmt.Printf("LeftSide does not match order side: %s != %s\n", x.allow.LeftSide, convertSideType(x.order.Side))
		return false
	}
	pQty := x.pos.Qty.F64()
	leftQty := pQty + x.allow.TradeQty
	if x.allow.LeftQty != leftQty {
		fmt.Printf("LeftQty does not match expected value: %f != %f\n", x.allow.LeftQty, leftQty)
		return false
	}
	leftPrice := (x.pos.Price*pQty + x.allow.TradePrice*x.allow.TradeQty) / leftQty
	if x.allow.LeftPrice != leftPrice {
		fmt.Printf("LeftPrice does not match expected value: %f != %f\n", x.allow.LeftPrice, leftPrice)
		return false
	}

	if !x._validateHasPosBegin() {
		return false
	}

	return x._validateMarginAllowAfter()
}

func (x *MarginAudit) validatePositionClose() bool {
	posQty := x.pos.Qty.F64()
	var closeQty float64
	var closeValue float64
	if x.allow.TradeQty > posQty {
		x.allow.LeftSide = convertSideType(x.order.Side)
		x.allow.LeftQty = x.allow.TradeQty - posQty
		x.allow.LeftPrice = x.allow.TradePrice
		closeQty = posQty
	} else {
		x.allow.LeftSide = convertSideType(x.pos.Type)
		x.allow.LeftQty = posQty - x.allow.TradeQty
		x.allow.LeftPrice = x.pos.Price
		closeQty = x.allow.TradeQty
	}

	if x.order.Side == "BUY" {
		closeValue = x.pos.Price*closeQty - x.allow.TradePrice*closeQty
	} else {
		closeValue = x.allow.TradePrice*closeQty - x.pos.Price*closeQty
	}
	tradeQuotePnl, tradePnlRate := x.ctx.quoteAsUsd(closeValue)
	if x.allow.TradeQuotePnl != tradeQuotePnl {
		fmt.Printf("TradeQuotePnl does not match expected value: %f != %f\n", x.allow.TradeQuotePnl, tradeQuotePnl)
		return false
	}
	if x.allow.TradeQuoteRate != tradePnlRate {
		fmt.Printf("TradeQuoteRate does not match expected value: %f != %f\n", x.allow.TradeQuoteRate, tradePnlRate)
		return false
	}
	if !x._validateHasPosBegin() {
		return false
	}
	return x._validateMarginAllowAfter()
}

func (x *MarginAudit) validateMarginLevel() bool {
	// remove processing order in customer account positions
	symbolId, err := B85Conv.Parse(x.order.Symbol)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse symbol ID: %s, error: %v", x.order.Symbol, err))
	}
	symbolMap := x.accPos.ProcessingPositionHashMap[StringInt64(symbolId)]
	if symbolMap == nil {
		fmt.Printf("No processing orders found for symbol: %s\n", x.order.Symbol)
		return false
	}
	_, ok := symbolMap[x.order.ClOrdID]
	if !ok {
		fmt.Printf("Order %s not found in processing orders for symbol %s\n", x.order.ClOrdID, x.order.Symbol)
		return false
	}
	delete(symbolMap, x.order.ClOrdID)

	_, _, totalMarginReq, match := x.cc.calAll()
	if !match {
		fmt.Printf("Margin level validation failed for symbol: %s\n", x.order.Symbol)
		return false
	}
	deltaQuoteNotion, deltaQuotePnl := x.allow.Delta()
	equity := x.accPos.Balance + x.accPos.CreditLimit + x.accPos.OpenPNL + deltaQuotePnl
	marginReq := totalMarginReq + deltaQuoteNotion/x.allow.Leverage
	marginLevel := calMarginLevel(equity, marginReq)
	if !equalsF64(marginLevel, x.allow.MarginLevel) {
		fmt.Printf("Margin Level mismatch: %f != %f\n", marginLevel, x.allow.MarginLevel)
		return false
	}
	return true
}

func (x *MarginAudit) _validateHasPosBegin() bool {
	var notionValue, pnl float64
	posQty := x.pos.Qty.F64()

	if x.pos.Type == "BUY" {
		if x.pos.PriceCurr != x.allow.BidPrice {
			fmt.Printf("Position current price does not match bid price: %f != %f\n", x.pos.PriceCurr, x.allow.BidPrice)
			return false
		}
		if x.pos.PriceCurr != x.ctx.Bid {
			fmt.Printf("Position current price does not match context bid price: %f != %f\n", x.pos.PriceCurr, x.ctx.Bid)
			return false
		}
		pnl = x.pos.PriceCurr*posQty - x.pos.Price*posQty
	} else {
		if x.pos.PriceCurr != x.allow.AskPrice {
			fmt.Printf("Position current price does not match ask price: %f != %f\n", x.pos.PriceCurr, x.allow.AskPrice)
			return false
		}
		if x.pos.PriceCurr != x.ctx.Ask {
			fmt.Printf("Position current price does not match context ask price: %f != %f\n", x.pos.PriceCurr, x.ctx.Ask)
			return false
		}
		pnl = x.pos.Price*posQty - x.pos.PriceCurr*posQty
	}

	notionValue = x.pos.PriceCurr * posQty
	beforeQuoteNotion, beforeNotionRate := x.ctx.quoteAsUsd(notionValue)
	if x.allow.BeforeQuoteNotion != beforeQuoteNotion {
		fmt.Printf("BeforeQuoteNotion does not match expected value: %f != %f\n", x.allow.BeforeQuoteNotion, beforeQuoteNotion)
		return false
	}
	if x.allow.QuoteBeforeNotionRate != beforeNotionRate {
		fmt.Printf("QuoteBeforeNotionRate does not match expected value: %f != %f\n", x.allow.QuoteBeforeNotionRate, beforeNotionRate)
		return false
	}

	beforePnlQuote, beforePnlRate := x.ctx.quoteAsUsd(pnl)
	if x.allow.BeforeQuotePnl != beforePnlQuote {
		fmt.Printf("BeforeQuotePnl does not match expected value: %f != %f\n", x.allow.BeforeQuotePnl, beforePnlQuote)
		return false
	}

	_beforePnlRate := x.pos.QuotePnl / x.pos.Pnl
	if _beforePnlRate != beforePnlRate {
		fmt.Printf("QuotePnlRate does not match expected value: %f != %f\n", _beforePnlRate, beforePnlRate)
		return false
	}

	return true
}

func (x *MarginAudit) _validateMarginAllowAfter() bool {
	notionValue := x.allow.LeftQty * x.allow.LeftPrice
	afterQuoteNotion, afterNotionRate := x.ctx.quoteAsUsd(notionValue)
	if x.allow.AfterQuoteNotion != afterQuoteNotion {
		fmt.Printf("AfterQuoteNotion does not match expected value: %f != %f\n", x.allow.AfterQuoteNotion, afterQuoteNotion)
		return false
	}
	if x.allow.QuoteAfterNotionRate != afterNotionRate {
		fmt.Printf("QuoteAfterNotionRate does not match expected value: %f != %f\n", x.allow.QuoteAfterNotionRate, afterNotionRate)
		return false
	}

	pnl := allowPnl(x.order, x.allow)
	afterQuotePnl, afterPnlRate := x.ctx.quoteAsUsd(pnl)
	if x.allow.AfterQuotePnl != afterQuotePnl {
		fmt.Printf("AfterQuotePnl does not match expected value: %f != %f\n", x.allow.AfterQuotePnl, afterQuotePnl)
		return false
	}
	if x.allow.QuotePnlRate != afterPnlRate {
		fmt.Printf("QuoteAfterPnlRate does not match expected value: %f != %f\n", x.allow.QuotePnlRate, afterPnlRate)
		return false
	}
	if x.allow.TradeQuotePnl != 0 {
		fmt.Printf("TradeQuotePnl should be zero, but got: %f\n", x.allow.TradeQuotePnl)
		return false
	}
	return true
}

func (x *MarginAudit) hasPrvPosition() bool {
	pos, ok := x.accPos.PositionMap[x.order.Symbol]
	if !ok {
		return false
	}
	if pos.Qty == 0 {
		return false
	}
	if pos.Qty < 0 {
		panic(fmt.Sprintf("Negative position Qty cannot be negative: %f, account:%s symbol:%s",
			pos.Qty, x.accPos.AccountId, x.order.Symbol))
	}
	x.pos = pos
	return true
}

func allowPnl(order *NewOrder, allow *MarginAllow) float64 {
	var pnl float64
	if order.Side == "BUY" {
		pnl = allow.LeftPrice*allow.LeftQty - allow.TradePrice*allow.LeftQty
	} else {
		pnl = allow.TradePrice*allow.LeftQty - allow.LeftPrice*allow.LeftQty
	}
	return pnl
}

func tradePrice(order *NewOrder, allow *MarginAllow) float64 {
	if order.Side == "BUY" {
		return allow.AskPrice
	}
	return allow.BidPrice
}

func reverseTradePrice(order *NewOrder, allow *MarginAllow) float64 {
	if order.Side == "BUY" {
		return allow.BidPrice
	}
	return allow.AskPrice
}

func convertSideType(side string) string {
	if side == "BUY" {
		return "1"
	} else if side == "SELL" {
		return "2"
	}
	panic(fmt.Sprintf("Unknown side type: %s", side))
}
