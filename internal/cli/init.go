package cli

import (
	"fmt"
	"mygit/internal/repo"
	"mygit/internal/utils"
	"os"
)

func InitCommand() error {
	if len(os.Args) > 3 {
		return usageError("init takes max 2 arguments")
	}

	if utils.Exists(".mygit") {
		return fmt.Errorf("repo already exists")
	}

	err := repo.Init()
	if err != nil {
		return fmt.Errorf("init: %w", err)
	}

	fmt.Println("Repository initialized")
	return nil
}
