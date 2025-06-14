package main

import (
	"fmt"
	"github.com/pinealctx/ranger/accpos/jda"
	"github.com/urfave/cli/v2"
)

var (
	tradeMarginCmd = &cli.Command{
		Name:  "tmargin",
		Usage: "scan trade margin from account positions",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "input",
				Usage: "input file containing account positions in JSON format",
			},
		},
		Action: tradeMarginAction,
	}
)

func tradeMarginAction(c *cli.Context) error {
	inputFile := c.String("input")
	if inputFile == "" {
		return fmt.Errorf("input file is required")
	}
	return checkTradeMargin(inputFile)
}

func checkTradeMargin(inputFile string) error {
	file, scanner, err := getFileLineScaner(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", inputFile, err)
	}
	defer func() {
		_ = file.Close()
	}()

	var orderCursor *jda.NewOrder
	var order *jda.NewOrder
	var marginAllow *jda.MarginAllow
	var accPos *jda.CustomerAccountPositions

	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum += 1
		order, err = jda.FigureNewOrder(line)
		if err != nil {
			return fmt.Errorf("failed to figure new order at line %d: %w", lineNum, err)
		}
		if order != nil {
			orderCursor = order
			continue
		}
		marginAllow, accPos, err = jda.FigureMarginAllowAndAccPos(line)
		if err != nil {
			return fmt.Errorf("failed to figure margin allow and account positions at line %d: %w", lineNum, err)
		}
		if marginAllow != nil && accPos != nil {
			marginAudit := jda.NewMarginAudit(orderCursor, marginAllow, accPos)
			validate := marginAudit.Validate()
			if !validate {
				fmt.Printf("Margin audit failed at line %d: %+v\n", lineNum, line)
			} else {
				fmt.Printf("Margin audit passed at line %d: %+v\n", lineNum, line)
			}
		}
	}
	err = scanner.Err()
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", inputFile, err)
	}
	return nil
}
