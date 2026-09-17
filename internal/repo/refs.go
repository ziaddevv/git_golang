package repo

import (
	"fmt"
	"mygit/internal/object"
	"mygit/internal/utils"
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
	finalRef, err := ResolveRef(ref)
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

func ResolveRef(ref string) (string, error) {
	currentRef := ref

	// Prevent infinite symbolic-ref loops.
	const maxDepth = 10

	for i := 0; i < maxDepth; i++ {
		refPath := filepath.Join(utils.RepoDir(), currentRef)

		// The final ref doesn't have to exist yet
		//
		// Example:
		//
		// HEAD -> refs/heads/main
		//
		// If main doesn't exist yet, we still return
		// "refs/heads/main" so AddRef can create it
		if !utils.Exists(refPath) {
			return currentRef, nil
		}

		data, err := utils.ReadFile(refPath)
		if err != nil {
			return "", fmt.Errorf(
				"failed to read ref %s: %w",
				currentRef,
				err,
			)
		}

		content := strings.TrimSpace(string(data))

		// Symbolic ref:
		//
		// ref: refs/heads/main
		if strings.HasPrefix(content, "ref: ") {
			currentRef = strings.TrimPrefix(content, "ref: ")

			if currentRef == "" {
				return "", fmt.Errorf(
					"invalid symbolic ref %s",
					ref,
				)
			}

			continue
		}

		// We reached a normal ref containing a hash.
		return currentRef, nil
	}

	return "", fmt.Errorf(
		"too many symbolic ref levels while resolving %s",
		ref,
	)
}
