package object

import (
	"bytes"
)

func CommitObject(treeHash string, parentHashs []string) {

	// construct content
	/*
		tree <tree-hash>
		parent <parent-commit-hash>     ← optional, can appear multiple times
		author <name> <email> <timestamp>
		committer <name> <email> <timestamp>

		<commit message>
	*/

	//validate the tree hash exists?

	// validate parents commits exists

	// aappends to
	var buf bytes.Buffer

	buf.WriteString("tree ")
	buf.WriteString(treeHash)
	buf.WriteByte('\n')

	for _, parent := range parentHashs {
		buf.WriteString("parent ")
		buf.WriteString(parent)
		buf.WriteByte('\n')
	}

}
