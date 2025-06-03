package main

import (
	"fmt"
	"github.com/pinealctx/ranger/fxb/b85"
	"github.com/urfave/cli/v2"
	"os"
)

var encCmd = &cli.Command{
	Name:   "enc",
	Usage:  "encode string to long",
	Action: encAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name: "s",
		},
	},
}

var decCmd = &cli.Command{
	Name:   "dec",
	Usage:  "decode long to string",
	Action: decAction,
	Flags: []cli.Flag{
		&cli.Int64Flag{
			Name: "i",
		},
	},
}

var app = &cli.App{
	Name:     "fxb",
	Usage:    "encode and decode string to long",
	Commands: cli.Commands{encCmd, decCmd},
}

var base85Coder = b85.NewBase85LongConverterS()

func encAction(c *cli.Context) error {
	s := c.String("s")
	fmt.Println("encode:", s)
	i, err := base85Coder.Parse(s)
	if err != nil {
		return err
	}
	fmt.Println("result:", i)
	return nil
}

func decAction(c *cli.Context) error {
	i := c.Int64("i")
	fmt.Println("decode:", i)
	fmt.Println("result:", base85Coder.AsString(i))
	return nil
}

func main() {
	var err = app.Run(os.Args)
	if err != nil {
		fmt.Println("run error:", err)
	}
}
