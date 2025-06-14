package main

import (
	"fmt"
	"github.com/pinealctx/ranger/accpos/jda"
	"github.com/urfave/cli/v2"
)

var (
	scanTradeCmd = &cli.Command{
		Name:  "sct",
		Usage: "scan account positions and compare them",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "input",
				Usage: "input account log file to scan",
			},
		},
		Action: scanTradeAction,
	}
)

func scanTradeAction(c *cli.Context) error {
	inputFile := c.String("input")
	if inputFile == "" {
		return cli.Exit("input file is required", 1)
	}

	accRecs, err := filterAccRecs(inputFile)
	if err != nil {
		return cli.Exit(fmt.Sprintf("failed to filter account records: %v", err), 1)
	}

	if len(accRecs) < 2 {
		return cli.Exit("not enough account records to compare", 1)
	}

	processBalanceRecCompare(accRecs)
	return nil
}

func processBalanceRecCompare(recs []AccRec) {
	size := len(recs)
	var diff jda.AccountMainDiff
	for i := 1; i < size; i++ {
		if recs[i].AccPos.Balance != recs[i-1].AccPos.Balance {
			recs[i].AccPos.MainDiffRatio(recs[i-1].AccPos, &diff)
			fmt.Printf("Line %d: %s\n", recs[i].Line, diff.String())
		}
	}
}
