package object

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"mygit/internal/index"
	"mygit/internal/repo"
	"mygit/internal/utils"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type TreeEntry struct {
	Mode string
	Path string
	Hash [20]byte
}

func newTreeEntry(mode, path string, sha [20]byte) *TreeEntry {
	return &TreeEntry{
		Mode: mode,
		Path: path,
		Hash: sha,
	}
}

type Tree struct {
	Entries []*TreeEntry
}

func (t *Tree) Add(entry *TreeEntry) {
	t.Entries = append(t.Entries, entry)
}

func BuildTree(idx *index.Index) *Tree {
	tree := &Tree{
		Entries: make([]*TreeEntry, 0, len(idx.Entries)),
	}

	for _, entry := range idx.Entries {
		mode := fmt.Sprintf("%o", entry.Mode)

		tree.Add(newTreeEntry(
			mode,
			entry.Path,
			entry.SHA,
		))
	}

	sort.Slice(tree.Entries, func(i, j int) bool {
		return tree.Entries[i].Path < tree.Entries[j].Path
	})

	return tree
}

// ----------------
type TrieNode struct {
	Children map[string]*TrieNode
	Entry    *TreeEntry
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		Children: make(map[string]*TrieNode),
		Entry:    nil,
	}
}

type Trie struct {
	Root *TrieNode
}

func (t *Trie) Insert(path []string, entry *TreeEntry) {
	cur := t.Root
	for _, str := range path {
		if cur.Children[str] == nil {
			cur.Children[str] = NewTrieNode()
		}
		cur = cur.Children[str]
	}
	// fmt.Println(cur.Children)
	cur.Entry = entry
}

func (t *Trie) ParseTreeObject(cur *TrieNode, fullPath string) (*TreeEntry, error) {
	if cur.IsFile == true {
		// create blob and return
		data, err := utils.ReadFile(fullPath)
		if err != nil {
			return nil, err
		}

		hash, err := repo.WriteObject("blob", data)

		var sha [20]byte
		raw, err := hex.DecodeString(hash)
		if err != nil {
			return nil, err
		}
		if len(raw) != 20 {
			return nil, fmt.Errorf("invalid sha1 length")
		}
		copy(sha[:], raw)

		mode := "100644"
		path := filepath.Base(fullPath)
		treeEntry := newTreeEntry(mode, path, sha)
		return treeEntry, nil

	}
	fmt.Println(cur.Children)
	// recurse on all children and return the hash
	var entries []*TreeEntry
	for key, val := range cur.Children {
		var nextPath string
		if fullPath == "" {
			nextPath = key
		} else {
			nextPath = fullPath + "/" + key
		}
		ObjEntry, err := t.ParseTreeObject(val, nextPath)
		if err != nil {
			return nil, err
		}
		entries = append(entries, ObjEntry)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	content := TreeContent(entries)
	hash, err := repo.WriteObject("tree", content)

	if err != nil {
		return nil, err
	}
	var sha [20]byte
	raw, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	if len(raw) != 20 {
		return nil, fmt.Errorf("invalid sha1 length")
	}
	copy(sha[:], raw)

	mode := "40000"
	path := filepath.Base(fullPath)
	treeEntry := newTreeEntry(mode, path, sha)
	return treeEntry, nil
	// take all hashes and make a tree object in /objects

}

func TreeContent(entries []*TreeEntry) []byte {
	var contentBuffer bytes.Buffer

	for _, entry := range entries {
		contentBuffer.WriteString(entry.Mode)
		contentBuffer.WriteByte(' ')

		contentBuffer.WriteString(entry.Path)
		contentBuffer.WriteByte('\x00')

		contentBuffer.Write(entry.Hash[:])
	}

	return contentBuffer.Bytes()
}

/*
strings.Split("/a//b", "/")
*/
func BuildTreeObject() {
	idx := repo.ReadIndex()
	tree := BuildTree(idx)

	// Initialize the Root node here
	trie := &Trie{
		Root: NewTrieNode(),
	}

	for _, entry := range tree.Entries {
		fmt.Printf("%s %s %s\n", entry.Mode, "a", entry.Path)
		trie.Insert(strings.Split(entry.Path, "/"), entry)
	}
	// trie.ParseTreeObject(trie.Root, "")
	rootEntry, err := trie.ParseTreeObject(trie.Root, "")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("%x\n", rootEntry.Hash)
}
