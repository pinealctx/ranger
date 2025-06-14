package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

var (
	exceptCmd = &cli.Command{
		Name:   "except",
		Usage:  "except the line that contains spec string",
		Action: exceptAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "input",
				Usage: "log file to except some string",
			},
			&cli.StringFlag{
				Name:  "except",
				Usage: "except the line that contains this string",
			},
		},
	}
)

func exceptAction(c *cli.Context) error {
	inputFile := c.String("input")
	if inputFile == "" {
		return cli.Exit("input file is required", 1)
	}
	exceptStr := c.String("except")
	if exceptStr == "" {
		return cli.Exit("except string is required", 1)
	}
	expLines, err := exceptLines(inputFile, exceptStr)
	if err != nil {
		return cli.Exit("failed to except lines: "+err.Error(), 1)
	}
	return writeLinesToFile(inputFile, expLines)
}

func exceptLines(inputFile, exceptStr string) ([]string, error) {
	file, scanner, err := getFileLineScaner(inputFile)
	if err != nil {
		return nil, cli.Exit("failed to open file "+inputFile+": "+err.Error(), 1)
	}
	defer func() {
		_ = file.Close()
	}()

	// collect lines, at the end, write them to file
	// old file content will be replaced
	// except the lines that contain exceptStr
	lines := make([]string, 0, 1024)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, exceptStr) {
			lines = append(lines, line)
		}
	}
	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return lines, nil
}

func writeLinesToFile(input string, lines []string) error {
	outputFile, err := os.Create(input)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", input, err)
	}
	defer func() {
		_ = outputFile.Close()
	}()

	// Write the lines to the file
	for _, line := range lines {
		if _, err = outputFile.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write line to file: %w", err)
		}
	}
	return nil
}
