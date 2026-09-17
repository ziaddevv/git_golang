package porcelain

import (
	"fmt"
	"mygit/internal/object"
	"mygit/internal/repo"
	"mygit/internal/tree"
)

func FullCommitObject(message string) error {

	/*
			1. Read the index
		       ↓
			2. Create a tree
				↓
			3. Find the current HEAD commit
				↓
			4. Use HEAD as the parent (if exists)
				↓
			5. Create a commit object
				↓
			6. Update HEAD/branch to the new commit
	*/

	treeHash, err := tree.BuildTreeObject()
	if err != nil {
		return fmt.Errorf("failed to write tree: %w", err)
	}

	_, parentCommit, err := repo.ResolveRef("HEAD")
	if err != nil {
		return fmt.Errorf("failed to resolve HEAD: %w", err)
	}

	var parents []string
	if parentCommit != "" {
		parents = []string{parentCommit}
	}

	commitHash, err := object.CommitObject(treeHash, message, parents)
	if err != nil {
		return err
	}

	if err := repo.AddRef("HEAD", commitHash); err != nil {
		return fmt.Errorf("failed to update HEAD: %w", err)
	}

	fmt.Println(commitHash)

	return nil
}
