package cmd

import (
	"github.com/pinealctx/ranger/walker/walker"
	"github.com/urfave/cli/v2"
)

var (
	FilterCmd = &cli.Command{
		Name:   "filter",
		Usage:  "filter all file paths match pattern",
		Action: filterAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "p",
				Usage: "file filter pattern",
			},
		},
	}
)

func filterAction(c *cli.Context) error {
	var fn = func() walker.Walker {
		return walker.NewFilterWalker(c.String("p"))
	}
	return outputPath(fn)
}
