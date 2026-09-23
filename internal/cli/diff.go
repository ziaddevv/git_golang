package cli

import (
	"mygit/internal/porcelain"
	"os"

	flag "github.com/spf13/pflag"
)

func DiffCommand() error {
	cmd := flag.NewFlagSet("diff", flag.ExitOnError)
	staged := cmd.Bool("staged", false, "diff of between staged and last commit")

	cmd.Parse(os.Args[2:])

	if *staged {
		err := porcelain.DiffStaged()
		if err != nil {
			return err
		}
	} else {

		err := porcelain.Diff()
		if err != nil {
			return err
		}
	}

	return nil
}
