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

	parentHashes := cmd.StringSliceP("parent", "p", []string{}, "List of parent commit hashes")

	message := cmd.StringP("message", "m", "", "the commit message")

	cmd.Parse(os.Args[2:])

	if cmd.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: mygit commit-tree <tree-hash> -m <message> [-p <parent>]...")
		os.Exit(1)
	}

	treeHash := cmd.Arg(0)

	if !object.IsValidObjectID(treeHash) {
		fmt.Fprintf(os.Stderr, "fatal: not a valid object name: %s\n", treeHash)
		os.Exit(1)
	}

	if *message == "" {
		fmt.Fprintln(os.Stderr, "fatal: commit message is required (-m)")
		os.Exit(1)
	}

	for _, p := range *parentHashes {
		if !object.IsValidObjectID(p) {
			fmt.Fprintf(os.Stderr, "fatal: not a valid parent hash: %s\n", p)
			os.Exit(1)
		}
	}

	err := object.CommitObject(treeHash, *message, *parentHashes)

	if err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

