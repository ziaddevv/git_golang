package object

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

type TreeParsedEntry struct {
	Mode string
	Type string
	Hash string
	Name string
}

func ParseTreeContent(content []byte) string {
	r := bufio.NewReader(bytes.NewReader(content))
	var out strings.Builder

	for {
		entry, err := ReadTreeEntry(r)
		if err != nil {
			break
		}
		// format: 100644 blob b6fc4c620b67d95f    test.txt
		fmt.Fprintf(&out, "%06s %s %s\t%s\n", entry.Mode, entry.Type, entry.Hash, entry.Name)
	}

	return out.String()
}

func ReadTreeEntry(r *bufio.Reader) (TreeParsedEntry, error) {
	var e TreeParsedEntry

	modeBytes, err := r.ReadBytes(' ')
	if err != nil {
		return e, err
	}
	e.Mode = string(modeBytes[:len(modeBytes)-1])

	nameBytes, err := r.ReadBytes('\x00')
	if err != nil {
		return e, err
	}
	e.Name = string(nameBytes[:len(nameBytes)-1])

	var sha [20]byte
	if _, err := io.ReadFull(r, sha[:]); err != nil {
		return e, err
	}
	e.Hash = fmt.Sprintf("%x", sha)

	if e.Mode == "40000" {
		e.Type = "tree"
	} else {
		e.Type = "blob"
	}

	return e, nil
}
func ReadTreeEntries(treeHash string) ([]TreeParsedEntry, error) {
	data, err := ReadObject(treeHash)
	if err != nil {
		return nil, err
	}

	objType, content, err := ParseObject(data)
	if err != nil {
		return nil, err
	}
	if objType != "tree" {
		return nil, fmt.Errorf("object %s is a %s, not a tree", treeHash, objType)
	}

	r := bufio.NewReader(bytes.NewReader(content))
	var entries []TreeParsedEntry

	for {
		entry, err := ReadTreeEntry(r)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			break
		}
		entries = append(entries, entry)
	}

	return entries, nil
}
