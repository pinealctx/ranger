package jda

import (
	"bytes"
	"fmt"
	"github.com/pinealctx/ranger/fxb/b85"
)

const (
	MaxMarginLevel = 999999.99

	SqUnknown  SymbolQuote = iota
	SqQuote                // use "USD" as quote currency
	SqBase                 // use "USD" as base currency
	SqRefQuote             // "XAUEUR" -> "EUR" as quote currency, "EURUSD", USD is the quote currency of "EURUSD"
	SqRefBase              // "AUDJPY" -> "JPY" as quote currency, "USDJPY", USD is the base currency of "USDJPY"
)

var (
	b85Conv = b85.NewBase85LongConverterS()
)

// SymbolQuote represents the type of symbol quote in the JDA system.
type SymbolQuote int

type Runtime struct {
	CtxMap        map[int64]*SymbolCtx
	SymbolMap     map[int64]string
	SymbolNameMap map[string]int64
	accPos        *CustomerAccountPositions
}

func NewRuntime(accPos *CustomerAccountPositions) *Runtime {
	x := &Runtime{
		CtxMap:        make(map[int64]*SymbolCtx),
		SymbolMap:     make(map[int64]string),
		SymbolNameMap: make(map[string]int64),
		accPos:        accPos,
	}
	x.buildSymbolMap()
	return x
}

func (x *Runtime) CheckWhole() bool {
	// check all positions
	var totalNotion, totalPnl, totalMarginReq float64
	match := true
	itemMatch := true

	for k, pos := range x.accPos.PositionMap {
		ctx := x.buildSymbolCtxByName(k)
		itemMatch = ctx.checkPosition(pos)
		if !itemMatch {
			match = false
		}
		totalNotion += pos.QuoteValue
		totalPnl += pos.QuotePnl
		totalMarginReq += pos.QuoteValue / ctx.Leverage
	}

	for k, ordMap := range x.accPos.ProcessingPositionHashMap {
		symbolId := int64(k)
		ctx := x.buildSymbolCtx(symbolId)
		for oid, ord := range ordMap {
			itemMatch = ctx.checkOpenOrd(symbolId, oid, ord)
			totalMarginReq += ord.QuoteValue / ctx.Leverage
			if !itemMatch {
				match = false
			}
		}
	}
	if !equalsUsd(totalNotion, x.accPos.NotionalValue) {
		fmt.Printf("Total Notional Value mismatch: %f != %f\n", totalNotion, x.accPos.NotionalValue)
		match = false
	}
	if !equalsUsd(totalPnl, x.accPos.OpenPNL) {
		fmt.Printf("Total PnL mismatch: %f != %f\n", totalPnl, x.accPos.OpenPNL)
		match = false
	}
	if !equalsUsd(totalMarginReq, x.accPos.MarginRequirement) {
		fmt.Printf("Total Margin Requirement mismatch: %f != %f\n", totalMarginReq, x.accPos.MarginRequirement)
		match = false
	}
	equity := x.accPos.Balance + x.accPos.CreditLimit + totalPnl
	if !equalsUsd(equity, x.accPos.Equity) {
		fmt.Printf("Equity mismatch: %f != %f\n", equity, x.accPos.Equity)
		match = false
	}
	availableForTrading := equity - totalMarginReq
	if !equalsUsd(availableForTrading, x.accPos.AvailableForMarginTrading) {
		fmt.Printf("Available for Margin Trading mismatch: %f != %f\n", availableForTrading, x.accPos.AvailableForMarginTrading)
		match = false
	}
	marginLevel := calMarginLevel(equity, totalMarginReq)
	if !ratioEqual(marginLevel, x.accPos.MarginRatio) {
		fmt.Printf("Margin Level mismatch: %f != %f\n", marginLevel, x.accPos.MarginRatio)
		match = false
	}
	return match
}

func (x *Runtime) buildSymbolCtxByName(symbolName string) *SymbolCtx {
	symbolId := x.symbolId(symbolName)
	return x.buildSymbolCtx(symbolId)
}

func (x *Runtime) buildSymbolCtx(symbolId int64) *SymbolCtx {
	if ctx, ok := x.CtxMap[symbolId]; ok {
		return ctx
	}
	symbolName := x.symbolName(symbolId)
	snap, ok := x.accPos.SymbolSnapMap[StringInt64(symbolId)]
	if !ok {
		panic("symbol snap ctx not exist: " + symbolName + ", id:" + fmt.Sprintf("%d", symbolId))
	}
	ctx := &SymbolCtx{
		SymbolId: symbolId,
		Name:     symbolName,
		Bid:      snap.Bid,
		Ask:      snap.Ask,
		Leverage: snap.Leverage,
	}
	base, quote := figureCurrency(symbolName)
	if quote == "USD" {
		ctx.Quote = SqQuote
		return ctx
	}

	if base == "USD" {
		ctx.Quote = SqBase
		return ctx
	}

	usdAsQ := usdAsQuote(quote)
	refId, refSnap, ok := x.figureRef(usdAsQ)
	if ok {
		ctx.feedRef(SqRefQuote, refId, usdAsQ, refSnap)
		return ctx
	}
	usdAsB := usdAsBase(quote)
	refId, refSnap, ok = x.figureRef(usdAsB)
	if ok {
		ctx.feedRef(SqRefBase, refId, usdAsB, refSnap)
		return ctx
	}
	panic("symbol snap not exist: " + symbolName + ", base: " + base + ", quote: " + quote)
}

func (x *Runtime) buildSymbolMap() {
	for k := range x.accPos.SymbolSnapMap {
		symbolId := int64(k)
		symbolName := b85Conv.AsString(symbolId)
		x.SymbolMap[symbolId] = symbolName
		x.SymbolNameMap[symbolName] = symbolId
	}
}

func (x *Runtime) symbolName(symbolId int64) string {
	if name, ok := x.SymbolMap[symbolId]; ok {
		return name
	}
	panic("symbol not exist")
}

func (x *Runtime) symbolId(symbolName string) int64 {
	if id, ok := x.SymbolNameMap[symbolName]; ok {
		return id
	}
	panic("symbol not exist")
}

func (x *Runtime) figureRef(refSymbolName string) (int64, *SymbolSnap, bool) {
	refId, ok := x.SymbolNameMap[refSymbolName]
	if !ok {
		return 0, nil, false
	}
	refSnap, ok := x.accPos.SymbolSnapMap[StringInt64(refId)]
	if !ok {
		panic("ref symbol snap not exist: " + refSymbolName + ", id: " + fmt.Sprintf("%d", refId))
	}
	return refId, refSnap, true
}

func usdAsQuote(currency string) string {
	sb := bytes.NewBuffer(make([]byte, 6))
	sb.Reset()
	sb.WriteString(currency)
	sb.WriteString("USD")
	return sb.String()
}

func usdAsBase(currency string) string {
	sb := bytes.NewBuffer(make([]byte, 6))
	sb.Reset()
	sb.WriteString("USD")
	sb.WriteString(currency)
	return sb.String()
}

// return base currency, quote currency
// e.g. "XAUEUR" -> ("XAU", "EUR")
func figureCurrency(symbolName string) (string, string) {
	l := len(symbolName)
	if l != 6 {
		panic("symbol name must be 6 characters long:" + symbolName)
	}
	return symbolName[:3], symbolName[3:]
}

func calMarginLevel(equity, marginReq float64) float64 {
	if equity < 0 {
		return equity
	}
	if equity == 0 {
		return 0
	}
	if marginReq <= 0 {
		return MaxMarginLevel
	}
	marginLevel := equity / marginReq
	if marginLevel > MaxMarginLevel {
		return MaxMarginLevel
	}
	return marginLevel
}
