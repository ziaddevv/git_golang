package worktree

import (
	"io/fs"
	"mygit/internal/utils"
	"os"
	"path/filepath"
)

type treeReader struct {
	baseDir string
	files   map[string][]byte
}

func (t *treeReader) processFile(path string, d fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if d.IsDir() && (d.Name() == ".git" || d.Name() == ".mygit") {
		return filepath.SkipDir
	}

	if !d.IsDir() {
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(t.baseDir, path)
		if err != nil {
			return err
		}

		relPath = filepath.ToSlash(relPath)
		t.files[relPath] = content
	}

	return nil
}

func ReadWorkTree() (map[string][]byte, error) {
	dirPath, err := utils.WorkDir()
	if err != nil {
		return nil, err
	}

	reader := &treeReader{
		baseDir: dirPath,
		files:   make(map[string][]byte),
	}

	err = filepath.WalkDir(dirPath, reader.processFile)
	if err != nil {
		return nil, err
	}

	return reader.files, nil
}
