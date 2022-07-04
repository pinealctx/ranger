package cmd

import (
	"github.com/pinealctx/ranger/walker/walker"
	"github.com/urfave/cli/v2"
)

var (
	PFilterCmd = &cli.Command{
		Name:   "pf",
		Usage:  "filter all paths contain pattern",
		Action: pFilterAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "p",
				Usage: "file filter pattern",
			},
		},
	}
)

func pFilterAction(c *cli.Context) error {
	var fn = func() walker.Walker {
		return walker.NewFilterPathWalker(c.String("p"))
	}
	return outputPath(fn)
}
