package porcelain

import (
	"fmt"
	"mygit/internal/object"
	"mygit/internal/repo"
	"mygit/internal/utils"
	"os"
	"path/filepath"
	"slices"
	"strings"
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

func ForceDeleteBranch(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: mygit branch -D <branch>")
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

func RenameBranch(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: mygit branch -m <old> <new>")
	}

	oldName := args[0]
	newName := args[1]

	branches, err := repo.ListBranches()
	if err != nil {
		return err
	}

	if !slices.Contains(branches, oldName) {
		return fmt.Errorf("fatal: branch '%s' not found", oldName)
	}

	if slices.Contains(branches, newName) {
		return fmt.Errorf("fatal: a branch named '%s' already exists", newName)
	}

	oldRef := filepath.Join(utils.RepoDir(), "refs/heads", oldName)
	newRef := filepath.Join(utils.RepoDir(), "refs/heads", newName)

	if err := os.Rename(oldRef, newRef); err != nil {
		return err
	}

	currentBranch, err := repo.GetCurrentBranch(utils.HeadPath())
	if err != nil {
		return err
	}

	if currentBranch == oldName {
		headContent := []byte("ref: refs/heads/" + newName + "\n")

		if err := utils.WriteFileSafely(utils.HeadPath(), headContent); err != nil {
			return err
		}
	}

	return nil
}

func ListBranchesVerbose() error {
	branches, err := repo.ListBranches()
	if err != nil {
		return err
	}

	currentBranch, err := repo.GetCurrentBranch(utils.HeadPath())
	if err != nil {
		return err
	}

	for _, branchName := range branches {
		ref := filepath.Join("refs/heads", branchName)

		_, commitHash, err := repo.ResolveRef(ref)
		if err != nil {
			return err
		}

		data, err := object.ReadObject(commitHash)
		if err != nil {
			return err
		}

		_, content, err := object.ParseObject(data)
		if err != nil {
			return err
		}

		commit := string(content)
		lines := strings.SplitN(commit, "\n\n", 2)

		subject := ""
		if len(lines) == 2 {
			message := strings.TrimSpace(lines[1])
			if message != "" {
				subject = strings.SplitN(message, "\n", 2)[0]
			}
		}

		prefix := " "
		if branchName == currentBranch {
			prefix = "*"
		}

		shortHash := commitHash[:7]

		fmt.Printf("%s %-10s %s %s\n",
			prefix,
			branchName,
			shortHash,
			subject,
		)
	}

	return nil
}
