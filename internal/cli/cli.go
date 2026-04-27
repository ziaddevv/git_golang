package cli

import (
	"fmt"
	"os"
)

// fun command parser
func Run() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: mygit <command>")
		return
	}

	switch os.Args[1] {
	case "init":
		InitCommand()
	case "hash-object": // a genirec object creator , like it can work with all objects + it does't unserstand the high level semantics , it just takes bytes and type and it write it in objects folder
		//TODO make the command takes input from file  or std
		// git hash-object file.txt      --> just hash object
		// git hash-object -w file.txt   --> hash and write
		// echo -n "hello" | git hash-object --stdin    ---> pipe text directly
		//echo -n "hello" | git hash-object -w --stdin  --> pipe text and write the object
		//printf "a.txt\nb.txt\n" | git hash-object --stdin-paths --> path names from stdin
		// hash, err := repo.WriteObject("blob", []byte("hello"))
		HashObjectCommand()
		// fmt.Println(hash, err)
	case "cat-file":
		CatFileCommand()

	case "update-index":
		// UpdateIndex()
		// fmt.Println(utils.ReadFile("/home/ziad/projects/test/tests/.git/index"))
		UpdateIndexCommand()
	}
}

/*


// check if there is a provided foldername
		var prfxName string
		if len(os.Args) > 2 {
			prfxName = os.Args[2]
		}

		// get current path

		working_dir, err := os.Getwd()

		if err != nil {
			fmt.Println("Error getting directory:", err)
			return
		}

		fmt.Println("Current working directory:", working_dir, "\n foldername", prfxName)
		// create the structure of .git

		// 0755 is the standard Unix permission (rwxr-xr-x)

		path := filepath.Join(prfxName, ".mygit")

		branchesFolder := filepath.Join(path, "branches")

		err = os.MkdirAll(branchesFolder, 0755)
		if err != nil {
			panic(err)
		}

		objectsFolder := filepath.Join(path, "objects")

		err = os.MkdirAll(objectsFolder, 0755)
		if err != nil {
			panic(err)
		}

		objectsInfoFolder := filepath.Join(path, "objects", "info")

		err = os.MkdirAll(objectsInfoFolder, 0755)
		if err != nil {
			panic(err)
		}
		objectsPackFolder := filepath.Join(path, "objects", "pack")

		err = os.MkdirAll(objectsPackFolder, 0755)
		if err != nil {
			panic(err)
		}

		hooksFolder := filepath.Join(path, "hooks")

		err = os.MkdirAll(hooksFolder, 0755)
		if err != nil {
			panic(err)
		}

		refsFolder := filepath.Join(path, "refs")

		err = os.MkdirAll(refsFolder, 0755)
		if err != nil {
			panic(err)
		}

		infoFolder := filepath.Join(path, "info")

		err = os.MkdirAll(infoFolder, 0755)
		if err != nil {
			panic(err)
		}

		// err = os.MkdirAll(path, 0755)
		// if err != nil {
		// 	panic(err)
		// }

		// err = os.MkdirAll(path, 0755)
		// if err != nil {
		// 	panic(err)
		// }
		// err = os.MkdirAll(path, 0755)
		// if err != nil {
		// 	panic(err)
		// }

		// put folders inside





*/
