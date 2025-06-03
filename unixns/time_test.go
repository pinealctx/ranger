package main

import (
	"fmt"
	"testing"
	"time"
)

func TestTs(t *testing.T) {
	anchor := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)
	ns := anchor.UnixNano()
	us := ns / 1000
	ms := ns / 1000000
	fmt.Printf("anchor: %v, ns: %v, us: %v, ms: %v\n", anchor, ns, us, ms)

	fmt.Println("----------convert as seconds----------")
	fmt.Println(convertSecToTime(ns))
	fmt.Println(convertSecToTime(us))
	fmt.Println(convertSecToTime(ms))

	fmt.Println("----------convert as microseconds----------")
	fmt.Println(convertUsToTime(ns))
	fmt.Println(convertUsToTime(us))
	fmt.Println(convertUsToTime(ms))

	fmt.Println("----------convert as milliseconds----------")
	fmt.Println(convertMsToTime(ns))
	fmt.Println(convertMsToTime(us))
	fmt.Println(convertMsToTime(ms))

	fmt.Println("----------convert as universe timestamp----------")
	fmt.Println(convertUniverseTsToTime(ns))
	fmt.Println(convertUniverseTsToTime(us))
	fmt.Println(convertUniverseTsToTime(ms))

	anchor = time.Now()
	ns = anchor.UnixNano()
	us = ns / 1000
	ms = ns / 1000000
	fmt.Println("----------convert as universe timestamp----------")
	fmt.Println(convertUniverseTsToTime(ns))
	fmt.Println(convertUniverseTsToTime(us))
	fmt.Println(convertUniverseTsToTime(ms))

	anchor = time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
	ns = anchor.UnixNano()
	us = ns / 1000
	ms = ns / 1000000
	fmt.Println("----------convert as universe timestamp----------")
	fmt.Println(convertUniverseTsToTime(ns))
	fmt.Println(convertUniverseTsToTime(us))
	fmt.Println(convertUniverseTsToTime(ms))
}

func convertSecToTime(sec int64) time.Time {
	return time.Unix(sec, 0).UTC()
}

func convertUsToTime(us int64) time.Time {
	return time.Unix(0, us*1000).UTC()
}

func convertMsToTime(ms int64) time.Time {
	return time.Unix(0, ms*1000000).UTC()
}

func convertUniverseTsToTime(u int64) time.Time {
	if u >= 1262304000000000000 {
		return time.Unix(0, u).UTC()
	} else if u >= 1262304000000000 {
		return convertUsToTime(u)
	} else if u >= 1262304000000 {
		return convertMsToTime(u)
	} else {
		return convertSecToTime(u)
	}
}
