package repo

import (
	"errors"
	"fmt"
	"mygit/internal/utils"
)

//   converting from []byte to string ---> []byte(mystring)
//   converting frokm string to []bytes ---> string(bybytes)

// [type] + " " + [size] + \0 + [content]          , sh1 hash is the object id

// TODO compress before storing

func HashObject(objType string, content []byte) ([]byte, string) {
	data := fmt.Sprintf("%s %d\x00%s", objType, len(content), string(content))
	rawBytes := []byte(data)
	hash := utils.Hash(rawBytes)
	return rawBytes, hash
}
func WriteObject(objType string, content []byte) (string, error) {

	if !utils.Exists(utils.RepoDir()) {
		return "", errors.New("not a repository")
	}

	rawBytes, hash := HashObject(objType, content)
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

// func BlobObject(,write bool)([]byte , err){
// 	var data []byte
// 	var err error

// 	return  data ,err
// }
// func TreeObject()([]byte , err){
// 	var data []byte
// 	var err error

// 	return  data ,err
// }
// func CommitObject()([]byte , err){
// 	var data []byte
// 	var err error

// 	return  data ,err
// }

func ReadObject(hash string) ([]byte, error) {
	path := utils.ObjectPath(hash)
	return utils.ReadFile(path)
}
