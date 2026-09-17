package cli

import (
	"fmt"
	"mygit/internal/object"
	"mygit/internal/repo"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

func UpdateRefCommand() {

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
		fmt.Fprintln(os.Stderr, "usage: mygit update-ref [-d] <ref> [<commit-hash>]")
		os.Exit(1)
	}

	refPath := cmd.Arg(0)

	if !isValidRefPath(refPath) {
		fmt.Fprintf(os.Stderr, "fatal: invalid ref path: %s\n", refPath)
		os.Exit(1)
	}

	if *delete {
		err := repo.RemoveRef(refPath)

		if err != nil {
			fmt.Fprintf(os.Stderr, "fatal: %s\n", err)
			os.Exit(1)
		}
		return
	}

	if cmd.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "usage: mygit update-ref <ref> <commit-hash>")
		os.Exit(1)
	}

	commitHash := cmd.Arg(1)

	if !object.IsValidObjectID(commitHash) {
		fmt.Fprintf(os.Stderr, "fatal: not a valid object name: %s\n", commitHash)
		os.Exit(1)
	}

	err := repo.AddRef(refPath, commitHash)

	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %s\n", err)
		return
	}
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
