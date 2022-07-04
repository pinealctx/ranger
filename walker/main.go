package main

import (
	"fmt"
	"github.com/pinealctx/ranger/walker/cmd"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	var app = &cli.App{
		Name:  "walker",
		Usage: "walk in specified directory",
		Commands: cli.Commands{
			cmd.DebugCmd,
			cmd.FilterCmd,
			cmd.PFilterCmd,
			cmd.PFilterExeCmd,
		},
	}
	var err = app.Run(os.Args)
	if err != nil {
		fmt.Println("run error:", err)
	}
}
