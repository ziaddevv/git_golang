package cli

import (
	"mygit/internal/porcelain"
	"os"

	flag "github.com/spf13/pflag"
)

func CommitCommand() error {

	cmd := flag.NewFlagSet("commit", flag.ExitOnError)

	message := cmd.StringP("message", "m", "", "the commit message")

	cmd.Parse(os.Args[2:])

	if *message == "" {
		return usageError("usage: mygit commit -m <message>")
	}

	return porcelain.FullCommitObject(*message)
}
