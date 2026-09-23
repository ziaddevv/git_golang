package cli

import (
	"fmt"
	"os"

	"mygit/internal/porcelain"

	flag "github.com/spf13/pflag"
)

func AddCommand() error {
	cmd := flag.NewFlagSet("add", flag.ExitOnError)

	cmd.Parse(os.Args[2:])

	if cmd.NArg() < 1 {
		return fmt.Errorf("usage: mygit add <file>...")
	}

	return porcelain.Add(cmd.Args())
}
