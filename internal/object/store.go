package object

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"mygit/internal/utils"
	"os"
	"strconv"
	"strings"
)

//   converting from []byte to string ---> []byte(mystring)
//   converting frokm string to []bytes ---> string(bybytes)

// [type] + " " + [size] + \0 + [content]          , sh1 hash is the object id

//	func HashObject(objType string, content []byte) ([]byte, string) {
//		data := fmt.Sprintf("%s %d\x00%s", objType, len(content), string(content))
//		rawBytes := []byte(data)
//		hash := utils.Hash(rawBytes)
//		return rawBytes, hash
//	}
func IsValidObjectType(objType string) bool {
	switch objType {
	case "blob", "tree", "commit", "tag":
		return true
	default:
		return false
	}
}

func IsValidObjectID(hash string) bool {
	if len(hash) != 40 {
		return false
	}

	for _, c := range hash {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false
		}
	}

	return true
}

func HashObject(objType string, content []byte) ([]byte, string) {
	header := []byte(fmt.Sprintf("%s %d\x00", objType, len(content)))
	rawBytes := append(header, content...)
	hash := utils.Hash(rawBytes)
	return rawBytes, hash
}

func WriteObject(objType string, content []byte) (string, error) {

	if !IsValidObjectType(objType) {
		return "", fmt.Errorf("invalid object type: %s", objType)
	}

	if !utils.Exists(utils.RepoDir()) {
		return "", errors.New("not a repository")
	}

	rawBytes, hash := HashObject(objType, content)
	path, err := utils.ObjectPath(hash)

	if err != nil {
		return "", errors.New("object path is incorrect")
	}

	fmt.Println(path)
	if utils.Exists(path) {
		return hash, nil
	}

	var buf bytes.Buffer

	writer := zlib.NewWriter(&buf)
	if _, err := writer.Write(rawBytes); err != nil {
		writer.Close()
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	compressed := buf.Bytes()

	fmt.Println("Original:", len(rawBytes))
	fmt.Println("Compressed:", len(compressed))

	err = utils.WWriteFileSafely(path, compressed)

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
	if !IsValidObjectID(hash) {
		return nil, fmt.Errorf("invalid object id: %s", hash)
	}

	path, err := utils.ObjectPath(hash)

	if err != nil {
		return []byte{}, errors.New("object path is incorrect")
	}

	compressed, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	reader, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	decompressed, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if _, _, err := ParseObject(decompressed); err != nil {
		return nil, err
	}

	actualHash := utils.Hash(decompressed)
	if !strings.EqualFold(actualHash, hash) {
		return nil, fmt.Errorf("object hash mismatch: expected %s, got %s", hash, actualHash)
	}

	return decompressed, nil
}
func SplitByByte(data []byte, target byte) ([]byte, []byte, bool) {
	for i, b := range data {
		if b == target {
			return data[:i], data[i+1:], true
		}
	}
	return data, nil, false // target not found
}

func ParseObject(object []byte) (string, []byte, error) {
	header, content, ok := SplitByByte(object, '\x00')
	if !ok {
		return "", nil, errors.New("invalid object: missing header terminator")
	}

	parts := strings.Split(string(header), " ")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", nil, fmt.Errorf("invalid object header: %q", string(header))
	}

	objType := parts[0]
	if !IsValidObjectType(objType) {
		return "", nil, fmt.Errorf("invalid object type: %s", objType)
	}

	size, err := strconv.Atoi(parts[1])
	if err != nil || size < 0 {
		return "", nil, fmt.Errorf("invalid object size: %s", parts[1])
	}

	if size != len(content) {
		return "", nil, fmt.Errorf("object size mismatch: header says %d, content has %d", size, len(content))
	}

	return objType, content, nil
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
