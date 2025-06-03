package jda

// MarketDataLiquidity 表示市场流动性信息
type MarketDataLiquidity struct {
	Px  NumericString `json:"px"`
	Qty NumericString `json:"qty"`
}

// MarketData 表示市场数据信息
type MarketData struct {
	Liquidity []MarketDataLiquidity `json:"liquidity"`
}

// InternalExecution 表示内部执行详情
type InternalExecution struct {
	AvgPx        NumericString `json:"avgPx"`
	AvgPxRcvd    NumericString `json:"avgPxRcvd"`
	ExpectedPx   NumericString `json:"expectedPx"`
	Slippage     NumericString `json:"slippage"`
	CumQty       NumericString `json:"cumQty"`
	LeavesQty    NumericString `json:"leavesQty"`
	OrderQty     NumericString `json:"orderQty"`
	Account      string        `json:"account"`
	Modified     bool          `json:"modified"`
	OrdStatus    string        `json:"ordStatus"`
	ClOrdId      string        `json:"clOrdId"`
	TimeInForce  string        `json:"timeInForce"`
	OrderID      string        `json:"orderID"`
	LastPrice    NumericString `json:"lastPrice"`
	LastQty      NumericString `json:"lastQty"`
	MarketData   MarketData    `json:"marketData"`
	Ask          NumericString `json:"ask"`
	Bid          NumericString `json:"bid"`
	CreateTime   TimeNano      `json:"createTime"`
	TransactTime TimeNano      `json:"transactTime"`
	GatewayTime  TimeNano      `json:"gatewayTime"`
	ReceivedTime TimeNano      `json:"receivedTime"`
}

// NewOrderSingle 表示新订单信息
type NewOrderSingle struct {
	EventTime      TimeNano      `json:"eventTime"`
	Account        string        `json:"account"`
	ClOrdID        string        `json:"clOrdID"`
	SendingTime    TimeNano      `json:"sendingTime"`
	Symbol         string        `json:"symbol"`
	OrdType        string        `json:"ordType"`
	ExecInst       string        `json:"execInst"`
	TimeInForce    string        `json:"timeInForce"`
	Side           string        `json:"side"`
	Price          NumericString `json:"price"`
	OrderQty       NumericString `json:"orderQty"`
	PositionEffect string        `json:"positionEffect"`
	PartyID        string        `json:"partyID"`
	ExpireDate     TimeNano      `json:"expireDate"`
	TransactTime   TimeNano      `json:"transactTime"`
}

// ExecutionReportInternal 表示执行报告内部信息
type ExecutionReportInternal struct {
	EventTime      TimeNano      `json:"eventTime"`
	ClOrdID        string        `json:"clOrdID"`
	OrigClOrdId    string        `json:"origClOrdId"`
	Text           string        `json:"text"`
	PartyId        string        `json:"partyId"`
	SubPartyId     string        `json:"subPartyId"`
	ExecID         string        `json:"execID"`
	OrderID        string        `json:"orderID"`
	Account        string        `json:"account"`
	PositionEffect string        `json:"positionEffect"`
	OrdType        string        `json:"ordType"`
	ExecType       string        `json:"execType"`
	OrdStatus      string        `json:"ordStatus"`
	Symbol         string        `json:"symbol"`
	Side           string        `json:"side"`
	TransactTime   TimeNano      `json:"transactTime"`
	SendingTime    TimeNano      `json:"sendingTime"`
	Commission     float64       `json:"commission"`
	UsdRate        float64       `json:"usdRate"`
	CommissionRate float64       `json:"commissionRate"`
	Fee            float64       `json:"fee"`
	ExchangeCcy    string        `json:"exchangeCcy"`
	Price          NumericString `json:"price"`
	AvgPx          NumericString `json:"avgPx"`
	LastPx         NumericString `json:"lastPx"`
	LastQty        NumericString `json:"lastQty"`
	LeavesQty      NumericString `json:"leavesQty"`
	MinQty         NumericString `json:"minQty"`
	CumQty         NumericString `json:"cumQty"`
	OrderQty       NumericString `json:"orderQty"`
}

// ClientReport 表示客户交易报告的完整结构
type ClientReport struct {
	EventTime               TimeNano                `json:"eventTime"`
	NewOrderSingle          NewOrderSingle          `json:"newOrderSingle"`
	ExecutionReportInternal ExecutionReportInternal `json:"executionReportInternal"`
	Slippage                NumericString           `json:"slippage"`
	ExpectedPx              NumericString           `json:"expectedPx"`
	Internal                []InternalExecution     `json:"internal"`
	Ask                     NumericString           `json:"ask"`
	Bid                     NumericString           `json:"bid"`
}
