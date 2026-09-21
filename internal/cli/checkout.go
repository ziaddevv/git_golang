package cli

import (
	"os"

	"mygit/internal/porcelain"

	flag "github.com/spf13/pflag"
)

func CheckoutCommand() error {
	cmd := flag.NewFlagSet("checkout", flag.ExitOnError)

	createBranch := cmd.BoolP("branch", "b", false, "create and checkout a new branch")

	cmd.Parse(os.Args[2:])

	args := cmd.Args()

	if *createBranch {
		return porcelain.CreateAndCheckoutBranch(args)
	}

	return porcelain.Checkout(args)
}
