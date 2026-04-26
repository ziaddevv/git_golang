package cli

import (
	"fmt"
	"mygit/internal/repo"
	"os"

	flag "github.com/spf13/pflag"
)

/*
flags

-t ----> print type
-s ---> print size
-p ----> if blob --> print content , else ---> print raw payload as text/bytes
*/
func CatFileCommand() {
	catFileCmd := flag.NewFlagSet("cat-file", flag.ExitOnError)

	typ := catFileCmd.BoolP("type", "t", false, "print type of the object")
	pretty := catFileCmd.BoolP("pretty", "p", false, "pretty print object content")
	size := catFileCmd.BoolP("size", "s", false, "print size of content")

	catFileCmd.Parse(os.Args[2:])

	if catFileCmd.NArg() < 1 {
		fmt.Println("usage: mygit cat-file [-t|-p|-s] <hash>")
		return
	}

	hash := catFileCmd.Arg(0)

	data, err := repo.ReadObject(hash)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	objType, content := repo.ParseObject(data)

	switch {
	case *typ:
		fmt.Println(objType)

	case *size:
		fmt.Println(len(content))

	case *pretty:
		switch objType {
		case "blob":
			fmt.Print(string(content))
		case "tree":

		case "commit":

		case "tag":

		default:
			fmt.Print(string(content))
		}

	default:
		fmt.Println("usage: mygit cat-file [-t|-p|-s] <hash>")
	}
}
