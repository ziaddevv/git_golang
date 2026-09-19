package cli

import (
	"mygit/internal/porcelain"
	"os"

	flag "github.com/spf13/pflag"
)

func LogCommand() error {
	cmd := flag.NewFlagSet("log", flag.ExitOnError)

	oneLine := cmd.Bool("oneline", false, "the commit logs")

	cmd.Parse(os.Args[2:])

	if *oneLine == false {
		porcelain.Log()
		// procelain.CommitLog()
	} else {
	}

	return nil
}
