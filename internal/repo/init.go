package repo

import (
	"fmt"
	"mygit/internal/utils"
	"os"
	"path/filepath"
	// "mygit/internal/utils"
)

func Init() error {

	var prfxName string
	if len(os.Args) > 2 {
		prfxName = os.Args[2]
	}

	workingDir, err := os.Getwd()

	if err != nil {
		fmt.Println("Error getting directory:", err)
		return err
	}

	// create the structure of .git

	path := filepath.Join(workingDir, prfxName, ".mygit")

	// fmt.Println("path:", path, "\n foldername", prfxName)

	//folders

	branchesFolder := filepath.Join(path, "branches")

	err = os.MkdirAll(branchesFolder, 0755)
	if err != nil {
		return err
	}

	objectsFolder := filepath.Join(path, "objects")

	err = os.MkdirAll(objectsFolder, 0755)
	if err != nil {
		return err
	}

	objectsInfoFolder := filepath.Join(path, "objects", "info")

	err = os.MkdirAll(objectsInfoFolder, 0755)
	if err != nil {
		return err
	}
	objectsPackFolder := filepath.Join(path, "objects", "pack")

	err = os.MkdirAll(objectsPackFolder, 0755)
	if err != nil {
		return err
	}

	hooksFolder := filepath.Join(path, "hooks")

	err = os.MkdirAll(hooksFolder, 0755)
	if err != nil {
		return err
	}

	refsFolder := filepath.Join(path, "refs")

	err = os.MkdirAll(refsFolder, 0755)
	if err != nil {
		return err
	}

	infoFolder := filepath.Join(path, "info")

	err = os.MkdirAll(infoFolder, 0755)
	if err != nil {
		return err
	}

	//files
	headFile := filepath.Join(path, "HEAD")
	err = utils.WriteFile(headFile, []byte("ref: refs/heads/master"))

	if err != nil {
		return err
	}

	descriptionFile := filepath.Join(path, "description")
	err = utils.WriteFile(descriptionFile, []byte("rename to give a description to the repository"))

	if err != nil {
		return err
	}

	return nil
}

