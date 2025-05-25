package main

import (
	"fmt"
	"image/color"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// A message used to indicate that activity has occurred. In the real world (for
// example, chat) this would contain actual data.

const CREATE_MODULE = "Cloning repo"
const CREATE_TEMPLATE = "Creating template"
const IS_DONE_MESSAGE = "Enjoy "

type step int

const (
	IS_CLONING step = iota
	IS_TEMPLATING
	IS_DONE
)

type responseMsg struct {
	message string
	step    step
}

// Simulate a process that sends events at an irregular interval in real time.
// In this case, we'll send events on the channel at a random interval between
// 100 to 1000 milliseconds. As a command, Bubble Tea will run this
// asynchronously.
func (m *Model) cloneRepo() {
	cloneRepo(m.moduleMetadata)

	m.sub <- responseMsg{message: CREATE_TEMPLATE, step: IS_TEMPLATING}
}

func (m *Model) templateRepo() {
	m.moduleMetadata.templateRepo()

	m.sub <- responseMsg{message: IS_DONE_MESSAGE, step: IS_DONE}
	close(m.sub)
}

// A command that waits for the activity on a channel.
func waitForActivity(sub chan responseMsg) tea.Cmd {
	return func() tea.Msg {
		rm := <-sub // Wait for activity

		return rm
	}
}

// A Model can be more or less any type of data. It holds all the data for a
// program, so often it's a struct. For this simple example, however, all
// we'll need is a simple integer.

type MessageGradient struct {
	message string
	blend   *[]color.Color
}

func (mg *MessageGradient) updateMessage(message string) {
	mg.message = message
}

func (mg *MessageGradient) rotateBlend() {
	mg.blend = rotateBlend(mg.blend)
}

func NewMessageGradient(message string) *MessageGradient {
	return &MessageGradient{message: message, blend: createBlendP(message)}

}

type Model struct {
	sub             chan responseMsg
	tick            int
	step            step
	messageGradient *MessageGradient
	moduleMetadata  *ModuleMetadata
}

// Init optionally returns an initial command we should run. In this case we
// want to start the timer.
func (m *Model) Init() tea.Cmd {
	m.messageGradient = NewMessageGradient(CREATE_MODULE + ".")

	m.sub = make(chan responseMsg)
	go m.cloneRepo()

	return tea.Batch(tick, waitForActivity(m.sub))
}

const divisor = 9

var total = divisor * 3

func debugPrint(str string) {
	f, err := os.OpenFile("debug.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(str + "\n"); err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

// Update is called when messages are received. The idea is that you inspect the
// message and send back an updated Model accordingly. You can also return
// a command, which is a function that performs I/O and returns a message.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		debugPrint("key")
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+z":
			return m, tea.Suspend
		}

	case responseMsg:
		m.messageGradient = NewMessageGradient(msg.message)
		m.step = msg.step

		debugPrint(fmt.Sprintf("received message: %d", m.step))

		if m.step == IS_TEMPLATING {
			go m.templateRepo()
			return m, tea.Batch(tick, waitForActivity(m.sub))
		} else if m.step == IS_DONE {
			return m, tea.Quit
		}

		return m, tick

	case tickMsg:
		debugPrint("tick")
		m.tick++

		if m.step == IS_CLONING {
			if m.tick%total < divisor {
				m.messageGradient.updateMessage(CREATE_MODULE + ".  ")
			} else if m.tick%total < divisor*2 {
				m.messageGradient.updateMessage(CREATE_MODULE + ".. ")
			} else {
				m.messageGradient.updateMessage(CREATE_MODULE + "...")
			}
		} else if m.step == IS_TEMPLATING {
			if m.tick%total < divisor {
				m.messageGradient.updateMessage(CREATE_TEMPLATE + ".  ")
			} else if m.tick%total < divisor*2 {
				m.messageGradient.updateMessage(CREATE_TEMPLATE + ".. ")
			} else {
				m.messageGradient.updateMessage(CREATE_TEMPLATE + "...")
			}
		}

		m.messageGradient.rotateBlend()
		return m, tick
	}
	return m, nil
}

// View returns a string based on data in the Model. That string which will be
// rendered to the terminal.
func (m *Model) View() string {
	return makeGradientWithBlend(m.messageGradient.message, lipgloss.NewStyle().MarginBottom(2), *m.messageGradient.blend)
}

// Messages are events that we respond to in our Update function. This
// particular one indicates that the timer has ticked.
type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Millisecond * 50)
	return tickMsg{}
}
