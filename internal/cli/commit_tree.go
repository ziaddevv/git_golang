package cli

import (
	"fmt"
	"mygit/internal/object"
	"os"

	flag "github.com/spf13/pflag"
)

func CommitTreeCommand() error {
	// mygit commit-tree <tree-hash> -m "message"
	// mygit commit-tree <tree-hash> -p <parent-commit> -m "message"
	cmd := flag.NewFlagSet("commit-tree", flag.ExitOnError)

	parentHashes := cmd.StringSliceP("parent", "p", []string{}, "List of parent commit hashes")

	message := cmd.StringP("message", "m", "", "the commit message")

	cmd.Parse(os.Args[2:])

	if cmd.NArg() < 1 {
		return usageError("usage: mygit commit-tree <tree-hash> -m <message> [-p <parent>]...")
	}

	treeHash := cmd.Arg(0)

	if !object.IsValidObjectID(treeHash) {
		return fmt.Errorf("not a valid object name: %s", treeHash)
	}

	if *message == "" {
		return fmt.Errorf("commit message is required (-m)")
	}

	for _, p := range *parentHashes {
		if !object.IsValidObjectID(p) {
			return fmt.Errorf("not a valid parent hash: %s", p)
		}
	}

	return object.CommitObject(treeHash, *message, *parentHashes)
}
