package cli

import (
	"os"

	"mygit/internal/porcelain"

	flag "github.com/spf13/pflag"
)

func BranchCommand() error {
	cmd := flag.NewFlagSet("branch", flag.ExitOnError)

	deleteBranch := cmd.BoolP("delete", "d", false, "delete branch")
	forceDelete := cmd.BoolP("delete-force", "D", false, "force delete branch")

	moveBranch := cmd.BoolP("move", "m", false, "rename branch")
	forceMove := cmd.BoolP("move-force", "M", false, "force rename branch")

	verbose := cmd.BoolP("verbose", "v", false, "show last commit")
	all := cmd.BoolP("all", "a", false, "list all branches")

	cmd.Parse(os.Args[2:])

	args := cmd.Args()

	_ = deleteBranch
	_ = forceDelete
	_ = moveBranch
	_ = forceMove
	_ = verbose
	_ = all

	switch {
	case *deleteBranch:
		return porcelain.DeleteBranch(args)

	case *forceDelete:
		return porcelain.ForceDeleteBranch(args)

	case *moveBranch:
		return porcelain.RenameBranch(args)

	case *verbose:
		return porcelain.ListBranchesVerbose()

	case *all:
		// return porcelain.ListAllBranches()

	default:

		if len(args) == 0 {
			return porcelain.ListBranches()
		}

		return porcelain.CreateBranch(args[0])
	}

	return nil
}
