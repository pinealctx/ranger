package main

import (
	"fmt"
	"time"
)

var (
	nyLocation   *time.Location
	beginTime    time.Time
	endTime      time.Time
	transactions []Transaction
)

func init() {
	var err error
	// Initialize the nyLocation variable with the New York time zone
	nyLocation, err = time.LoadLocation("America/New_York")
	if err != nil {
		panic(err) // If loading the location fails, panic to indicate a critical error
	}
	beginTime = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime = time.Date(2036, 1, 1, 0, 0, 0, 0, time.UTC)
	transactions = make([]Transaction, 0, 64)
}

func main() {
	// Call the function to figure out the transactions
	figureTransactions()
	// Print the transactions
	printTransactions()
}

type Transaction struct {
	T      time.Time
	IsDst  bool // IsDst indicates whether the time is in Daylight Saving Time
	Offset int  // Offset is the time zone offset in seconds
}

func figureTransactions() {
	interval := time.Hour
	var prevOffset int
	var curOffset int
	var nyTime time.Time
	var minTime = beginTime.Add(interval)
	_, prevOffset = beginTime.In(nyLocation).Zone()

	for t := minTime; t.Before(endTime); t = t.Add(interval) {
		nyTime = t.In(nyLocation)
		_, curOffset = nyTime.Zone()

		if curOffset != prevOffset {
			// If the offset has changed, we have a new transaction
			transactions = append(transactions, Transaction{
				T:      t,
				IsDst:  nyTime.IsDST(),
				Offset: curOffset,
			})
		}
		prevOffset = curOffset
	}
}

func printTransactions() {
	for _, tx := range transactions {
		weekday := tx.T.Weekday()
		weekIdx := weekdayOccurInMonth(tx.T)
		if tx.IsDst {
			fmt.Printf("DST begin: UTC:%+v, NY:%+v, week:%+v, weekday:%+v, Unix:%+v, Offset:%+v\n",
				tx.T.UTC(), tx.T.In(nyLocation), weekIdx, weekday, tx.T.Unix(), tx.Offset)
		} else {
			fmt.Printf("DST end: UTC:%+v, NY:%+v, week:%+v, weekday:%+v, Unix:%+v, Offset:%+v\n",
				tx.T.UTC(), tx.T.In(nyLocation), weekIdx, weekday, tx.T.Unix(), tx.Offset)
		}
	}

	for _, tx := range transactions {
		weekday := tx.T.Weekday()
		weekIdx := weekdayOccurInMonth(tx.T)
		if tx.IsDst {
			fmt.Printf("// DST begin: UTC:%+v, NY:%+v, week:%+v, weekday:%+v\n",
				tx.T.UTC(), tx.T.In(nyLocation), weekIdx, weekday)
			fmt.Printf("{%+vL, %+v},\n", tx.T.Unix(), tx.Offset)
		} else {
			fmt.Printf("// DST end: UTC:%+v, NY:%+v, week:%+v, weekday:%+v\n",
				tx.T.UTC(), tx.T.In(nyLocation), weekIdx, weekday)
			fmt.Printf("{%+vL, %+v},\n", tx.T.Unix(), tx.Offset)
		}
	}
}

func weekdayOccurInMonth(t time.Time) int {
	firstDayOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	weekdayDiff := int(t.Weekday() - firstDayOfMonth.Weekday())
	if weekdayDiff < 0 {
		weekdayDiff += 7
	}
	firstSameWeekday := firstDayOfMonth.AddDate(0, 0, weekdayDiff)
	dayDiff := t.Day() - firstSameWeekday.Day()
	weekDiff := dayDiff / 7
	if dayDiff%7 != 0 {
		weekDiff++
	}
	return weekDiff + 1
}
