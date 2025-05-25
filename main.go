package main

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Tab struct {
	icon                 string
	shouldShowDefault    bool
	shouldShowOnLeftSide bool
}

func NewTab() *Tab {
	return &Tab{
		icon:                 "",
		shouldShowDefault:    false,
		shouldShowOnLeftSide: false,
	}
}

type ModuleMetadata struct {
	renderModes       []string
	moduleTitle       string
	moduleDescription string
	tab               *Tab
}

// given module_title -> Module Title
// given moduleTitle -> Module Title
func (mm *ModuleMetadata) makeModuleTitleHumanReadable() {
	mm.moduleTitle = cases.Title(language.English).String(strings.ReplaceAll(mm.moduleTitle, "_", " "))
}

func (mm *ModuleMetadata) getModuleNameKebabCase() string {
	return strings.ToLower(strings.ReplaceAll(mm.moduleTitle, " ", "-"))
}

func (mm *ModuleMetadata) getModuleElementNameKebabCase() string {
	return strings.ToLower(strings.ReplaceAll(mm.moduleTitle, " ", "-")) + "-element"
}

// turn Module Title -> ModuleTitle
func (mm *ModuleMetadata) getModuleNameClassName() string {
	return strings.ReplaceAll(cases.Title(language.English).String(mm.moduleTitle), " ", "") + "Element"
}

// turn Module Title -> moduleTitle
func (mm *ModuleMetadata) getModuleNameVariable() string {
	return strings.ReplaceAll(cases.Lower(language.English).String(mm.moduleTitle[:1])+mm.moduleTitle[1:], " ", "") + "Element"
}

// turn package_name
func (mm *ModuleMetadata) getModulePackageName() string {
	return "bookera-" + mm.getModuleNameKebabCase()
}

func NewModuleMetadata() *ModuleMetadata {
	return &ModuleMetadata{
		renderModes:       []string{},
		moduleTitle:       "",
		moduleDescription: "",
	}
}

func (mm *ModuleMetadata) hasSidePanel() bool {
	hasSidePanel := false
	for _, renderMode := range mm.renderModes {
		if renderMode == string(SIDE_PANEL) {
			hasSidePanel = true
		}
	}

	mm.tab = NewTab()

	return hasSidePanel
}

func main() {

	runForm()

}
