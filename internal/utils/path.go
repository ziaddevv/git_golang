package utils

import (
	"errors"
	"os"
	"path/filepath"
)

// func EnsureDir(path string) bool{
// 	inf , err := os.Stat(path)

// 	return nil
// }
// TODO : function to check for folders existance
// TODO : function to get each path to our folders eg objects , branches etc

func RepoDir() string {
	return ".mygit"
}
func WorkDir() (string, error) {
	return os.Getwd()
}

func ObjectsDir() string {
	return filepath.Join(RepoDir(), "objects")
}

func ObjectPath(hash string) (string, error) {
	if len(hash) < 2 {
		return "", errors.New("not a repository")
	}
	return filepath.Join(ObjectsDir(), hash[0:2], hash[2:]), nil
}

func IndexPath() string {
	return "INDEX"
}
