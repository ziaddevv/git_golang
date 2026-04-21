package utils

import "os"

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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
