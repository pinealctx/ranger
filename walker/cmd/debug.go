package cmd

import (
	"github.com/pinealctx/ranger/walker/walker"
	"github.com/urfave/cli/v2"
)

var (
	DebugCmd = &cli.Command{
		Name:   "debug",
		Usage:  "just print all file paths",
		Action: debugAction,
	}
)

func debugAction(_ *cli.Context) error {
	var fn = func() walker.Walker {
		return walker.NewDebugWalker()
	}
	return outputPath(fn)
}
