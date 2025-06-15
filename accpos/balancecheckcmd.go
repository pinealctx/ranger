package main

import (
	"fmt"
	"github.com/pinealctx/ranger/accpos/jda"
	"github.com/urfave/cli/v2"
)

var (
	balanceCheckCmd = &cli.Command{
		Name:  "bck",
		Usage: "scan balance check from account positions",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "input",
				Usage: "input file containing account positions in JSON format",
			},
		},
		Action: balanceCheckAction,
	}
)

func balanceCheckAction(c *cli.Context) error {
	inputFile := c.String("input")
	if inputFile == "" {
		return fmt.Errorf("input file is required")
	}
	return checkBalance(inputFile)
}

func checkBalance(inputFile string) error {
	file, scanner, err := getFileLineScaner(inputFile)
	if err != nil {
		return fmt.Errorf("can't open file %s: %w", inputFile, err)
	}
	defer func() {
		_ = file.Close()
	}()

	var prvNode *jda.MixCliRptPosition
	var curNode *jda.MixCliRptPosition
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum += 1

		curNode, err = jda.FigureCliRptAndAcc(line)
		if err != nil {
			return fmt.Errorf("failed to figure client report and account at line %d: %w", lineNum, err)
		}
		if curNode != nil && prvNode != nil && !curNode.Before && prvNode.Before {
			audit := jda.NewTradeAudit(prvNode.CliReport, prvNode.AccPos, curNode.CliReport, curNode.AccPos)
			validate := audit.Validate()
			if !validate {
				fmt.Printf("balance check failed at line %d: %+v\n\n", lineNum, line)
				break
			} else {
				fmt.Printf("balance check passed at line %d\n\n", lineNum)
			}
		}
		if curNode != nil {
			prvNode = curNode
		}
	}
	return nil
}
