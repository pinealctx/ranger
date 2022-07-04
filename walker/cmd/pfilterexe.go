package cmd

import (
	"github.com/pinealctx/ranger/walker/exe"
	"github.com/pinealctx/ranger/walker/walker"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path/filepath"
)

var (
	PFilterExeCmd = &cli.Command{
		Name:   "pfe",
		Usage:  "filter all paths contain pattern amd exec a command in such path",
		Action: pFilterExeAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "p",
				Usage: "file filter pattern",
			},
			&cli.StringFlag{
				Name:  "c",
				Usage: "command",
			},
		},
	}
)

func pFilterExeAction(c *cli.Context) error {
	var cwd, err = os.Getwd()
	if err != nil {
		return err
	}
	log.Println("current directory:", cwd)
	var lps = walker.NewFilterPathWalker(c.String("p"))
	err = filepath.Walk(cwd, lps.Walk)
	if err != nil {
		return err
	}

	var command = c.String("c")
	var ps = lps.Paths()
	for _, p := range ps {
		err = exe.RunExe(p, command)
		log.Println("run command:", command, " in dir:", p, " error:", err)
		if err != nil {
			return err
		}
	}
	return nil
}
