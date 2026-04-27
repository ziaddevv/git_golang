package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func ObjectExists(hash string) bool {
	objPath := ObjectPath(hash)
	return Exists(objPath)
}

func CreateDir(name string, path string) error {
	fullPath := filepath.Join(path, name)

	err := os.MkdirAll(fullPath, 0755)
	if err != nil {
		return err
	}

	return nil

}
func CreateEmptyFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	return file.Close()
}

func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
func WWriteFileSafely(path string, data []byte) error { // same as write to fie but it checks the folder existing and creating it if not
	folderPath := filepath.Dir(path)

	err := os.MkdirAll(folderPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func AppendFile(path string, data []byte) error {
	file, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
