package cli

import (
	"mygit/internal/index"
	"os"

	flag "github.com/spf13/pflag"
)

func UpdateIndexCommand() error {
	cmd := flag.NewFlagSet("update-index", flag.ExitOnError)
	add := cmd.Bool("add", false, "add  to index")
	remove := cmd.Bool("remove", false, "remove from index")
	// refresh := cmd .Bool("refresh",false,"refresh stat info without changing hashes")
	cacheInfo := cmd.String("cacheinfo", "", "add an entry directly")

	cmd.Parse(os.Args[2:])

	if *add {
		if *cacheInfo != "" {

		}

		if cmd.NArg() < 1 {
			return usageError("usage: mygit update-index --add <file>")
		}

		path := cmd.Arg(0)

		return index.UpdateIndexAdd(path)
	}

	if *remove {
		if cmd.NArg() < 1 {
			return usageError("usage: mygit update-index --remove <file>")
		}

		path := cmd.Arg(0)

		return index.UpdateIndexRemove(path)
	}

	if *cacheInfo != "" {
		return index.UpdateIndexCacheInfo(*cacheInfo)
	}

	return usageError("usage: mygit update-index [--add|--remove|--cacheinfo]")
}
