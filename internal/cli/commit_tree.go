package cli

import (
	"fmt"
	"mygit/internal/object"
	"os"

	flag "github.com/spf13/pflag"
)

func CommitTreeCommand() {
	// mygit commit-tree <tree-hash> -m "message"
	// mygit commit-tree <tree-hash> -p <parent-commit> -m "message"
	cmd := flag.NewFlagSet("commit-tree", flag.ExitOnError)

	parentHashs := cmd.StringSliceP("parent", "p", []string{}, "List of parent commit hashes")

	message := cmd.StringP("message", "m", "", "the commit message")

	cmd.Parse(os.Args[2:])

	tree_hash := cmd.Arg(0)

	err := object.CommitObject(tree_hash, *message, *parentHashs)

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
