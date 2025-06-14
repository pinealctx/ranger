package jda

import (
	"fmt"
	"testing"
)

func TestFigureNewOrder(t *testing.T) {
	var s = `2025-06-13T02:55:49.949Z 02:55:49.949 [main/ServiceStarter/core-event-loop]  INFO  com.xsyphon.csmanager.service.CustomerManagerService - NewOrderSingle Pre check: ts:1749783349949140991, eventTime:1749783349946959781, transactTime:1749783349649000000, sendingTime:1749783349946937968, timeout ns:600000000000, timeout reject: false, event:{"eventTime":"2025-06-13T02:55:49.946959781","account":"Test_TMAC","clOrdID":"17497833499166908","sendingTime":"2025-06-13T02:55:49.946937968","symbol":"XAUUSD","ordType":"MARKET","execInst":"NONE","timeInForce":"IMMEDIATE_OR_CANCEL","side":"SELL","price":"0.00","orderQty":"1K","positionEffect":"O","partyID":"trader-1749783349649","transactTime":"2025-06-13T02:55:49.649"}`
	order, err := FigureNewOrder(s)
	if err != nil {
		t.Fatalf("Failed to figure new order: %v", err)
	}
	fmt.Printf("order: %v\n", order)
}
