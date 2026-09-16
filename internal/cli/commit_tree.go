package cli

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

func CommitTreeCommand() {
	// mygit commit-tree <tree-hash> -m "message"
	// mygit commit-tree <tree-hash> -p <parent-commit> -m "message"
	cmd := flag.NewFlagSet("commit-tree", flag.ExitOnError)

	parent := cmd.StringP("parent", "p", "", "the parent commit")

	message := cmd.StringP("message", "m", "", "the commit message")

	cmd.Parse(os.Args[2:])

	fmt.Println(*parent, *message)

	fmt.Println(os.Args)

	tree_hash := cmd.Arg(0)

	fmt.Println(tree_hash)
}
