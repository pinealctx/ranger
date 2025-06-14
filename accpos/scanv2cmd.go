package main

import (
	"encoding/json"
	"fmt"
	"github.com/pinealctx/ranger/accpos/jda"
	"github.com/urfave/cli/v2"
)

type AccRec struct {
	AccPos *jda.CustomerAccountPositions
	Line   int
}

var (
	scanV2Cmd = &cli.Command{
		Name:  "sv2",
		Usage: "scan account positions and compare them",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "input",
				Usage: "input account log file to scan",
			},
			&cli.Float64Flag{
				Name:  "kpi",
				Usage: "KPI threshold for account position comparison",
				Value: 0.001, // Default threshold value
			},
		},
		Action: scanV2Action,
	}
)

func scanV2Action(c *cli.Context) error {
	inputFile := c.String("input")
	if inputFile == "" {
		return cli.Exit("input file is required", 1)
	}
	kpi := c.Float64("kpi")

	accRecs, err := filterAccRecs(inputFile)
	if err != nil {
		return cli.Exit(fmt.Sprintf("failed to filter account records: %v", err), 1)
	}

	if len(accRecs) < 2 {
		return cli.Exit("not enough account records to compare", 1)
	}

	processAccRecCompare(accRecs, kpi)
	return nil
}

func filterAccRecs(inputFile string) ([]AccRec, error) {
	file, scanner, err := getFileLineScaner(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", inputFile, err)
	}
	defer func() {
		_ = file.Close()
	}()

	accRecs := make([]AccRec, 0, 512)
	// 逐行处理
	lineNum := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineNum += 1
		_, jsonData, match := matchAccountAndPosition(line, "", accPosReg)
		if match {
			accData := jda.CustomerAccountPositions{}
			err = json.Unmarshal([]byte(jsonData), &accData)
			if err != nil {
				return nil, fmt.Errorf("error unmarshalling JSON on line %d: %w", lineNum, err)
			}
			accRecs = append(accRecs, AccRec{
				AccPos: &accData,
				Line:   lineNum,
			})
		}
	}
	return accRecs, scanner.Err()
}

func processAccRecCompare(recs []AccRec, kpi float64) {
	size := len(recs)
	var diff jda.AccountMainDiff
	for i := 1; i < size; i++ {
		recs[i].AccPos.MainDiffRatio(recs[i-1].AccPos, &diff)
		if !diff.IsLess(kpi) {
			fmt.Printf("Line %d: %s\n", recs[i].Line, diff.String())
		}
	}
}
