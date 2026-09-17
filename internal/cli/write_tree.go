package cli

import (
	"fmt"
	"mygit/internal/tree"
)

/*
 tree object

*/

func WriteTreeCommand() error {
	// get the entries from index file
	//create a tree object
	_, err := tree.BuildTreeObject()
	if err != nil {
		return fmt.Errorf("write-tree: %w", err)
	}
	return nil
}
