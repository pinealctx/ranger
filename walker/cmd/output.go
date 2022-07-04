package cmd

import (
	"github.com/pinealctx/ranger/walker/walker"
	"log"
	"os"
	"path/filepath"
)

func outputPath(wGen func() walker.Walker) error {
	var cwd, err = os.Getwd()
	if err != nil {
		return err
	}
	log.Println("current directory:", cwd)
	var lps = wGen()
	err = filepath.Walk(cwd, lps.Walk)
	if err != nil {
		return err
	}
	var ps = lps.Paths()
	for _, p := range ps {
		log.Println(p)
	}
	return nil
}
