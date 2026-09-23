package object

import (
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

	err = utils.WriteFileSafely(path, compressed)

	if err != nil {
		return "", err
	}

	return hash, nil
}

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

func Content(objectHash string) ([]byte, error) {

	data, err := ReadObject(objectHash)

	if err != nil {
		return []byte{}, err
	}

	_, content, err := ParseObject(data)

	if err != nil {
		return []byte{}, err
	}

	return content, nil

}
