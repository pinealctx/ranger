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
)

func main() {
	var app = &cli.App{
		Name:  "unix nano timestamp tool",
		Usage: "convert unix timestamp to UTC time",
		Commands: cli.Commands{
			convertNsCmd,
			parseUnixNanoCmd,
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
