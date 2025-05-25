package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func cloneRepo(mm *ModuleMetadata) bool {
	repoURL := "https://github.com/serranolabs-io/bookera-module-template.git"
	dirName := mm.getModuleNameKebabCase()
	fmt.Println("clone repo")

	err := os.Mkdir(dirName, 0755)
	if err != nil {
		log.Fatalf("Failed to create directory: %v, please retry again", err)
	}

	cmd := exec.Command("git", "clone", repoURL, dirName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to clone repository: %v, please retry again", err)
	}

	return true
}
