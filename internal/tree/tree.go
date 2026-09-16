package tree

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"mygit/internal/index"
	"mygit/internal/object"
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
	// index here is just a list of entries which holds info about each tracked file
	// tree is a list of entries which constructs the tree object , each entry holds--> mode path hash
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
	IsFile   bool
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		Children: make(map[string]*TrieNode),
		Entry:    nil,
		IsFile:   false,
	}
}

type Trie struct {
	Root *TrieNode
}

func (t *Trie) Insert(path []string, entry *TreeEntry) error {
	cur := t.Root
	for i, str := range path {
		if str == "" || str == "." || str == ".." {
			return fmt.Errorf("invalid tree path: %s", entry.Path)
		}
		if cur.IsFile {
			return fmt.Errorf("file/directory conflict: %s", entry.Path)
		}
		if cur.Children[str] == nil {
			cur.Children[str] = NewTrieNode()
		}
		cur = cur.Children[str]
		if i == len(path)-1 && len(cur.Children) > 0 {
			return fmt.Errorf("file/directory conflict: %s", entry.Path)
		}
	}
	// fmt.Println(cur.Children)
	cur.Entry = &TreeEntry{
		Mode: entry.Mode,
		Path: filepath.Base(entry.Path),
		Hash: entry.Hash,
	}
	cur.IsFile = true
	return nil
}

func (t *Trie) ParseTreeObject(cur *TrieNode, fullPath string) (*TreeEntry, error) {
	if cur.IsFile == true {
		return cur.Entry, nil

	}
	// fmt.Println(cur.Children)
	// for key, value := range cur.Children {

	// }
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
		return treeSortName(entries[i]) < treeSortName(entries[j])
	})

	content := TreeContent(entries)
	hash, err := object.WriteObject("tree", content)

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

func treeSortName(entry *TreeEntry) string {
	if entry.Mode == "40000" {
		return entry.Path + "/"
	}
	return entry.Path
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
	idx, err := index.ReadIndex()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	tree := BuildTree(idx)

	// Initialize the Root node here
	trie := &Trie{
		Root: NewTrieNode(),
	}

	for _, entry := range tree.Entries {
		fmt.Printf("%s %s %s\n", entry.Mode, "a", entry.Path)
		if _, err := object.ReadObject(fmt.Sprintf("%x", entry.Hash)); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		if err := trie.Insert(strings.Split(entry.Path, "/"), entry); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	}
	// trie.ParseTreeObject(trie.Root, "")
	rootEntry, err := trie.ParseTreeObject(trie.Root, "")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("%x\n", rootEntry.Hash)

	fmt.Println("-------------------")
	PrintTrie(trie.Root, "")
}

// //////////////////////////////////////////////////////////////
// ------------
func PrintTrie(root *TrieNode, prefix string) {
	printTrieHelper(root, prefix, "")
}

func printTrieHelper(node *TrieNode, prefix, name string) {
	if node == nil {
		return
	}

	// Print current node
	marker := ""
	if node.IsFile {
		marker = " [FILE]"
	}

	if name != "" {
		fmt.Printf("%s%s%s\n", prefix, name, marker)
	} else {
		fmt.Println("root")
	}

	// Print Entry if exists
	if node.Entry != nil {
		fmt.Printf("%s   └─ Entry: %+v\n", prefix, node.Entry)
	}

	// Print children
	for key, child := range node.Children {
		newPrefix := prefix + "│   "
		if key == getLastKey(node.Children) {
			newPrefix = prefix + "    "
		}
		printTrieHelper(child, newPrefix, key)
	}
}

// Helper to make tree look nicer (find last key)
func getLastKey(m map[string]*TrieNode) string {
	for k := range m {
		return k // returns the last key in map iteration (not perfect but good enough for display)
	}
	return ""
}
