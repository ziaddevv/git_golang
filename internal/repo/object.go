package repo

import (
	"errors"
	"fmt"
	"mygit/internal/utils"
)

//   converting from []byte to string ---> []byte(mystring)
//   converting frokm string to []bytes ---> string(bybytes)

// [type] + " " + [size] + \0 + [content]          , sh1 hash is the object id

func WriteObject(objType string, content []byte) (string, error) {

	//TODO but this in the function calling this
	if !utils.Exists(utils.RepoDir()) {
		return "", errors.New("not a repository")
	}

	data := fmt.Sprintf("%s %d\x00%s", objType, len(content), string(content))
	rawBytes := []byte(data)
	fmt.Println("rawbytes", rawBytes)
	// data := append([]byte(objType))

	hash := utils.Hash(rawBytes)
	path := utils.ObjectPath(hash)

	if utils.Exists(path) {
		return hash, nil
	}

	err := utils.WWriteFileSafely(path, rawBytes)

	if err != nil {
		return "", err
	}

	return hash, nil
}

func ReadObject(hash string) ([]byte, error) {
	path := utils.ObjectPath(hash)
	return utils.ReadFile(path)
}
