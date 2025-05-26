package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func createDirName(dirName string) string {
	if useDebugMode {
		test := "test"
		os.Mkdir(test, 0755)
		dirName = path.Join(test, dirName)
	} else {
		os.Mkdir(dirName, 0755)

	}

	return dirName
}

func cloneRepo(mm *ModuleMetadata) bool {
	repoURL := "https://github.com/serranolabs-io/bookera-module-template.git"

	dirName := mm.getModuleNameKebabCase()

	dirName = createDirName(dirName)

	cmd := exec.Command("git", "clone", repoURL, dirName)

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to clone repository: %v, please retry again", err)
	}

	// Remove .git directory
	err = os.RemoveAll(path.Join(dirName, ".git"))
	if err != nil {
		log.Fatalf("Failed to remove .git directory: %v", err)
	}

	// Remove .gitignore file
	err = os.Remove(path.Join(dirName, ".gitignore"))
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Failed to remove .gitignore file: %v", err)
	}

	return true
}

func renameFile(filePath string, search string, newString string) string {
	if strings.Contains(filePath, search) {
		newFilePath := strings.ReplaceAll(filePath, search, newString)
		err := os.Rename(filePath, newFilePath)
		if err != nil {
			log.Fatalf("Failed to rename file: %v", err)
		}
		filePath = newFilePath
	}

	return filePath
}

func (mm *ModuleMetadata) applyTemplateToFile(filePath string) {
	contents, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file %s", err)
	}

	strContents := string(contents)

	filePath = renameFile(filePath, "{module_name_kc}", mm.getModuleElementNameKebabCase())
	filePath = renameFile(filePath, "{module_element_kc}", mm.getModuleElementNameKebabCase())

	strContents = strings.ReplaceAll(strContents, "{package_name}", mm.getModulePackageName())
	strContents = strings.ReplaceAll(strContents, "{module_name_kc}", mm.getModuleNameKebabCase())
	strContents = strings.ReplaceAll(strContents, "{module_element_kc}", mm.getModuleElementNameKebabCase())
	strContents = strings.ReplaceAll(strContents, "$ModuleElementName", mm.getModuleNameClassName())
	strContents = strings.ReplaceAll(strContents, "$moduleElementName", mm.getModuleNameVariable())
	strContents = strings.ReplaceAll(strContents, "{module_name_hr}", mm.moduleTitle)
	strContents = strings.ReplaceAll(strContents, "{description}", mm.moduleDescription)
	strContents = strings.ReplaceAll(strContents, "`{renderModes}`", mm.renderRenderModes())

	if mm.tab != nil {
		strContents = strings.ReplaceAll(strContents, "{tab.icon}", mm.tab.icon)
		if mm.tab.shouldShowDefault {
			strContents = strings.ReplaceAll(strContents, ".removeTab()", "")
		}

		showOnLeftSide := "left"
		if !mm.tab.shouldShowOnLeftSide {
			showOnLeftSide = "right"
		}

		strContents = strings.ReplaceAll(strContents, "{shouldShowLeftSide}", showOnLeftSide)
	}

	debugPrint("file: " + filePath + "\n" + strContents)

	os.WriteFile(filePath, []byte(strContents), 0755)

}

func (mm *ModuleMetadata) templateRepo() {
	dirName := mm.getModuleNameKebabCase()
	dirName = createDirName(dirName)
	debugPrint("Templating repo")

	err := filepath.Walk(dirName, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			mm.applyTemplateToFile(filePath)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking through directory: %v", err)
	}

	debugPrint("Finished templating repo")

}
