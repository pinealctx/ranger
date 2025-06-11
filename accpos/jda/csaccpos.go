package jda

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

// StringInt64Map 使用 StringInt64 作为键的映射类型
type StringInt64Map[T any] map[StringInt64]T

// CustomerAccountDiff Position different
type CustomerAccountDiff struct {
	MarginRatioDiff               float64 `json:"marginDiffRatio"`               // 头寸差异率
	AvailableForMarginTradingDiff float64 `json:"availableForMarginTradingDiff"` // 可用于保证金交易的差异
	MarginRequirementDiff         float64 `json:"marginRequirementDiff"`         // 保证金要求差异
	EquityDiff                    float64 `json:"equityDiff"`                    // 权益差异
	OpenPNLDiff                   float64 `json:"openPNLDiff"`                   // 未实现盈亏差异
	NotionalValueDiff             float64 `json:"notionalValueDiff"`             // 名义价值差异
}

func (x *CustomerAccountDiff) Less(ratio float64) bool {
	return math.Abs(x.MarginRatioDiff) < ratio &&
		math.Abs(x.AvailableForMarginTradingDiff) < ratio &&
		math.Abs(x.MarginRequirementDiff) < ratio &&
		math.Abs(x.EquityDiff) < ratio &&
		math.Abs(x.OpenPNLDiff) < ratio &&
		math.Abs(x.NotionalValueDiff) < ratio
}

func (x *CustomerAccountDiff) OnlyPnlNotLess(ratio float64) bool {
	// 仅检查未实现盈亏是否小于给定比例
	return math.Abs(x.MarginRatioDiff) < ratio &&
		math.Abs(x.AvailableForMarginTradingDiff) < ratio &&
		math.Abs(x.MarginRequirementDiff) < ratio &&
		math.Abs(x.EquityDiff) < ratio &&
		math.Abs(x.NotionalValueDiff) < ratio
}

func (x *CustomerAccountDiff) String() string {
	return fmt.Sprintf("marginLevel:%+v, available:%+v, marginReq:%+v, equity:%+v, openPnl:%+v, notionalValue:%+v",
		x.MarginRatioDiff, x.AvailableForMarginTradingDiff, x.MarginRequirementDiff, x.EquityDiff, x.OpenPNLDiff, x.NotionalValueDiff)
}

func (x *CustomerAccountDiff) OutRatioString(ratio float64) string {
	// if math.Abs(field) >= ratio, append it as string builder
	sb := strings.Builder{}
	if math.Abs(x.MarginRatioDiff) >= ratio {
		sb.WriteString(fmt.Sprintf("marginLevel:%+v ", x.MarginRatioDiff))
	}
	if math.Abs(x.AvailableForMarginTradingDiff) >= ratio {
		sb.WriteString(fmt.Sprintf("available:%+v ", x.AvailableForMarginTradingDiff))
	}
	if math.Abs(x.MarginRequirementDiff) >= ratio {
		sb.WriteString(fmt.Sprintf("marginReq:%+v ", x.MarginRequirementDiff))
	}
	if math.Abs(x.EquityDiff) >= ratio {
		sb.WriteString(fmt.Sprintf("equity:%+v ", x.EquityDiff))
	}
	if math.Abs(x.OpenPNLDiff) >= ratio {
		sb.WriteString(fmt.Sprintf("openPnl:%+v ", x.OpenPNLDiff))
	}
	if math.Abs(x.NotionalValueDiff) >= ratio {
		sb.WriteString(fmt.Sprintf("notionalValue:%+v ", x.NotionalValueDiff))
	}
	if sb.Len() == 0 {
		return "no significant differences"
	}
	return sb.String()
}

// CustomerAccountPositions 代表账户头寸数据的顶层结构
type CustomerAccountPositions struct {
	EventTime                 TimeNano                          `json:"eventTime"`
	Exchange                  string                            `json:"exchange"`
	PositionMode              string                            `json:"positionMode"`
	AccountId                 string                            `json:"accountId"`
	TransactionId             string                            `json:"transactionId"`
	WrittenTime               TimeNano                          `json:"writtenTime"`
	EodTime                   TimeNano                          `json:"eodTime"`
	LastEventTime             TimeNano                          `json:"lastEventTime"`
	PositionMap               map[string]*Position              `json:"positionMap"`
	SymbolSnapMap             StringInt64Map[*SymbolSnap]       `json:"symbolSnapMap"`
	IsPNLUpdate               bool                              `json:"isPNLUpdate"`
	Currency                  string                            `json:"currency"`
	Balance                   float64                           `json:"balance"`
	MarginRatio               float64                           `json:"marginRatio"`
	AvailableForMarginTrading float64                           `json:"availableForMarginTrading"`
	MarginRequirement         float64                           `json:"marginRequirement"`
	Equity                    float64                           `json:"equity"`
	Leverage                  int                               `json:"leverage"`
	OpenPNL                   float64                           `json:"openPNL"`
	ClosedPNL                 float64                           `json:"closedPNL"`
	NotionalValue             float64                           `json:"notionalValue"`
	CreditLimit               float64                           `json:"creditLimit"`
	LastBalanceTrxId          int                               `json:"lastBalanceTrxId"`
	Fee                       float64                           `json:"fee"`
	OpenOrders                interface{}                       `json:"openOrders"`
	MarginStatus              string                            `json:"marginStatus"`
	ProcessingPositionHashMap StringInt64Map[map[string]*Order] `json:"processingPositionHashMap"`
}

func (x *CustomerAccountPositions) JsonString() string {
	buf, _ := json.Marshal(x)
	return string(buf)
}

func (x *CustomerAccountPositions) DiffRatio(other *CustomerAccountPositions) *CustomerAccountDiff {
	return &CustomerAccountDiff{
		MarginRatioDiff:               diffRatio(x.MarginRatio, other.MarginRatio),
		AvailableForMarginTradingDiff: diffRatio(x.AvailableForMarginTrading, other.AvailableForMarginTrading),
		MarginRequirementDiff:         diffRatio(x.MarginRequirement, other.MarginRequirement),
		EquityDiff:                    diffRatio(x.Equity, other.Equity),
		OpenPNLDiff:                   diffRatio(x.OpenPNL, other.OpenPNL),
		NotionalValueDiff:             diffRatio(x.NotionalValue, other.NotionalValue),
	}
}

func (x *CustomerAccountPositions) SameProcessingOrdsSource(other *CustomerAccountPositions) bool {
	// assumes that x and other are not nil
	var xL1Size = len(x.ProcessingPositionHashMap)
	var otherL1Size = len(other.ProcessingPositionHashMap)
	if xL1Size != otherL1Size {
		return false
	}
	for k, v := range x.ProcessingPositionHashMap {
		// check if the key exists in the other map
		if otherV, ok := other.ProcessingPositionHashMap[k]; !ok {
			return false
		} else {
			// check if the values are the same
			if len(v) != len(otherV) {
				return false
			}
			for orderId, order := range v {
				if otherOrder, ook := otherV[orderId]; !ook || !order.SameSource(otherOrder) {
					return false
				}
			}
		}
	}
	return true
}

func (x *CustomerAccountPositions) SamePositionMapSource(other *CustomerAccountPositions) bool {
	var xL1Size = len(x.PositionMap)
	var otherL1Size = len(other.PositionMap)
	if xL1Size != otherL1Size {
		return false
	}
	for k, v := range x.PositionMap {
		// check if the key exists in the other map
		if otherV, ok := other.PositionMap[k]; !ok {
			return false
		} else {
			// check if the values are the same
			if !v.SameSource(otherV) {
				return false
			}
		}
	}
	return true
}

const openPnlPrecision = 0.015

func (x *CustomerAccountPositions) OpenPnlMatch() bool {
	totalPnl := 0.0
	for _, v := range x.PositionMap {
		totalPnl += v.QuotePnl
	}
	return math.Abs(totalPnl-x.OpenPNL) < openPnlPrecision
}

func (x *CustomerAccountPositions) SameSource(other *CustomerAccountPositions) bool {
	return x.Exchange == other.Exchange &&
		x.PositionMode == other.PositionMode &&
		x.AccountId == other.AccountId &&
		x.IsPNLUpdate == other.IsPNLUpdate &&
		x.Currency == other.Currency &&
		x.Balance == other.Balance &&
		x.Leverage == other.Leverage &&
		x.ClosedPNL == other.ClosedPNL &&
		x.CreditLimit == other.CreditLimit &&
		x.LastBalanceTrxId == other.LastBalanceTrxId &&
		x.Fee == other.Fee &&
		x.MarginStatus == other.MarginStatus
}

func (x *CustomerAccountPositions) DiffSourceString(other *CustomerAccountPositions) string {
	var sb strings.Builder
	if x.Exchange != other.Exchange {
		sb.WriteString(fmt.Sprintf("Exchange: %s -> %s; ", other.Exchange, x.Exchange))
	}
	if x.PositionMode != other.PositionMode {
		sb.WriteString(fmt.Sprintf("PositionMode: %s -> %s; ", other.PositionMode, x.PositionMode))
	}
	if x.AccountId != other.AccountId {
		sb.WriteString(fmt.Sprintf("AccountId: %s -> %s; ", other.AccountId, x.AccountId))
	}
	if x.IsPNLUpdate != other.IsPNLUpdate {
		sb.WriteString(fmt.Sprintf("IsPNLUpdate: %t -> %t; ", other.IsPNLUpdate, x.IsPNLUpdate))
	}
	if x.Currency != other.Currency {
		sb.WriteString(fmt.Sprintf("Currency: %s -> %s; ", other.Currency, x.Currency))
	}
	if x.Balance != other.Balance {
		sb.WriteString(fmt.Sprintf("Balance: %.2f, %.2f -> %.2f; ", other.Balance-x.Balance, other.Balance, x.Balance))
	}
	if x.Leverage != other.Leverage {
		sb.WriteString(fmt.Sprintf("Leverage: %d -> %d; ", other.Leverage, x.Leverage))
	}
	if x.ClosedPNL != other.ClosedPNL {
		sb.WriteString(fmt.Sprintf("ClosedPNL: %.2f -> %.2f; ", other.ClosedPNL, x.ClosedPNL))
	}
	if x.CreditLimit != other.CreditLimit {
		sb.WriteString(fmt.Sprintf("CreditLimit: %.2f -> %.2f; ", other.CreditLimit, x.CreditLimit))
	}
	if x.LastBalanceTrxId != other.LastBalanceTrxId {
		sb.WriteString(fmt.Sprintf("LastBalanceTrxId: %d -> %d; ", other.LastBalanceTrxId, x.LastBalanceTrxId))
	}
	if x.Fee != other.Fee {
		sb.WriteString(fmt.Sprintf("Fee: %.2f -> %.2f; ", other.Fee, x.Fee))
	}
	if x.MarginStatus != other.MarginStatus {
		sb.WriteString(fmt.Sprintf("MarginStatus: %s -> %s; ", other.MarginStatus, x.MarginStatus))
	}
	if sb.Len() == 0 {
		return "No significant differences"
	}
	return sb.String()
}

// Position 代表单个交易对的头寸
type Position struct {
	EventTime     TimeNano      `json:"eventTime"`
	Exchange      string        `json:"exchange"`
	PositionMode  string        `json:"positionMode"`
	AccountId     string        `json:"accountId"`
	PositionId    string        `json:"positionId"`
	TransactionId string        `json:"transactionId"`
	StrategyId    string        `json:"strategyId"`
	PlanExecId    string        `json:"planExecId"`
	StepId        string        `json:"stepId"`
	OpenTime      TimeNano      `json:"openTime"`
	Symbol        string        `json:"symbol"`
	Type          string        `json:"type"`
	Price         float64       `json:"price"`
	Qty           NumericString `json:"qty"`
	QtyAvailable  float64       `json:"qtyAvailable"`
	PriceCurr     float64       `json:"priceCurr"`
	Pnl           float64       `json:"pnl"`
	QuotePnl      float64       `json:"quotePnl"`
	QuoteValue    float64       `json:"quoteValue"`
	Swap          float64       `json:"swap"`
	Status        string        `json:"status"`
	Verified      bool          `json:"verified"`
	OriginTradeId string        `json:"originTradeId"`
	Reason        string        `json:"reason"`
	CreatedAt     TimeNano      `json:"createdAt"`
	UpdatedAt     TimeNano      `json:"updatedAt"`
}

func (x *Position) SameSource(other *Position) bool {
	// assumes that x and other are not nil
	return x.Exchange == other.Exchange &&
		x.PositionMode == other.PositionMode &&
		x.AccountId == other.AccountId &&
		x.PositionId == other.PositionId &&
		x.StrategyId == other.StrategyId &&
		x.PlanExecId == other.PlanExecId &&
		x.StepId == other.StepId &&
		x.OpenTime == other.OpenTime &&
		x.Symbol == other.Symbol &&
		x.Type == other.Type &&
		x.Price == other.Price &&
		x.Qty == other.Qty &&
		x.QtyAvailable == other.QtyAvailable &&
		x.Swap == other.Swap &&
		x.Status == other.Status &&
		x.Verified == other.Verified &&
		x.OriginTradeId == other.OriginTradeId &&
		x.Reason == other.Reason &&
		x.CreatedAt == other.CreatedAt
}

// SymbolSnap 代表交易对的快照数据
type SymbolSnap struct {
	Bid         float64 `json:"bid"`
	Ask         float64 `json:"ask"`
	Leverage    float64 `json:"leverage"`
	SendingTime int64   `json:"sendingTime"`
}

// Order 代表处理中的订单
type Order struct {
	Side         string        `json:"side"`
	OrdType      string        `json:"ordType"`
	TimeInForce  string        `json:"timeInForce"`
	ExpireDate   string        `json:"expireDate"`
	TransactTime TimeNano      `json:"transactTime"`
	EventTime    TimeNano      `json:"eventTime"`
	Qty          NumericString `json:"qty"`
	Price        NumericString `json:"price"`
	StopPrice    float64       `json:"stopPrice"`
	CumQty       float64       `json:"cumQty"`
	AvgPrice     float64       `json:"avgPrice"`
	QuoteValue   float64       `json:"quoteValue"`
}

func (x *Order) SameSource(other *Order) bool {
	// assumes that x and other are not nil
	return x.Side == other.Side &&
		x.OrdType == other.OrdType &&
		x.TimeInForce == other.TimeInForce &&
		x.ExpireDate == other.ExpireDate &&
		x.TransactTime == other.TransactTime &&
		x.EventTime == other.EventTime &&
		x.Qty == other.Qty &&
		x.Price == other.Price &&
		x.StopPrice == other.StopPrice &&
		x.CumQty == other.CumQty &&
		x.AvgPrice == other.AvgPrice
}

func diffRatio(x, y float64) float64 {
	if y == 0 {
		return (x - y) / 1 // Avoid division by zero, return a large number
	}
	return (x - y) / y
}
