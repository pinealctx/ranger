package main

import (
	"encoding/json"
	"fmt"
	"github.com/pinealctx/ranger/accpos/jda"
	"github.com/urfave/cli/v2"
)

var (
	checkAccPosCmd = &cli.Command{
		Name:  "check-acc",
		Usage: "scan/check account positions from json file",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "input",
				Usage: "input file containing account positions in JSON format",
			},
		},
		Action: checkAccPosAction,
	}
)

func checkAccPosAction(c *cli.Context) error {
	inputFile := c.String("input")
	if inputFile == "" {
		return cli.Exit("input file is required", 1)
	}
	// Call the function to check account positions
	return checkAccPositions(inputFile)
}

func checkAccPositions(inputFile string) error {
	count := 0
	errCount := 0

	err := scanFileAsLine(inputFile, func(line string) error {
		var err error
		var acc = &jda.CustomerAccountPositions{}
		if err = json.Unmarshal([]byte(line), acc); err != nil {
			return err
		}

		count += 1
		if len(acc.SymbolSnapMap) > 0 {
			r := jda.NewRuntime(acc)
			defer func() {
				e := recover()
				if e != nil {
					fmt.Printf("panic recovered: %v\n\n", e)
					fmt.Printf("Account positions: %+v\n\n", line)
					errCount += 1
				}
			}()

			match := r.CheckWhole()
			if !match {
				fmt.Printf("Account positions mismatch: %+v\n\n", line)
				errCount += 1
			}
		}
		return nil
	})
	fmt.Printf("error count: %d, total accounts: %d, accurate:%+v\n",
		errCount, count, float64(count-errCount)/float64(count))
	return err
}
