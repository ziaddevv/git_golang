package cli

import (
	"errors"
	"fmt"
	"os"
)

// UsageError indicates wrong flags or missing arguments
type UsageError struct {
	Message string
}

func (e *UsageError) Error() string {
	return e.Message
}

func usageError(msg string) error {
	return &UsageError{Message: msg}
}

func Run() {

	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: mygit <command>")
		return
	}

	var err error

	switch os.Args[1] {
	case "init":
		err = InitCommand()
	case "hash-object":
		err = HashObjectCommand()
	case "cat-file":
		err = CatFileCommand()
	case "update-index":
		err = UpdateIndexCommand()
	case "write-tree":
		err = WriteTreeCommand()
	case "commit-tree":
		err = CommitTreeCommand()
	case "update-ref":
		err = UpdateRefCommand()
	case "commit":
		err = CommitCommand()
	case "log":
		err = LogCommand()
	case "status":
		err = StatusCommand()
	case "branch":
		err = BranchCommand()
	case "checkout":
		err = CheckoutCommand()
	default:
		fmt.Fprintf(os.Stderr, "mygit: '%s' is not a mygit command\n", os.Args[1])
		return
	}

	if err != nil {
		var usageErr *UsageError
		if errors.As(err, &usageErr) {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		fmt.Fprintln(os.Stderr, "fatal:", err)
	}
}
