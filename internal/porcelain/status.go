package porcelain

import (
	"fmt"
	"mygit/internal/index"
	"mygit/internal/object"
	"mygit/internal/repo"
	"mygit/internal/utils"
	"mygit/internal/worktree"
	"path/filepath"
)

type Snapshot map[string]string

type ChangeType int

const (
	Modified ChangeType = iota
	Added
	Deleted
)

func (c ChangeType) String() string {
	return [...]string{"Modified", "Added", "Deleted"}[c]
}

type Change struct {
	Path string
	Type ChangeType
}

// just for colors
const (
	Green  = "\033[32m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Reset  = "\033[0m"
)

func Status() error {
	/*
		so in status we have to compare three stages
		1-  Head ---> last commit on current branch
		2-  INDEX ---> index file that hold entries
		3-  working tree ---> scan all the working tree files and get the blobhash

		then compare all those

		1- compare HEAD <--> INDEX

			what changes since last commit
				foo.go --> hash1
				foo.go --> hash2

		2- compare INDEX <---> working tree

			what changes in my filestha ti haven't staged yet
	*/

	idx, err := index.ReadIndex()
	if err != nil {
		return err
	}

	indexMap := ParseIndex(idx)

	_, commitHash, err := repo.ResolveRef("HEAD")
	if err != nil {
		return err
	}

	headMap := make(Snapshot)

	if commitHash != "" {
		lastCommit, err := object.GetCommitbyHash(commitHash)
		if err != nil {
			return err
		}

		treeEntries, err := object.ReadTreeEntries(lastCommit.TreeHash)
		if err != nil {
			return err
		}

		err = ParseTreeObject(treeEntries, "", headMap)
		if err != nil {
			return err
		}
	}

	workDirEntries, err := worktree.ReadWorkTree()
	if err != nil {
		return err
	}

	workDirMap := ParseWorkTree(workDirEntries)

	stagedChanges := Compare(headMap, indexMap)
	unstagedChanges := Compare(indexMap, workDirMap)

	PrintStagedChanges(stagedChanges)

	fmt.Println()

	PrintUnstagedChanges(unstagedChanges)

	fmt.Println()

	PrintUntrackedFiles(unstagedChanges)

	return nil
}

func PrintStagedChanges(changes []Change) {
	fmt.Println("Changes to be committed:")

	for _, change := range changes {
		switch change.Type {
		case Added:
			fmt.Printf("  %snew file: %s%s\n",
				utils.Green, change.Path, utils.Reset)

		case Modified:
			fmt.Printf("  %smodified: %s%s\n",
				utils.Green, change.Path, utils.Reset)

		case Deleted:
			fmt.Printf("  %sdeleted: %s%s\n",
				utils.Green, change.Path, utils.Reset)
		}
	}
}

func PrintUnstagedChanges(changes []Change) {
	fmt.Println("Changes not staged for commit:")

	for _, change := range changes {
		switch change.Type {
		case Modified:
			fmt.Printf("  %smodified: %s%s\n",
				utils.Red, change.Path, utils.Reset)

		case Deleted:
			fmt.Printf("  %sdeleted: %s%s\n",
				utils.Red, change.Path, utils.Reset)
		}
	}
}

func PrintUntrackedFiles(changes []Change) {
	fmt.Println("Untracked files:")

	for _, change := range changes {
		if change.Type == Added {
			fmt.Printf("  %s%s%s\n",
				utils.Red, change.Path, utils.Reset)
		}
	}
}

func ParseIndex(idx *index.Index) Snapshot {

	result := make(Snapshot)

	for _, entry := range idx.Entries {
		hash := fmt.Sprintf("%x", entry.SHA)
		result[entry.Path] = hash
	}

	return result
}

func ParseTreeObject(treeEntries []object.TreeParsedEntry, fullPath string, result Snapshot) error {

	for _, entry := range treeEntries {

		currentPath := filepath.Join(fullPath, entry.Name)

		switch entry.Type {
		case "blob":
			result[currentPath] = entry.Hash

		case "tree":
			treeEntries, err := object.ReadTreeEntries(entry.Hash)
			if err != nil {
				return err
			}

			err = ParseTreeObject(
				treeEntries,
				currentPath,
				result,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ParseWorkTree(workDirEntries map[string][]byte) Snapshot {

	result := make(Snapshot)

	for path, content := range workDirEntries {
		_, hash := object.HashObject("blob", content)
		result[path] = hash
	}

	return result
}

func Compare(track1 map[string]string, track2 map[string]string) []Change {

	var result []Change

	paths := make(map[string]struct{})

	for path := range track1 {
		paths[path] = struct{}{}
	}

	for path := range track2 {
		paths[path] = struct{}{}
	}

	for path := range paths {

		value1, exists1 := track1[path]
		value2, exists2 := track2[path]

		switch {
		case !exists1 && exists2:
			result = append(result, Change{
				Path: path,
				Type: Added,
			})

		case exists1 && !exists2:
			result = append(result, Change{
				Path: path,
				Type: Deleted,
			})

		case exists1 && exists2 && value1 != value2:
			result = append(result, Change{
				Path: path,
				Type: Modified,
			})

		case exists1 && exists2 && value1 == value2:
			continue
		}
	}

	return result
}
