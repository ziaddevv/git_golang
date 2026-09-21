package repo

import (
	"fmt"
	"mygit/internal/object"
	"mygit/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

func RemoveRef(ref string) error {
	refPath := filepath.Join(utils.RepoDir(), ref)

	if !utils.Exists(refPath) {
		return fmt.Errorf("ref %s does not exist", ref)
	}

	return utils.RemoveFile(refPath)
}

func AddRef(ref string, commitHash string) error {

	//  validate commit hash

	data, err := object.ReadObject(commitHash)
	if err != nil {
		return fmt.Errorf("invalid object %s: %w", commitHash, err)
	}

	objType, _, err := object.ParseObject(data)
	if err != nil {
		return fmt.Errorf("corrupted object %s: %w", commitHash, err)
	}

	if objType != "commit" {
		return fmt.Errorf(
			"object %s is a %s, not a commit",
			commitHash,
			objType,
		)
	}
	/*
		resolve symbolic refs.
		HEAD:
		ref: refs/heads/main
	*/
	finalRef, _, err := ResolveRef(ref)
	if err != nil {
		return err
	}

	refPath := filepath.Join(utils.RepoDir(), finalRef)

	parentDir := filepath.Dir(refPath)

	if err := utils.EnsureDir(parentDir); err != nil {
		return fmt.Errorf(
			"failed to create ref directory: %w",
			err,
		)
	}

	content := []byte(commitHash + "\n")

	if err := utils.WriteFileSafely(refPath, content); err != nil {
		return fmt.Errorf(
			"failed to update ref %s: %w",
			finalRef,
			err,
		)
	}

	return nil
}

func ResolveRef(ref string) (string, string, error) {
	currentRef := ref

	const maxDepth = 10

	for i := 0; i < maxDepth; i++ {
		refPath := filepath.Join(utils.RepoDir(), currentRef)

		if !utils.Exists(refPath) {
			return currentRef, "", nil
		}

		data, err := utils.ReadFile(refPath)
		if err != nil {
			return "", "", fmt.Errorf(
				"failed to read ref %s: %w",
				currentRef,
				err,
			)
		}

		content := strings.TrimSpace(string(data))

		// Symbolic ref:
		// ref: refs/heads/main
		if strings.HasPrefix(content, "ref: ") {
			currentRef = strings.TrimPrefix(content, "ref: ")

			if currentRef == "" {
				return "", "", fmt.Errorf(
					"invalid symbolic ref %s",
					ref,
				)
			}

			continue
		}

		// Normal ref containing an object hash.
		return currentRef, content, nil
	}

	return "", "", fmt.Errorf(
		"too many symbolic ref levels while resolving %s",
		ref,
	)
}

func ListBranches() ([]string, error) {
	entries, err := os.ReadDir(utils.RefsDir())
	if err != nil {
		return nil, err
	}

	var branches []string
	for _, e := range entries {
		if !e.IsDir() {
			branches = append(branches, e.Name())
		}
	}

	return branches, nil
}

func GetCurrentBranch(headPath string) (string, error) {
	content, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}

	head := strings.TrimSpace(string(content))

	if strings.HasPrefix(head, "ref: refs/heads/") {
		return strings.TrimPrefix(head, "ref: refs/heads/"), nil
	}

	return head, nil
}
