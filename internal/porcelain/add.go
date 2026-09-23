package porcelain

import (
	"mygit/internal/index"
	"os"
	"path/filepath"
)

func Add(paths []string) error {
	for _, path := range paths {
		err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && info.Name() == ".mygit" {
				return filepath.SkipDir
			}

			if info.IsDir() {
				return nil
			}

			return index.UpdateIndexAdd(p)
		})

		if err != nil {
			return err
		}
	}
	return nil
}
