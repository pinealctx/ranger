package main

import (
	"encoding/json"
	"fmt"
	"github.com/pinealctx/ranger/accpos/jda"
	"github.com/urfave/cli/v2"
)

var (
	scanAccPosCmd = &cli.Command{
		Name:   "scan-acc",
		Usage:  "scan account positions from json file",
		Action: scanAccPosAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "input",
				Usage: "input file containing account positions in JSON format",
			},
			&cli.Float64Flag{
				Name:  "threshold",
				Usage: "threshold for diff ratio, default is 0.001",
				Value: 0.001,
			},
		},
	}
)

func scanAccPosAction(c *cli.Context) error {
	inputFile := c.String("input")
	if inputFile == "" {
		return cli.Exit("input file is required", 1)
	}
	threshold := c.Float64("threshold")
	// Call the function to scan account positions
	return scanAccountPositions(inputFile, threshold)
}

func scanAccountPositions(inputFile string, threshold float64) error {
	var preAcc *jda.CustomerAccountPositions
	return scanFileAsLine(inputFile, func(line string) error {
		var err error
		var acc = &jda.CustomerAccountPositions{}
		if err = json.Unmarshal([]byte(line), acc); err != nil {
			return err
		}

		if !acc.OpenPnlMatch() {
			fmt.Printf("Account positions open PnL mismatch: %+v\n\n", acc.WrittenTime)
			fmt.Printf("%+v\n\n", acc.JsonString())
		}

		var needLog bool
		if preAcc != nil {
			if !acc.SameSource(preAcc) {
				fmt.Printf("Account positions from different sources: %+v to %+v\n", preAcc.WrittenTime, acc.WrittenTime)
				fmt.Printf("%+v\n\n", acc.DiffSourceString(preAcc))
			} else {
				diff := acc.DiffRatio(preAcc)
				needLog = false
				if !diff.Less(threshold) {
					if !diff.OnlyPnlNotLess(threshold) {
						needLog = true
					} else {
						if !preAcc.OpenPnlMatch() || !acc.OpenPnlMatch() {
							needLog = true
						}
					}
					if needLog {
						fmt.Printf("Account positions from same source: %+v -> %+v, diff: %+v\n\n",
							preAcc.WrittenTime, acc.WrittenTime, diff.OutRatioString(threshold))
						fmt.Printf("%+v\n\n", preAcc.JsonString())
						fmt.Printf("%+v\n\n", acc.JsonString())
					}
				}
			}
		}
		preAcc = acc
		return nil
	})
}
