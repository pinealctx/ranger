package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"strconv"
	"time"
)

var (
	convertNsCmd = &cli.Command{
		Name:   "n",
		Usage:  "convert nano unix timestamp to UTC time",
		Action: unixNanoToUTCAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "t",
			},
		},
	}

	parseUnixNanoCmd = &cli.Command{
		Name:   "p",
		Usage:  "parse UTC time to nano unix timestamp",
		Action: parseUTCtoUnixNanoAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name: "t",
			},
		},
	}

	convertNYTimeCmd = &cli.Command{
		Name:   "ny",
		Usage:  "convert a specific time string of NewYork timezone to UTC time",
		Action: convertNYTimeAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "t",
				Usage: `Specify a New York time, like"2025-07-01T13:00:00""`,
			},
		},
	}
)

func main() {
	var app = &cli.App{
		Name:  "unix nano timestamp tool",
		Usage: "convert unix timestamp to UTC time",
		Commands: cli.Commands{
			convertNsCmd,
			parseUnixNanoCmd,
			convertNYTimeCmd,
		},
	}
	var err = app.Run(os.Args)
	if err != nil {
		fmt.Println("run error:", err)
	}
}

// input nano unix timestamp, output time in UTC
func unixNanoToUTCAction(c *cli.Context) error {
	var t = c.String("t")
	if t == "" {
		return fmt.Errorf("empty timestamp")
	}
	// parse string to int64
	var ts, err = strconv.ParseInt(t, 10, 64)
	if err != nil {
		return err
	}
	// convert to time
	ut := time.Unix(0, ts).UTC()
	fmt.Println("UTC time:", ut)
	return nil
}

// input time in UTC, output nano unix timestamp
func parseUTCtoUnixNanoAction(c *cli.Context) error {
	var t = c.String("t")
	if t == "" {
		return fmt.Errorf("empty timestamp")
	}
	var ut, err = time.Parse(time.DateTime, t)
	if err != nil {
		return err
	}
	// convert to nano unix timestamp
	ts := ut.UnixNano()
	fmt.Println("nano unix timestamp:", ts)
	return nil
}

// convert a NewYork time to UTC time
func convertNYTimeAction(c *cli.Context) error {
	var t = c.String("t")
	if t == "" {
		return fmt.Errorf("empty time string")
	}
	timezone, err := time.LoadLocation("America/New_York")
	if err != nil {
		return err
	}
	tm, err := time.ParseInLocation("2006-01-02T15:04:05", t, timezone)
	if err != nil {
		return err
	}
	// convert to UTC
	ut := tm.UTC()
	fmt.Printf("New York time is %+v, UTC time is %+v\n", tm, ut)
	return nil
}
