package repo

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mygit/internal/utils"
	"os"
	"strings"
)

//   converting from []byte to string ---> []byte(mystring)
//   converting frokm string to []bytes ---> string(bybytes)

// [type] + " " + [size] + \0 + [content]          , sh1 hash is the object id

// TODO compress before storing

//	func HashObject(objType string, content []byte) ([]byte, string) {
//		data := fmt.Sprintf("%s %d\x00%s", objType, len(content), string(content))
//		rawBytes := []byte(data)
//		hash := utils.Hash(rawBytes)
//		return rawBytes, hash
//	}
func HashObject(objType string, content []byte) ([]byte, string) {
	header := []byte(fmt.Sprintf("%s %d\x00", objType, len(content)))
	rawBytes := append(header, content...)
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

	return os.ReadFile(path)
}
func SplitByByte(data []byte, target byte) ([]byte, []byte) {
	for i, b := range data {
		if b == target {
			return data[:i], data[i+1:]
		}
	}
	return data, nil // target not found
}

func ParseObject(object []byte) (string, []byte) {
	header, content := SplitByByte(object, '\x00')
	objType := strings.Split(string(header), " ")[0]
	return objType, content
}

type TreeParsedEntry struct {
	Mode string
	Type string
	Hash string
	Name string
}

func ParseTreeContent(content []byte) string {
	r := bufio.NewReader(bytes.NewReader(content))
	var out strings.Builder

	for {
		entry, err := ReadTreeEntry(r)
		if err != nil {
			break
		}
		// format: 100644 blob b6fc4c620b67d95f    test.txt
		fmt.Fprintf(&out, "%06s %s %s\t%s\n", entry.Mode, entry.Type, entry.Hash, entry.Name)
	}

	return out.String()
}

func ReadTreeEntry(r *bufio.Reader) (TreeParsedEntry, error) {
	var e TreeParsedEntry

	modeBytes, err := r.ReadBytes(' ')
	if err != nil {
		return e, err
	}
	e.Mode = string(modeBytes[:len(modeBytes)-1])

	nameBytes, err := r.ReadBytes('\x00')
	if err != nil {
		return e, err
	}
	e.Name = string(nameBytes[:len(nameBytes)-1])

	var sha [20]byte
	if _, err := io.ReadFull(r, sha[:]); err != nil {
		return e, err
	}
	e.Hash = fmt.Sprintf("%x", sha)

	if e.Mode == "40000" {
		e.Type = "tree"
	} else {
		e.Type = "blob"
	}

	return e, nil
}
