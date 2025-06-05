package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	var app = &cli.App{
		Name:  "accpos",
		Usage: "filter account data from log file",
		Commands: cli.Commands{
			filterAccCmd,
			scanAccPosCmd,
			checkAccPosCmd,
		},
	}
	var err = app.Run(os.Args)
	if err != nil {
		println("run error:", err.Error())
	}
}
