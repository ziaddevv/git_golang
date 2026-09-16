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

func HashObjectCommand() {
	cmd := flag.NewFlagSet("hash-object", flag.ExitOnError)

	write := cmd.BoolP("write", "w", false, "write object into database")
	stdin := cmd.Bool("stdin", false, "read content from stdin")
	objType := cmd.StringP("type", "t", "blob", "object type")

	cmd.Parse(os.Args[2:])

	valid := map[string]bool{
		"blob":   true,
		"tree":   true,
		"commit": true,
		"tag":    true,
	}

	if !valid[*objType] {
		fmt.Fprintf(os.Stderr, "invalid object type: %s\n", *objType)
		os.Exit(1)
	}

	var data []byte
	var err error

	if *stdin {
		data, err = io.ReadAll(os.Stdin)
	} else {
		if cmd.NArg() < 1 {
			fmt.Fprintln(os.Stderr, "missing file operand")
			os.Exit(1)
		}

		fileName := cmd.Arg(0)
		data, err = utils.ReadFile(fileName)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var hash string

	if *write {
		fmt.Println("writing ", data)
		hash, err = object.WriteObject(*objType, data)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	} else {
		_, hash = object.HashObject(*objType, data)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(hash)
}
