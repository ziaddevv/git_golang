package cli

import (
	"fmt"
	"mygit/internal/repo"
	"mygit/internal/utils"
	"os"
)

func InitCommand() {
	if len(os.Args) > 3 {
		fmt.Println("init takes mx  2 arguments")
		return
	}

	if utils.Exists(".mygit") {
		fmt.Println("repo already exists")
		return
	}

	err := repo.Init()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("Repository initialized")
}
