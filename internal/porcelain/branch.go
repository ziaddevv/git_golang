package porcelain

import (
	"fmt"
	"mygit/internal/repo"
	"mygit/internal/utils"
	"path/filepath"
	"slices"
)

func ListBranches() error {
	branches, err := repo.ListBranches()
	if err != nil {
		return err
	}
	currentBranch, err := repo.GetCurrentBranch(utils.HeadPath())
	if err != nil {
		return err
	}
	for _, branchName := range branches {

		if currentBranch != branchName {
			fmt.Println(" ", branchName)
		} else {
			fmt.Println("*", branchName)
		}
	}

	return nil
}

func CreateBranch(branchName string) error {

	branches, err := repo.ListBranches()

	if err != nil {
		return err
	}

	if exists := slices.Contains(branches, branchName); exists {
		return fmt.Errorf("fatal: a branch named %s already exists", branchName)
	}
	_, commitHash, err := repo.ResolveRef("HEAD")

	if err != nil {
		return err
	}
	refPath := filepath.Join("refs/heads", branchName)
	err = repo.AddRef(refPath, commitHash)
	if err != nil {
		return err
	}

	return nil
}

func DeleteBranch(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: mygit branch -d <branch>")
	}

	branchName := args[0]

	branches, err := repo.ListBranches()
	if err != nil {
		return err
	}

	if !slices.Contains(branches, branchName) {
		return fmt.Errorf("fatal: branch '%s' not found", branchName)
	}

	currentBranch, err := repo.GetCurrentBranch(utils.HeadPath())
	if err != nil {
		return err
	}

	if currentBranch == branchName {
		return fmt.Errorf(
			"fatal: cannot delete branch '%s' checked out",
			branchName,
		)
	}

	refPath := filepath.Join("refs/heads", branchName)

	return repo.RemoveRef(refPath)
}
