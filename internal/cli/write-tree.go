package cli

import (
	"mygit/internal/object"
)

/*
 tree object






*/

func WriteTreeCommand() {
	// get the entries from index file
	//create a tree object
	object.BuildTreeObject()
}
