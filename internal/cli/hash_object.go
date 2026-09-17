package cli

import (
	"fmt"
	"io"
	"mygit/internal/object"
	"mygit/internal/utils"
	"os"

	flag "github.com/spf13/pflag" // standard flag is bad for me
)

// parse command
// validate command structure and flags
// what type os objects? -t flag ,

func HashObjectCommand() error {
	cmd := flag.NewFlagSet("hash-object", flag.ExitOnError)

	write := cmd.BoolP("write", "w", false, "write object into database")
	stdin := cmd.Bool("stdin", false, "read content from stdin")
	objType := cmd.StringP("type", "t", "blob", "object type")

	cmd.Parse(os.Args[2:])

	if !object.IsValidObjectType(*objType) {
		return usageError(fmt.Sprintf("invalid object type: %s", *objType))
	}

	var data []byte
	var err error

	if *stdin {
		data, err = io.ReadAll(os.Stdin)
	} else {
		if cmd.NArg() < 1 {
			return usageError("missing file operand")
		}

		fileName := cmd.Arg(0)
		data, err = utils.ReadFile(fileName)
	}

	if err != nil {
		return err
	}

	var hash string

	if *write {
		fmt.Println("writing ", data)
		hash, err = object.WriteObject(*objType, data)
		if err != nil {
			return err
		}
	} else {
		_, hash = object.HashObject(*objType, data)
	}

	fmt.Println(hash)
	return nil
}
