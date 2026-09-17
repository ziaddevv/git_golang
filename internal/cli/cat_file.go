package cli

import (
	"fmt"
	"mygit/internal/object"
	"os"

	flag "github.com/spf13/pflag"
)

/*
flags

-t ----> print type
-s ---> print size
-p ----> if blob --> print content , else ---> print raw payload as text/bytes
*/
func CatFileCommand() error {
	catFileCmd := flag.NewFlagSet("cat-file", flag.ExitOnError)

	typ := catFileCmd.BoolP("type", "t", false, "print type of the object")
	pretty := catFileCmd.BoolP("pretty", "p", false, "pretty print object content")
	size := catFileCmd.BoolP("size", "s", false, "print size of content")

	catFileCmd.Parse(os.Args[2:])

	if catFileCmd.NArg() < 1 {
		return usageError("usage: mygit cat-file [-t|-p|-s] <hash>")
	}

	hash := catFileCmd.Arg(0)

	data, err := object.ReadObject(hash)
	if err != nil {
		return err
	}
	objType, content, err := object.ParseObject(data)
	if err != nil {
		return err
	}

	switch {
	case *typ:
		fmt.Println(objType)

	case *size:
		fmt.Println(len(content))

	case *pretty:
		switch objType {
		case "blob":
			os.Stdout.Write(content)
		case "tree":
			// fmt.Print(string(content))
			fmt.Print(object.ParseTreeContent(content))
		case "commit":
			os.Stdout.Write(content)

		case "tag":
			fmt.Print(content)

		default:
			fmt.Print(string(content))
		}

	default:
		return usageError("usage: mygit cat-file [-t|-p|-s] <hash>")
	}

	return nil
}
