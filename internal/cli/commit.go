package cli

import (
	"fmt"
	"mygit/internal/porcelain"

	flag "github.com/spf13/pflag"
)

func CommitCommand() error {

	cmd := flag.NewFlagSet("commit-tree", flag.ExitOnError)

	message := cmd.StringP("message", "m", "", "the commit message")

	fmt.Println(message)

	err := porcelain.FullCommitObject(*message)

	return err
}
