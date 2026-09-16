package cli

import (
	"fmt"
	"mygit/internal/index"
	"os"

	flag "github.com/spf13/pflag"
)

func UpdateIndexCommand() {
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
			fmt.Println("usage: mygit update-index --add <file>")
			return
		}
		
		path := cmd.Arg(0)
		
		err := index.UpdateIndexAdd(path)
		if err != nil {
			fmt.Println("error:", err)
		}
		return

	}

	if *remove {
		if cmd.NArg() < 1 {
			fmt.Println("usage: mygit update-index --remove <file>")
			return
		}

		path := cmd.Arg(0)

		err := index.UpdateIndexRemove(path)
		if err != nil {
			fmt.Println("error:", err)
		}
		return

	}

	if *cacheInfo != "" {
		err := index.UpdateIndexCacheInfo(*cacheInfo)
		if err != nil {
			fmt.Println("error:", err)
		}
		return
	}
	fmt.Println("usage: mygit update-index [--add|--remove|--cacheinfo]")
}

