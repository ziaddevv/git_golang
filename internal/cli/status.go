package cli

import (
	"mygit/internal/porcelain"
	"os"

	flag "github.com/spf13/pflag"
)

func StatusCommand() {
	cmd := flag.NewFlagSet("status", flag.ExitOnError)

	cmd.Parse(os.Args[2:])

	porcelain.Status()
}
