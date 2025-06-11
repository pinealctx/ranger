package main

import (
	"encoding/json"
	"fmt"
	"github.com/pinealctx/ranger/accpos/jda"
	"github.com/pinealctx/ranger/fxb/b85"
	"github.com/urfave/cli/v2"
	"regexp"
)

var (
	checkEodCmd = &cli.Command{
		Name:   "check-eod",
		Usage:  "check end of day account positions",
		Action: checkEodAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "input",
				Usage: "input eod tier price and account positions in JSON format",
			},
			&cli.StringFlag{
				Name:  "tier",
				Usage: "tier name",
			},
		},
	}

	tierPricesRegex  = regexp.MustCompile(`EOD Tier Prices sent total: \d+, detail:(\{.*?\]\})`)
	eodPositionRegex = regexp.MustCompile(`EOD for account:([^,]+), position details: (.+)$`)

	b85Conv = b85.NewBase85LongConverterS()
)

func checkEodAction(ctx *cli.Context) error {
	inputFile := ctx.String("input")
	tierName := ctx.String("tier")
	if inputFile == "" || tierName == "" {
		return cli.Exit("input file and tier are required", 1)
	}
	file, scanner, err := getFileLineScaner(inputFile)
	if err != nil {
		return cli.Exit(err.Error(), 1)
	}
	defer func() {
		_ = file.Close()
	}()

	var eodPrices *jda.TierPriceEod
	eodPositions := make([]*jda.CustomerAccountPositions, 0, 8)
	for scanner.Scan() {
		line := scanner.Text()
		if eodPrices == nil {
			eodPrices = figureEodPrice(line)
			if eodPrices == nil {
				continue
			}
		}
		eodPosition := figureEodPosition(line)
		if eodPosition != nil {
			eodPositions = append(eodPositions, eodPosition)
		}
	}

	if eodPrices == nil {
		return cli.Exit("no EOD tier prices found in the input file", 1)
	}
	if len(eodPositions) == 0 {
		return cli.Exit("no EOD account positions found in the input file", 1)
	}

	checkEodData(tierName, eodPrices, eodPositions)
	return nil
}

func checkEodData(tierName string, eodPrices *jda.TierPriceEod, eodPositions []*jda.CustomerAccountPositions) {
	m := make(map[int64]*jda.TierPrice)
	for _, price := range eodPrices.TierPrices {
		if price.TierId == tierName {
			idx, err := b85Conv.Parse(price.SymbolId)
			if err != nil {
				panic(err)
			}
			m[idx] = price
		}
	}

	for k, v := range m {
		fmt.Printf("%s - %+v \n", b85Conv.AsString(k), v)
	}

	var preEodTime jda.TimeNano
	for _, accPos := range eodPositions {
		for k, snap := range accPos.SymbolSnapMap {
			p := m[k.Int64()]
			if p == nil {
				fmt.Printf("EOD position mismatch for account: %s, symbol: %s\n",
					accPos.AccountId, b85Conv.AsString(k.Int64()))
				continue
			}
			if p.Bid != snap.Bid || p.Ask != snap.Ask || p.SendingTime.Int64() != snap.SendingTime {
				fmt.Printf("EOD position mismatch for account: %s, symbol: %s\n",
					accPos.AccountId, b85Conv.AsString(k.Int64()))
			}
			if preEodTime == 0 {
				preEodTime = accPos.EodTime
			} else if preEodTime != accPos.EodTime {
				fmt.Printf("EOD position mismatch for account: %s, symbol: %s\n",
					accPos.AccountId, b85Conv.AsString(k.Int64()))
			}
		}
	}

	for _, price := range eodPrices.TierPrices {
		if price.SendingTime >= preEodTime {
			fmt.Printf("EOD position mismatch for account: %s, symbol: %s\n",
				price.TierId, price.SymbolId)
		}
	}
	fmt.Printf("EOD tier prices length: %d, EOD positions length: %d\n",
		len(eodPrices.TierPrices), len(eodPositions))
}

func figureEodPrice(line string) *jda.TierPriceEod {
	return figureMatchAndJsonUnmarshal(tierPricesRegex, line, &jda.TierPriceEod{}, 2)
}

func figureEodPosition(line string) *jda.CustomerAccountPositions {
	return figureMatchAndJsonUnmarshal(eodPositionRegex, line, &jda.CustomerAccountPositions{}, 3)
}

func figureMatchAndJsonUnmarshal[T any](re *regexp.Regexp, line string, v *T, size int) *T {
	matches := re.FindStringSubmatch(line)
	if len(matches) != size {
		return nil
	}
	err := json.Unmarshal([]byte(matches[size-1]), v)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal JSON: %s, error: %v", matches[1], err))
	}
	return v
}
