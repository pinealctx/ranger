package jda

import (
	"fmt"
	"math"
)

const (
	firstPos  = 0
	sameDire  = 1
	closeDire = 2
)

type TradeAudit struct {
	beforeCliRpt *ClientReport
	beforeAccPos *CustomerAccountPositions

	afterCliRpt *ClientReport
	afterAccPos *CustomerAccountPositions

	tradeRst  int
	symbolCtx *SymbolCtx
	oldPrice  float64

	lastQty   NumericString
	lastPrice NumericString

	symbol  string
	ordSide string
}

type Quote struct {
	Value float64
	Rate  float64
}

func NewQuote(value, rate float64) *Quote {
	return &Quote{
		Value: value,
		Rate:  rate,
	}
}

func NewTradeAudit(beforeRpt *ClientReport, beforeAccPos *CustomerAccountPositions,
	afterRpt *ClientReport, afterAccPos *CustomerAccountPositions) *TradeAudit {
	x := &TradeAudit{
		beforeCliRpt: beforeRpt,
		beforeAccPos: beforeAccPos,
		afterCliRpt:  afterRpt,
		afterAccPos:  afterAccPos,
	}

	x.lastQty = x.beforeCliRpt.ExecutionReportInternal.LastQty
	x.lastPrice = x.beforeCliRpt.ExecutionReportInternal.LastPx

	x.symbol = x.beforeCliRpt.NewOrderSingle.Symbol
	x.ordSide = x.beforeCliRpt.NewOrderSingle.Side
	x.figureTradeResult()
	cc := NewRuntime(x.afterAccPos)
	x.symbolCtx = cc.buildSymbolCtxByName(x.symbol)
	return x
}

func (x *TradeAudit) Validate() bool {
	switch x.tradeRst {
	case firstPos:
		if !x.validateNewPosition() {
			fmt.Printf("new position validation failed for symbol %s\n", x.symbol)
			return false
		}
		commission, ok := x.validateCommission()
		if !ok {
			fmt.Printf("new position commission validation failed for symbol %s\n", x.symbol)
			return false
		}
		if !x.validateBalance(-commission) {
			fmt.Printf("new position balance validation failed for symbol %s\n", x.symbol)
			return false
		}
		fmt.Printf("new position validation succeeded for symbol %s\n", x.symbol)
	case sameDire:
		if !x.validateSamePosition() {
			fmt.Printf("same position validation failed for symbol %s\n", x.symbol)
			return false
		}
		commission, ok := x.validateCommission()
		if !ok {
			fmt.Printf("same position commission validation failed for symbol %s\n", x.symbol)
			return false
		}
		if !x.validateBalance(-commission) {
			fmt.Printf("same position balance validation failed for symbol %s\n", x.symbol)
			return false
		}
		fmt.Printf("same position validation succeeded for symbol %s\n", x.symbol)
	case closeDire:
		quote, ok := x.validateReversePosition()
		if !ok {
			fmt.Printf("close position validation failed for symbol %s\n", x.symbol)
			return false
		}
		commission, ok := x.validateCommission()
		if !ok {
			fmt.Printf("close position commission validation failed for symbol %s\n", x.symbol)
			return false
		}
		var fee float64
		if x.symbolCtx.Quote != SqQuote {
			// charge fee
			fee = 0.003 * math.Abs(quote.Value)
		}
		if !x.validateBalance(quote.Value - commission - fee) {
			fmt.Printf("close position balance validation failed for symbol %s\n", x.symbol)
			return false
		}
		fmt.Printf("close position validation succeeded for symbol %s\n", x.symbol)
	default:
		panic(fmt.Sprintf("unknown trade result: %+v", x.tradeRst))
	}
	return true
}

func (x *TradeAudit) validateNewPosition() bool {
	pos, ok := x.afterAccPos.PositionMap[x.symbol]
	if !ok {
		fmt.Printf("position does not exist for symbol %s\n", x.symbol)
		return false
	}

	if pos.Type != x.ordSide {
		fmt.Printf("position type mismatch for symbol %s: exp: %+v, pos: %+v\n",
			x.symbol, x.ordSide, pos.Type)
		return false
	}

	if pos.Qty != x.lastQty {
		fmt.Printf("position quantity mismatch for symbol %s: exp: %+v, pos: %+v\n",
			x.symbol, x.lastQty, pos.Qty)
		return false
	}

	lPrice := x.lastPrice.F64()
	if pos.Price != lPrice {
		fmt.Printf("position price mismatch for symbol %s: exp: %+v, pos: %+v\n",
			x.symbol, lPrice, pos.Price)
		return false
	}
	return true
}

func (x *TradeAudit) validateSamePosition() bool {
	posOld, ok := x.beforeAccPos.PositionMap[x.symbol]
	if !ok {
		fmt.Printf("old position does not exist for symbol %s\n", x.symbol)
		return false
	}
	posNew, ok := x.afterAccPos.PositionMap[x.symbol]
	if !ok {
		fmt.Printf("new position does not exist for symbol %s\n", x.symbol)
		return false
	}

	qty := posOld.Qty + x.lastQty
	if qty != posNew.Qty {
		fmt.Printf("position quantity mismatch for symbol %s: old: %+v, lastQty: %+v, new: %+v\n",
			x.symbol, posOld.Qty, x.lastQty, posNew.Qty)
		return false
	}

	lastPrice, lastQty := x.lastPrice.F64(), x.lastQty.F64()
	oldQty := posOld.Qty.F64()
	price := (posOld.Price*oldQty + lastPrice*lastQty) / (oldQty + lastQty)
	if posNew.Price != price {
		fmt.Printf("position price mismatch for symbol %s: old: %+v, lastPrice: %+v, new: %+v\n",
			x.symbol, posOld.Price, lastPrice, posNew.Price)
		return false
	}
	if posNew.Type != posOld.Type {
		fmt.Printf("position type mismatch for symbol %s: old: %+v, new: %+v\n",
			x.symbol, posOld.Type, posNew.Type)
		return false
	}
	if posNew.Type != x.beforeCliRpt.NewOrderSingle.Side {
		fmt.Printf("position type mismatch for symbol %s: exp: %+v, pos: %+v\n",
			x.symbol, x.beforeCliRpt.NewOrderSingle.Side, posNew.Type)
		return false
	}
	return true
}

func (x *TradeAudit) validateReversePosition() (*Quote, bool) {
	posOld, ok := x.beforeAccPos.PositionMap[x.symbol]
	if !ok {
		fmt.Printf("old position does not exist for symbol %s\n", x.symbol)
		return nil, false
	}
	x.oldPrice = posOld.Price
	posNew, ok := x.afterAccPos.PositionMap[x.symbol]
	if !ok {
		fmt.Printf("new position does not exist for symbol %s\n", x.symbol)
		return nil, false
	}

	var leftType string
	var leftPrice NumericString
	var leftQty NumericString
	var closeQuotePnl, closeQuoteRate float64
	if x.lastQty > posOld.Qty {
		leftType = x.ordSide
		leftPrice = x.lastPrice
		leftQty = x.lastQty - posOld.Qty
		closeQuotePnl, closeQuoteRate = x.figureClosePnl(posOld.Qty.F64())
	} else if x.lastQty < posOld.Qty {
		leftType = posOld.Type
		leftPrice = NumericString(posOld.Price)
		leftQty = posOld.Qty - x.lastQty
		closeQuotePnl, closeQuoteRate = x.figureClosePnl(x.lastQty.F64())
	} else {
		leftQty = 0
		if posNew.Qty == 0 {
			closeQuotePnl, closeQuoteRate = x.figureClosePnl(posOld.Qty.F64())
			return NewQuote(closeQuotePnl, closeQuoteRate), true // no position left, valid reverse
		}
		fmt.Printf("position quantity mismatch for symbol %s: old: %+v, old: %+v\n",
			x.symbol, posOld.Qty, x.lastQty)
		return nil, false
	}

	if posNew.Price != leftPrice.F64() {
		fmt.Printf("position price mismatch for symbol %s: exp: %+v, pos: %+v\n",
			x.symbol, leftPrice, posNew.Price)
		return nil, false
	}
	if posNew.Qty != leftQty {
		fmt.Printf("position quantity mismatch for symbol %s: exp: %+v, pos: %+v\n",
			x.symbol, leftQty, posNew.Qty)
		return nil, false
	}
	if posNew.Type != leftType {
		fmt.Printf("position type mismatch for symbol %s: exp: %+v, pos: %+v\n",
			x.symbol, leftType, posNew.Type)
		return nil, false
	}
	return NewQuote(closeQuotePnl, closeQuoteRate), true
}

func (x *TradeAudit) validateCommission() (float64, bool) {
	exp := x.afterCliRpt.ExecutionReportInternal
	if exp.CommissionRate == 0 {
		return 0, true // no commission, valid case
	}
	commission, quoteRate := x.figureCommission()
	if !equalsUsd(commission, exp.Commission) {
		fmt.Printf("commission mismatch for symbol %s: exp: %+v, cal: %+v\n",
			x.symbol, exp.Commission, commission)
		return 0, false
	}
	if !equalsF64(quoteRate, exp.UsdRate) {
		fmt.Printf("USD rate mismatch for symbol %s: exp: %+v, cal: %+v\n",
			x.symbol, exp.UsdRate, quoteRate)
		return 0, false
	}
	return commission, true
}

func (x *TradeAudit) validateBalance(balanceDelta float64) bool {
	_balanceDelta := x.afterAccPos.Balance - x.beforeAccPos.Balance
	if !equalsUsd(_balanceDelta, balanceDelta) {
		fmt.Printf("balance delta mismatch for symbol %s: exp: %+v, cal: %+v\n",
			x.symbol, _balanceDelta, balanceDelta)
		return false
	}
	return true
}

func (x *TradeAudit) figureClosePnl(closeQty float64) (float64, float64) {
	lastPrice := x.lastPrice.F64()

	var pnl float64
	if x.ordSide == "BUY" {
		pnl = x.oldPrice*closeQty - lastPrice*closeQty
	} else {
		pnl = lastPrice*closeQty - x.oldPrice*closeQty
	}
	return x.symbolCtx.quoteAsUsd(pnl)
}

func (x *TradeAudit) figureCommission() (float64, float64) {
	exp := x.afterCliRpt.ExecutionReportInternal
	tradeValue := exp.LastPx.F64() * exp.LastQty.F64()
	quoteValue, quoteRate := x.symbolCtx.quoteAvgUsd(tradeValue)
	commission := (exp.CommissionRate * quoteValue) / 1000000
	if commission < 0.01 {
		commission = 0.01
	}
	return commission, quoteRate
}

func (x *TradeAudit) figureTradeResult() {
	pos, ok := x.beforeAccPos.PositionMap[x.symbol]
	if !ok {
		x.tradeRst = firstPos
		return
	}
	if pos.Qty == 0 {
		x.tradeRst = firstPos
		return
	}
	if pos.Qty < 0 {
		panic(fmt.Sprintf("negative position symbol:%+v, qty: %+v", x.symbol, pos.Qty))
	}

	if pos.Type == x.ordSide {
		x.tradeRst = sameDire
		return
	}
	x.tradeRst = closeDire
}
