package object

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

type Commit struct {
	Hash          string
	TreeHash      string
	ParentCommits []string
	Author        string
	Commiter      string
	Date          int64
	Message       string
}

func NewCommit(hash, treeHash, author, commiter string, date int64, message string, parents []string) *Commit {
	return &Commit{
		Hash:          hash,
		TreeHash:      treeHash,
		Author:        author,
		Commiter:      commiter,
		Date:          date,
		Message:       message,
		ParentCommits: parents,
	}
}

func GetCommitbyHash(commitHash string) (*Commit, error) {
	fileContent, err := content(commitHash)
	if err != nil {
		return nil, err
	}

	return ParseCommit(commitHash, fileContent)
}

// Added hash parameter
func ParseCommit(hash string, content []byte) (*Commit, error) {
	lines := strings.Split(string(content), "\n")

	commit := &Commit{
		Hash:          hash,
		ParentCommits: []string{},
	}

	inMessage := false
	var messageBuilder strings.Builder

	for _, line := range lines {
		if inMessage {
			messageBuilder.WriteString(line + "\n")
			continue
		}

		if line == "" {
			inMessage = true
			continue
		}

		if strings.HasPrefix(line, "tree ") {
			commit.TreeHash = strings.TrimPrefix(line, "tree ")

		} else if strings.HasPrefix(line, "parent ") {
			commit.ParentCommits = append(commit.ParentCommits, strings.TrimPrefix(line, "parent "))

		} else if strings.HasPrefix(line, "author ") {
			authorLine := strings.TrimPrefix(line, "author ")

			fields := strings.Fields(authorLine)
			if len(fields) >= 2 {
				timestampStr := fields[len(fields)-2]
				if parsedTime, err := strconv.ParseInt(timestampStr, 10, 64); err == nil {
					commit.Date = parsedTime
				}
				commit.Author = strings.Join(fields[:len(fields)-2], " ")
			}

		} else if strings.HasPrefix(line, "committer ") {
			committerLine := strings.TrimPrefix(line, "committer ")

			fields := strings.Fields(committerLine)
			if len(fields) >= 2 {
				commit.Commiter = strings.Join(fields[:len(fields)-2], " ")
			}
		}
	}

	commit.Message = strings.TrimSpace(messageBuilder.String())

	return commit, nil
}

func CommitObject(treeHash string, commitMessage string, parentHashes []string) (string, error) {

	// construct content
	/*
		tree <tree-hash>
		parent <parent-commit-hash>     ← optional, can appear multiple times
		author <name> <email> <timestamp>
		committer <name> <email> <timestamp>

		<commit message>
	*/

	// validate the tree hash exists and is actually a tree
	treeData, err := ReadObject(treeHash)
	if err != nil {
		return "", fmt.Errorf("invalid tree object %s: %w", treeHash, err)
	}

	objType, _, err := ParseObject(treeData)
	if err != nil {
		return "", fmt.Errorf("corrupted tree object %s: %w", treeHash, err)
	}

	if objType != "tree" {
		return "", fmt.Errorf("object %s is a %s, not a tree", treeHash, objType)
	}

	// validate parent commits exist and are actually commits
	for _, parentHash := range parentHashes {
		parentData, err := ReadObject(parentHash)
		if err != nil {
			return "", fmt.Errorf("invalid parent commit %s: %w", parentHash, err)
		}

		parentType, _, err := ParseObject(parentData)

		if err != nil {
			return "", fmt.Errorf("corrupted parent object %s: %w", parentHash, err)
		}

		if parentType != "commit" {
			return "", fmt.Errorf("object %s is a %s, not a commit", parentHash, parentType)
		}
	}

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
		return "", fmt.Errorf("failed to load .mygit/config: %w", err)
	}

	name := cfg.Section("user").Key("name").String()
	email := cfg.Section("user").Key("email").String()

	if name == "" || email == "" {
		return "", fmt.Errorf("user.name and user.email must be set in .mygit/config")
	}

	now := time.Now()
	timestamp := strconv.FormatInt(now.Unix(), 10)
	timezone := now.Format("-0700")

	fmt.Println(name, email)

	writeIdentity(&buf, "author", name, email, timestamp, timezone)
	writeIdentity(&buf, "committer", name, email, timestamp, timezone)

	buf.WriteByte('\n')
	buf.WriteString(commitMessage)
	buf.WriteByte('\n')

	fmt.Println(buf.String())
	hash, err := WriteObject("commit", buf.Bytes())

	if err != nil {
		return "", fmt.Errorf("failed to write commit object: %w", err)
	}
	fmt.Println(hash)

	return hash, nil
}

func writeIdentity(buf *bytes.Buffer, role, name, email, timestamp, timezone string) {
	buf.WriteString(role)
	buf.WriteString(" ")
	buf.WriteString(name)
	buf.WriteString(" <")
	buf.WriteString(email)
	buf.WriteString("> ")
	buf.WriteString(timestamp)
	buf.WriteString(" ")
	buf.WriteString(timezone)
	buf.WriteByte('\n')
}
