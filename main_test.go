package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {

	main()
	assert.Equal(t, 123, 123, "EQUAL")

}

func setupModuleMetadata() *ModuleMetadata {
	mm := &ModuleMetadata{moduleTitle: "module_title"}
	mm.makeModuleTitleHumanReadable()
	return mm
}

func TestMakeModuleTitleHumanReadable(t *testing.T) {
	mm := setupModuleMetadata()
	assert.Equal(t, "Module Title", mm.moduleTitle, "Module title should be human-readable")
}

func TestGetModuleNameKebabCase(t *testing.T) {
	mm := setupModuleMetadata()
	result := mm.getModuleNameKebabCase()
	assert.Equal(t, "module-title", result, "Module name should be in kebab-case")
}

func TestGetModuleElementNameKebabCase(t *testing.T) {
	mm := setupModuleMetadata()
	result := mm.getModuleElementNameKebabCase()
	assert.Equal(t, "module-title-element", result, "Module element name should be in kebab-case with '-element' suffix")
}

func TestGetModuleNameClassName(t *testing.T) {
	mm := setupModuleMetadata()
	result := mm.getModuleNameClassName()
	assert.Equal(t, "ModuleTitleElement", result, "Module name should be in class name format with 'Element' suffix")
}

func TestGetModuleNameVariable(t *testing.T) {
	mm := setupModuleMetadata()
	result := mm.getModuleNameVariable()
	assert.Equal(t, "moduleTitleElement", result, "Module name should be in variable format with 'Element' suffix")
}

func TestGetModulePackageName(t *testing.T) {
	mm := setupModuleMetadata()
	result := mm.getModulePackageName()
	assert.Equal(t, "bookera-module-title", result, "Module package name should be prefixed with 'bookera-' and in kebab-case")
}
