package porcelain

import (
	"fmt"
	"path/filepath"

	"mygit/internal/repo"
	"mygit/internal/utils"
)

func Checkout(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: mygit checkout <branch>")
	}

	branchName := args[0]

	branches, err := repo.ListBranches()
	if err != nil {
		return err
	}

	found := false
	for _, branch := range branches {
		if branch == branchName {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("error: branch '%s' not found", branchName)
	}

	_, commitHash, err := repo.ResolveRef(
		filepath.Join("refs/heads", branchName),
	)
	if err != nil {
		return err
	}

	if commitHash == "" {
		return fmt.Errorf("error: branch '%s' has no commit", branchName)
	}

	headContent := []byte("ref: refs/heads/" + branchName + "\n")

	return utils.WriteFileSafely(
		utils.HeadPath(),
		headContent,
	)
}

func CreateAndCheckoutBranch(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: mygit checkout -b <branch>")
	}

	branchName := args[0]

	if err := CreateBranch(branchName); err != nil {
		return err
	}

	return Checkout(args)
}
