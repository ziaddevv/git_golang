package cli

import (
	"fmt"
	"mygit/internal/object"
	"mygit/internal/repo"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

func UpdateRefCommand() error {

	/*
		create or set a reference
		mygit update-ref refs/heads/main <commit-hash>


		delete a ref
		git update-ref -d refs/heads/feature

		update only if it contains a spefic commit (safe)
		git update-ref refs/heads/main <new-hash> <old-hash>

		if head is given resolve the filename inside it
		git update-ref HEAD <commit-hash>

		and HEAD contains ref: refs/heads/main

	*/
	cmd := flag.NewFlagSet("update-ref", flag.ExitOnError)
	delete := cmd.BoolP("delete", "d", false, "delete a specific reference")

	cmd.Parse(os.Args[2:])

	if cmd.NArg() < 1 {
		return usageError("usage: mygit update-ref [-d] <ref> [<commit-hash>]")
	}

	refPath := cmd.Arg(0)

	if !isValidRefPath(refPath) {
		return fmt.Errorf("invalid ref path: %s", refPath)
	}

	if *delete {
		return repo.RemoveRef(refPath)
	}

	if cmd.NArg() < 2 {
		return usageError("usage: mygit update-ref <ref> <commit-hash>")
	}

	commitHash := cmd.Arg(1)

	if !object.IsValidObjectID(commitHash) {
		return fmt.Errorf("not a valid object name: %s", commitHash)
	}

	return repo.AddRef(refPath, commitHash)
}

func isValidRefPath(ref string) bool {
	if ref == "HEAD" {
		return true
	}
	if strings.HasPrefix(ref, "refs/") {
		return true
	}
	return false
}
