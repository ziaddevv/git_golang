package object

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/ini.v1"
)

func CommitObject(treeHash string, commitMessage string, parentHashes []string) error {

	// construct content
	/*
		tree <tree-hash>
		parent <parent-commit-hash>     ← optional, can appear multiple times
		author <name> <email> <timestamp>
		committer <name> <email> <timestamp>

		<commit message>
	*/

	//validate the tree hash exists?
	_, err := ReadObject(treeHash)

	if err != nil {
		return err
	}
	// validate parents commits exists
	for _, parentHash := range parentHashes {
		_, err := ReadObject(parentHash)
		if err != nil {
			return err
		}
	}
	// aappends to
	var buf bytes.Buffer

	buf.WriteString("tree ")
	buf.WriteString(treeHash)
	buf.WriteByte('\n')

	for _, parent := range parentHashes {
		buf.WriteString("parent ")
		buf.WriteString(parent)
		buf.WriteByte('\n')
	}

	// read from config file
	cfg, err := ini.Load(".mygit/config")

	if err != nil {
		return errors.New("failed to find Configuration file")
	}

	name := cfg.Section("user").Key("name").String()
	email := cfg.Section("user").Key("email").String()
	now := time.Now()
	// timestampStr := strconv.FormatInt(now.Unix(), 10)

	fmt.Println(name, email)

	buf.WriteString("author ")
	buf.WriteString(name)
	buf.WriteString(" ")
	buf.WriteString("<")
	buf.WriteString(email)
	buf.WriteString(">")
	buf.WriteString(" ")
	buf.WriteString(strconv.FormatInt(now.Unix(), 10))

	buf.WriteString(" ")

	buf.WriteString(now.Format("-0700"))
	buf.WriteByte('\n')

	buf.WriteString("committer ")
	buf.WriteString(name)
	buf.WriteString(" ")
	buf.WriteString("<")
	buf.WriteString(email)
	buf.WriteString(">")
	buf.WriteString(" ")
	buf.WriteString(strconv.FormatInt(now.Unix(), 10))

	buf.WriteString(" ")

	buf.WriteString(now.Format("-0700"))
	buf.WriteByte('\n')
	buf.WriteByte('\n')

	buf.WriteString(commitMessage)
	buf.WriteByte('\n')

	fmt.Println(buf.String())
	hash, err := WriteObject("commit", buf.Bytes())

	if err != nil {
		return err
	}
	fmt.Println(hash)

	return nil
}
