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
			4. Use HEAD as the parent
				↓
			5. Create a commit object
				↓
			6. Update HEAD/branch to the new commit
	*/
	treeHash, err := tree.BuildTreeObject()

	if err != nil {
		return err
	}
	_, parentCommit, err := repo.ResolveRef("HEAD")

	if err != nil {
		return err
	}

	commitHash, err := object.CommitObject(treeHash, message, []string{parentCommit})

	if err != nil {
		return err
	}

	repo.AddRef("HEAD", commitHash)

	fmt.Println(commitHash)

	return nil
}
