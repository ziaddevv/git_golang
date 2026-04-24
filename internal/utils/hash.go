package utils

import (
	"crypto/sha1"
	"fmt"
)

func Hash(data []byte) string {
	sum := sha1.Sum(data)
	return fmt.Sprintf("%x", sum)
}
