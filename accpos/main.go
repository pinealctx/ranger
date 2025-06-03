package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"regexp"
)

var (
	filterAccCmd = &cli.Command{
		Name:   "filter-acc",
		Usage:  "filter account data",
		Action: filterAccAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "input",
				Usage: "log file to filter",
			},
			&cli.StringFlag{
				Name:  "account",
				Usage: "filter by account ID, if not provided, all accounts will be filtered",
				Value: "", // Default to empty string to filter all accounts
			},
		},
	}
)

func main() {
	var app = &cli.App{
		Name:  "accpos",
		Usage: "filter account data from log file",
		Commands: cli.Commands{
			filterAccCmd,
		},
	}
	var err = app.Run(os.Args)
	if err != nil {
		println("run error:", err.Error())
	}
}

func filterAccAction(c *cli.Context) error {
	inputFile := c.String("input")
	if inputFile == "" {
		return fmt.Errorf("input file is required")
	}
	return filterAccData(inputFile, c.String("account"))
}

func filterAccData(inputFile string, account string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	// 模式1: outAccountPositions for account:FTWW266
	pattern := `\{"eventTime":"[^"]+","exchange":"[^"]+","positionMode":"[^"]+","accountId":"([^"]+)".*\}`

	re1 := regexp.MustCompile(pattern)

	// 创建扫描器
	scanner := bufio.NewScanner(file)

	// 设置更大的缓冲区，因为JSON可能很长
	const maxCapacity = 4 * 1024 * 1024 // 4MB
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	// 逐行处理日志
	for scanner.Scan() {
		line := scanner.Text()

		// 使用正则表达式匹配包含 "outAccountPositions for account:" 的行
		_, jsonData, match := matchAccountAndPosition(line, account, re1)
		//fmt.Printf("%s - %v -%s\n", accountId, match, jsonData)
		if match {
			fmt.Printf("%s\n", jsonData)
		}
	}
	return scanner.Err()
}

func matchAccountAndPosition(line string, account string, re *regexp.Regexp) (string, string, bool) {

	matches := re.FindStringSubmatch(line)
	if len(matches) >= 2 {
		// matches[0] 是完整的匹配（整个JSON）
		// matches[1] 是第一个捕获组（accountId 的值）
		jsonData := matches[0]
		accountID := matches[1]
		if account == "" || account == accountID {
			return accountID, jsonData, true
		}
		return accountID, jsonData, false
	}
	return "", "", false
}
