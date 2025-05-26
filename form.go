package main

import (
	"errors"
	"fmt"
	"image/color"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/gamut"
)

var primaryColor = lipgloss.Color("#ffb87e")
var textColor = lipgloss.NewStyle().Foreground(primaryColor)
var backgroundColor = lipgloss.NewStyle().Background(primaryColor)

const MAX_TITLE_LENGTH = 25
const MAX_DESCRIPTION_LENGTH = 400

type RenderMode string

const (
	SIDE_PANEL    RenderMode = "renderInSidePanel"
	MODULE_DAEMON RenderMode = "renderInDaemon"
	PANEL         RenderMode = "renderInPanel"
	SETTINGS      RenderMode = "renderInSettings"
)

func createBlend(str string) []color.Color {
	return gamut.Blends(primaryColor, lipgloss.Color("#4c9999"), len(str))
}

func createBlendP(str string) *[]color.Color {
	blend := gamut.Blends(primaryColor, lipgloss.Color("#4c9999"), len(str))
	return &blend
}

func printGradient(str string, style lipgloss.Style) {
	makeGradient(str, style)
}

func makeGradient(str string, style lipgloss.Style) string {
	base := lipgloss.NewStyle()
	return (style.Render(rainbow(base, str, createBlend(str))))
}

func makeGradientWithBlend(str string, style lipgloss.Style, blend []color.Color) string {
	base := lipgloss.NewStyle()
	return (style.Render(rainbow(base, str, blend)))
}

func runFirstStep() *ModuleMetadata {

	mm := NewModuleMetadata()

	welcome := "Welcome to the Bookera Module TUI"
	printGradient(welcome, lipgloss.NewStyle().Padding(2).MarginBottom(2))

	form := huh.NewForm(
		// Gather some final details about the order.
		huh.NewGroup(
			huh.NewInput().
				Title(textColor.Render("What’s the title of your module?")).
				Value(&mm.moduleTitle).

				// Validating fields is easy. The form will mark erroneous fields
				// and display error messages accordingly.
				Validate(func(str string) error {
					for _, char := range str {
						if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == ' ') {
							return errors.New("please make a name with that contains only letters in the alphabet or a space")
						}
					}
					if len(str) > MAX_TITLE_LENGTH {
						return fmt.Errorf("sorry, name is too long, max is %d and yours is %d", MAX_TITLE_LENGTH, len(str))
					}
					return nil
				}),
			huh.NewText().
				Title(textColor.Render("Please provide a description for your module!")).
				Value(&mm.moduleDescription),

			huh.NewMultiSelect[string]().
				Title(textColor.Render("Render Modes (Where would you like your module rendered?)")).
				Options(
					huh.NewOption("Side Panel - Add a tab to your module so it can be viewed in side panel", string(SIDE_PANEL)),
					huh.NewOption("Module Daemon - for event listeners & what not", string(MODULE_DAEMON)),
					huh.NewOption("Panel - classic", string(PANEL)),
					huh.NewOption("Settings - Put the settings in your module here", string(SETTINGS)),
				).
				Value(&mm.renderModes).Validate(func(s []string) error {
				if len(s) == 0 {
					return errors.New("you must select one, lol")
				}

				return nil
			}),
		),
	)

	form.WithTheme(huh.ThemeBase())

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	mm.makeModuleTitleHumanReadable()
	return mm
}

func runSidePanelStep() *Tab {

	tab := NewTab()

	welcome := backgroundColor.Padding(1).MarginBottom(2).Render("Tabs")
	fmt.Println(welcome)

	form := huh.NewForm(
		// Gather some final details about the order.
		huh.NewGroup(
			huh.NewInput().
				Title(textColor.Render("What icon would you like to use? Icons found here https://shoelace.style/components/icon/")).
				Value(&tab.icon),
			huh.NewConfirm().
				Title(textColor.Render("Would you like to show the tab on default?")).
				Value(&tab.shouldShowDefault),
			huh.NewConfirm().
				Title(textColor.Render("Would you like to place this tab on the left or right side")).
				Affirmative("Show on left side").Negative("Show on right side").
				Value(&tab.shouldShowOnLeftSide),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	return tab
}

func rainbow(base lipgloss.Style, s string, colors []color.Color) string {
	var str string
	for i, ss := range s {
		color, _ := colorful.MakeColor(colors[i%len(colors)])
		str = str + base.Foreground(lipgloss.Color(color.Hex())).Render(string(ss))
	}
	return str
}

func rotateBlend(blend *[]color.Color) *[]color.Color {
	if len(*blend) == 0 {
		return blend
	}

	previousColor := (*blend)[0]
	var nextColor color.Color
	for i := 0; i < len(*blend)-1; i++ {
		nextColor = (*blend)[i+1]
		(*blend)[i+1] = previousColor
		previousColor = nextColor
	}

	(*blend)[0] = nextColor

	return blend
}

func runForm() {
	mm := runFirstStep()

	if mm.hasSidePanel() {
		mm.tab = runSidePanelStep()
	}

	p := tea.NewProgram((&Model{moduleMetadata: mm}))

	if _, err := p.Run(); err != nil {
		panic("Error in running form")
	}

}
